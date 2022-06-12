package database

import (
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client = nil

var urlEncoding b64.Encoding = *b64.URLEncoding.WithPadding(b64.NoPadding)

func GetFreeURLID() (pageID int32, err error) {
	db, err := GetDatabase()
	if err != nil {
		panic(err)
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}

	cursor, err := db.Collection("free_pages").Aggregate(context.TODO(), pipeline)
	if err != nil {
		return -1, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	pageID, ok := results[0]["_id"].(int32)

	if !ok {
		err = errors.New("failed to retrieve _id from free_page result")
	}
	return
}

type IllegalURLError string

func (e IllegalURLError) Error() string {
	return "URL input is not allowed: " + string(e)
}

type URLTakenError string

func (e URLTakenError) Error() string {
	return "URL is not unique" + string(e)
}

// func createPage(alias string) (err error) {
// 	urlID, err := GetFreeURLID()
// 	if err != nil {
// 		return err
// 	}

// 	return
// }

func GetClient() (client *mongo.Client, err error) {
	if mongoClient != nil {
		return mongoClient, nil
	}
	usernameFile := os.Getenv("MONGO_INITDB_ROOT_USERNAME_FILE")
	passwordFile := os.Getenv("MONGO_INITDB_ROOT_PASSWORD_FILE")
	hostName := os.Getenv("HOSTNAME")

	if usernameFile == "" {
		mongoClient = nil
		return mongoClient, errors.New("MONGO_INITDB_ROOT_USERNAME_FILE env variable empty")
	}
	if passwordFile == "" {
		mongoClient = nil
		return mongoClient, errors.New("MONGO_INITDB_ROOT_PASSWORD_FILE env variable empty")
	}
	if hostName == "" {
		mongoClient = nil
		return mongoClient, errors.New("HOSTNAME env variable empty")
	}

	usernameData, err := os.ReadFile(usernameFile)
	if err != nil {
		mongoClient = nil
		return mongoClient, err
	}
	passwordData, err := os.ReadFile(passwordFile)
	if err != nil {
		mongoClient = nil
		return mongoClient, err
	}
	connectionUrl := fmt.Sprintf("mongodb://%s:%s@%s:27017/?authSource=admin", string(usernameData), string(passwordData), hostName)
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connectionUrl))
	return mongoClient, err
}

func GetDatabase() (database *mongo.Database, err error) {
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		return nil, errors.New("MONGO_INITDB_DATABASE env variable empty")
	}
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	database = client.Database(databaseName)
	return
}

// func ReserveNewPage() (URL string, err error) {
// 	db, err := GetDatabase()
// 	if err != nil {
// 		return
// 	}

// 	filter := bson.D{{Key: "_id", Value: urlID}}
// 	update := bson.D{
// 		{Key: "$set", Value: bson.D{
// 			{Key: "dateAdded", Value: time.Now()},
// 		},
// 		},
// 	}
// 	err = db.Collection("pages").FindOneAndUpdate(
// 		context.TODO(),
// 		filter,
// 		update,
// 		// opts,
// 	).Decode(&updatedDocument)
// 	return
// }
