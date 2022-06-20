package database

import (
	"context"
	"errors"
	"fmt"
	"linkshare_api/graph/model"
	"linkshare_api/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IllegalURLError string

func (e IllegalURLError) Error() string {
	return "URL input must be 1-30 alphanumeric or '_','-' characters"
}

type URLTakenError string

func (e URLTakenError) Error() string {
	return "URL is taken: " + string(e)
}

type LinkShareDB struct {
	client    *mongo.Client
	db        *mongo.Database
	pages     *mongo.Collection
	InsertOne func(ctx context.Context, document interface{},
		opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

func NewLinkShareDB(ctx context.Context) (linksDB *LinkShareDB, err error) {
	linksDB = new(LinkShareDB)
	client, err := GetClient(ctx)
	if err != nil {
		err = fmt.Errorf("no client: %w", err)
		return
	}
	linksDB.client = client

	db, err := GetDatabase(client)
	if err != nil {
		err = fmt.Errorf("no db: %w", err)
		return
	}
	linksDB.db = db
	linksDB.pages = db.Collection("pages")

	linksDB.InsertOne = linksDB.pages.InsertOne

	return
}
func (linksDB *LinkShareDB) Disconnect(ctx context.Context) {
	linksDB.client.Disconnect(ctx)
}

func GetClient(ctx context.Context) (client *mongo.Client, err error) {
	connectionUrl := utils.GetConf().ConnectionURL
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	return mongoClient, err
}

func GetDatabase(client *mongo.Client) (database *mongo.Database, err error) {
	database = client.Database(utils.GetConf().DBName)
	return
}

// This func doesn't validate if page is valid base64 encoding
func (linksDB *LinkShareDB) CreateNewPage(ctx context.Context, URL string, userID string) (createdPage *model.Page, err error) {
	createdURL := ""
	// If the user did not input a custom URL, create a random one
	isCustomURL := len(URL) != 0
	if isCustomURL {
		createdURL = URL
		if !utils.IsValidURL(URL) {
			err = IllegalURLError("")
			return
		}
	} else {
		createdURL = utils.GetRandomURL(6)
	}

	updateMap := bson.M{
		"_id":         createdURL,
		"dateAdded":   primitive.NewDateTimeFromTime(time.Now()),
		"description": "",
		"title":       "",
		"user_id":     userID,
		"links":       []string{},
		"schema":      utils.GetConf().SchemaVersion,
	}

	_, err = linksDB.InsertOne(ctx, updateMap)

	// Custom URLs don't need retries, we fail if it's taken
	if err != nil && isCustomURL {
		if strings.Contains(err.Error(), "E11000") {
			return nil, URLTakenError(URL)
		} else {
			return nil, err
		}
	}

	for misses := 0; err != nil && misses < 2; misses += 1 {
		// E11000 corresponse to DuplicateKey error
		// https://github.com/mongodb/mongo/blob/master/src/mongo/base/error_codes.yml
		if !strings.Contains(err.Error(), "E11000") {
			return nil, fmt.Errorf("failed to create new page: %w", err)
		}
		createdURL = utils.GetRandomURL(6)
		updateMap["_id"] = createdURL
		_, err = linksDB.InsertOne(ctx, updateMap)
		if err != nil && misses == 1 {
			err = errors.New("congrats you've won the lottery")
		}
	}

	createdPage = new(model.Page)
	createdPage.User = userID
	createdPage.URL = createdURL
	return
}
