package database

import (
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"linkshare_api/conf"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var urlEncoding b64.Encoding = *b64.URLEncoding.WithPadding(b64.NoPadding)

func createNewPage(alias string) (pageID int32, err error) {
	client, err := GetClient()
	if err != nil {
		err = fmt.Errorf("no client: %w", err)
		return
	}
	defer client.Disconnect(context.TODO())
	db, err := GetDatabase(client)
	if err != nil {
		err = fmt.Errorf("no db: %w", err)
		return
	}
	unusedPages := db.Collection("unusedPagesIDs")
	pages := db.Collection("pages")

	// Get a random page ID
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}
	cursor, err := unusedPages.Aggregate(context.TODO(), pipeline)
	if err != nil {
		err = fmt.Errorf("aggregation failure: %w", err)
		return
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		err = fmt.Errorf("failed to grab results from cursor: %w", err)
		return
	}
	pageID, ok := results[0]["_id"].(int32)
	if !ok {
		err = errors.New("failed to retrieve _id from free_page result")
		return
	}

	// Delete the random page ID from unusedPages
	var deletedDocument = make(bson.M)
	err = unusedPages.FindOneAndDelete(context.TODO(), bson.M{"_id": pageID}).Decode(deletedDocument)
	if err != nil {
		// we don't care if we discard this ID, we have many so just return
		err = fmt.Errorf("delete error: %w", err)
		return
	}

	updateMap := bson.M{
		"_id":       pageID,
		"dateAdded": primitive.NewDateTimeFromTime(time.Now()),
		"schema":    conf.GetConf().SchemaVersion,
	}

	if len(alias) > 0 {
		updateMap["alias"] = alias
	}

	res, err := pages.InsertOne(context.TODO(), updateMap)

	if err != nil {
		err = fmt.Errorf("insertion error: %w", err)
		return
	}

	return res.InsertedID.(int32), err
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
	connectionUrl := conf.GetConf().ConnectionURL
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUrl))
	return mongoClient, err
}

func GetDatabase(client *mongo.Client) (database *mongo.Database, err error) {
	database = client.Database(conf.GetConf().DBName)
	return
}
