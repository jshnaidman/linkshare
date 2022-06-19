package database

import (
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"linkshare_api/conf"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var urlEncoding b64.Encoding = *b64.URLEncoding.WithPadding(b64.NoPadding)

// taken from base64.encodeURL
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

type IllegalURLError string

func (e IllegalURLError) Error() string {
	return "URL input is not allowed: " + string(e)
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

func NewLinkShareDB() (linksDB *LinkShareDB, err error) {
	linksDB = new(LinkShareDB)
	client, err := GetClient()
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
func (linksDB *LinkShareDB) Disconnect() {
	linksDB.client.Disconnect(context.TODO())
}

// This func doesn't validate if page is valid base64 encoding
func (linksDB *LinkShareDB) CreateNewPage(URL string, userID primitive.ObjectID) (createdURL string, err error) {
	// If the user did not input a custom URL, create a random one
	isCustomURL := len(URL) != 0
	if isCustomURL {
		createdURL = URL
	} else {
		createdURL = GetRandomURL()
	}

	updateMap := bson.M{
		"_id":         createdURL,
		"dateAdded":   primitive.NewDateTimeFromTime(time.Now()),
		"description": "",
		"title":       "",
		"user_id":     userID,
		"links":       []string{},
		"schema":      conf.GetConf().SchemaVersion,
	}

	_, err = linksDB.InsertOne(context.TODO(), updateMap)

	// Custom URLs don't need retries, we fail if it's taken
	if err != nil && isCustomURL {
		if strings.Contains(err.Error(), "E11000") {
			return "", URLTakenError(URL)
		} else {
			return "", err
		}
	}

	for misses := 0; err != nil && misses < 2; misses += 1 {
		// E11000 corresponse to DuplicateKey error
		// https://github.com/mongodb/mongo/blob/master/src/mongo/base/error_codes.yml
		if !strings.Contains(err.Error(), "E11000") {
			return "", fmt.Errorf("failed to create new page: %w", err)
		}
		createdURL = GetRandomURL()
		updateMap["_id"] = createdURL
		_, err = linksDB.InsertOne(context.TODO(), updateMap)
		if err != nil && misses == 1 {
			err = errors.New("congrats you've won the lottery")
		}
	}

	return
}

func GetRandomURL() string {
	// generate a random 6 character string
	sb := strings.Builder{}
	sb.Grow(6)
	for i := 0; i < 6; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
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
