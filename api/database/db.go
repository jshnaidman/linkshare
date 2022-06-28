package database

import (
	"context"
	"errors"
	"fmt"
	"linkshare_api/contextual"
	"linkshare_api/graph/model"
	"linkshare_api/utils"
	"strings"

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
	Client   *mongo.Client
	DB       *mongo.Database
	Pages    *mongo.Collection
	Users    *mongo.Collection
	Sessions *mongo.Collection
}

type InsertOneFunc func(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
type FindOneAndUpdateFunc func(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult

func NewLinkShareDB(ctx context.Context) (linksDB *LinkShareDB, err error) {
	linksDB = new(LinkShareDB)
	client, err := GetClient(ctx)
	if err != nil {
		err = fmt.Errorf("no client: %w", err)
		return
	}
	linksDB.Client = client

	db, err := GetDatabase(client)
	if err != nil {
		err = fmt.Errorf("no db: %w", err)
		return
	}
	linksDB.DB = db
	linksDB.Pages = db.Collection("pages")
	linksDB.Users = db.Collection("users")
	linksDB.Sessions = db.Collection("sessions")

	return
}
func (linksDB *LinkShareDB) Disconnect(ctx context.Context) {
	linksDB.Client.Disconnect(ctx)
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
// TODO: Need to add created pageID to owning user
func (linksDB *LinkShareDB) CreatePage(ctx context.Context, URL string, userID primitive.ObjectID,
	insertOnePage InsertOneFunc) (createdPage *model.Page, err error) {
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

	insertMap := bson.M{
		"_id":         createdURL,
		"dateAdded":   utils.DateTimeNow(),
		"description": "",
		"title":       "",
		"user_id":     userID,
		"links":       []string{},
		"schema":      utils.GetConf().SchemaVersion,
	}

	_, err = insertOnePage(ctx, insertMap)

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
		insertMap["_id"] = createdURL
		_, err = insertOnePage(ctx, insertMap)
		if err != nil && misses == 1 {
			err = errors.New("congrats you've won the lottery")
		}
	}

	createdPage = new(model.Page)
	createdPage.OwningUserID = userID
	createdPage.URL = createdURL
	return
}

func (linksDB *LinkShareDB) UpsertUserByGoogleID(ctx context.Context, user *model.User,
	findOneUserAndUpdate FindOneAndUpdateFunc) (updatedUser *model.User, err error) {

	updateOption := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	updateData := bson.M{"$set": *user}

	filter := bson.M{"google_id": user.GoogleID}

	updatedUser = &model.User{}

	// updated user will have Id
	err = findOneUserAndUpdate(context.TODO(), filter, updateData, updateOption).Decode(updatedUser)

	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return updatedUser, err
}

// func (linksDB *LinkShareDB) FindUser(user *model.User) {

// }

func (linksDB *LinkShareDB) CreateSession(ctx context.Context, session *contextual.Session,
	insertOne InsertOneFunc) (err error) {

	_, err = insertOne(ctx, *session)

	return
}
