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

type IllegalURLError string

func (e IllegalURLError) Error() string {
	return "URL input is not allowed: " + string(e)
}

type URLTakenError string

func (e URLTakenError) Error() string {
	return "URL is taken: " + string(e)
}

// This func doesn't validate if page is valid base64 encoding
func createNewPage(alias string) (pageID uint32, err error) {
	// TODO - fix this

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
	hasAliasArg := len(alias) > 0

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
	tmpID, ok := results[0]["_id"].(int32)
	if !ok {
		err = errors.New("failed to retrieve _id from free_page result")
		return
	}
	pageID = uint32(tmpID)

	// Let's delete the random page and check if the decoded pageID / alias are taken in parallel
	deleteErrChan := make(chan error, 1)
	go func() {
		// Delete the random page ID from unusedPages
		var deletedDocument = make(bson.M)
		err = unusedPages.FindOneAndDelete(context.TODO(), bson.M{"_id": pageID}).Decode(deletedDocument)
		if err != nil {
			// we don't care if we discard this ID, we have many so just return
			deleteErrChan <- fmt.Errorf("delete error: %w", err)
		} else {
			deleteErrChan <- nil
		}
	}()

	pageURL := EncodePageID(pageID)

	var checkAliasQuery bson.M

	if hasAliasArg {
		checkAliasQuery = bson.M{
			"$or": []bson.M{
				{"alias": pageURL},
				{"alias": alias},
			},
		}
	} else {
		checkAliasQuery = bson.M{
			"alias": pageURL,
		}
	}

	result := pages.FindOne(context.TODO(), checkAliasQuery)
	if result == nil {
		//todo
		return
	}

	updateMap := bson.M{
		"_id":       pageID,
		"dateAdded": primitive.NewDateTimeFromTime(time.Now()),
		"schema":    conf.GetConf().SchemaVersion,
	}

	if hasAliasArg {
		updateMap["alias"] = alias
	}

	res, err := pages.InsertOne(context.TODO(), updateMap)

	if err != nil {
		err = fmt.Errorf("insertion error: %w", err)
		return
	}

	return uint32(res.InsertedID.(int32)), err
}

func EncodePageID(pageID uint32) (URL string) {
	data := make([]byte, 4)
	// mongodb stores in bson which is little endian except for timestamp and counter
	binary.LittleEndian.PutUint32(data, pageID)
	return urlEncoding.EncodeToString(data)
}

func DecodeURL(URL string) (pageID uint32, err error) {
	decodedURL, err := urlEncoding.DecodeString(URL)
	if len(decodedURL) > 4 {
		err = errors.New("decoded URL does not fit into int32")
		return
	}
	// mongodb stores in bson which is little endian except for timestamp and counter
	pageID = binary.LittleEndian.Uint32(decodedURL)

	return
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
