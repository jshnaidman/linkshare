package database

import (
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var urlEncoding b64.Encoding = *b64.URLEncoding.WithPadding(b64.NoPadding)

func createNewPage(alias string) (updatedDocument bson.M, err error) {
	client, err := GetClient()
	if err != nil {
		return
	}
	defer client.Disconnect(context.TODO())
	db, err := GetDatabase(client)
	if err != nil {
		return
	}
	pages := db.Collection("pages")

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "dateAdded", Value: nil}}}},
		bson.D{{Key: "$limit", Value: 10000}},
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}
	cursor, err := pages.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation failure: %w", err)
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	pageID, ok := results[0]["_id"].(int32)
	if !ok {
		return nil, errors.New("failed to retrieve _id from free_page result")
	}
	updateMap := bson.M{
		"$set": bson.M{
			"dateAdded": primitive.NewDateTimeFromTime(time.Now()),
		},
	}
	if len(alias) > 0 {
		updateMap["$set"].(bson.M)["alias"] = alias
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := pages.FindOneAndUpdate(context.TODO(), bson.M{"_id": pageID}, updateMap, opts)
	err = res.Decode(&updatedDocument)

	return
}

func GetURL(urlID int) (URL string) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(urlID))
	return urlEncoding.EncodeToString(data)
}

type IllegalURLError string

func (e IllegalURLError) Error() string {
	return "URL input is not allowed: " + string(e)
}

type URLTakenError string

func (e URLTakenError) Error() string {
	return "URL is not unique" + string(e)
}

// db.pages.aggregate([{$match: {"dateAdded": null}}, {$limit: 5000}, {$sample: {size: 1}}])

func GetClient() (client *mongo.Client, err error) {
	username := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	password := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	hostName := os.Getenv("HOSTNAME")

	if hostName == "" {
		return nil, errors.New("HOSTNAME env variable empty")
	}

	connectionUrl := fmt.Sprintf("mongodb://%s:%s@%s:27017/?authSource=admin", username, password, hostName)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUrl))
	return mongoClient, err
}

func GetDatabase(client *mongo.Client) (database *mongo.Database, err error) {
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		return nil, errors.New("MONGO_INITDB_DATABASE env variable empty")
	}
	database = client.Database(databaseName)
	return
}
