package database

import (
	"context"
	"errors"
	"fmt"
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

// TODO: Should probably change all CRUD operations to be model.Page methods.

// This func doesn't validate if page is valid base64 encoding
// TODO: Need to add created pageID to owning user
func (linksDB *LinkShareDB) CreatePage(ctx context.Context, URL string, userID primitive.ObjectID,
	insertOnePage utils.InsertOneFunc) (createdPage *model.Page, err error) {
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
		"date_added":  utils.DateTimeNow(),
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

// db.sessions.aggregate([{$match: {_id: 'VDMYBF72TWDR6SONKKX4M2FCAHEZT57QYELW22UUKQR7FD45H2SQ'}}, {$lookup: {from: "users", localField: "user_id", foreignField:"_id", as: "user"}}, {$unwind: "$user"}, {$replaceRoot: {newRoot: "$user"}}]).explain()
// db.sessions.aggregate([  {$replaceRoot: {newRoot: "$user"}}]).explain()
func FindUserForSession(ctx context.Context, sessionID string, sessionAggregate utils.AggregateFunc) (user *model.User, err error) {

	user = &model.User{}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{
			Key: "_id", Value: sessionID}},
		}},
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"}, {Key: "localField", Value: "user_id"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "user"},
		}}},
		bson.D{{Key: "$unwind", Value: "$user"}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: "$user"},
		},
		}},
	}

	cursor, err := sessionAggregate(ctx, pipeline)
	if err != nil {
		return
	}

	var results []model.User
	if err = cursor.All(context.TODO(), &results); err != nil {
		utils.LogError("findUserForSession - %s", err)
		return nil, err
	}
	if len(results) > 1 {
		err = fmt.Errorf("got multiple users for session key: %v", results)
	}
	if len(results) == 1 {
		user = &results[0]
	}
	return
}

// func (linksDB *LinkShareDB) FindUserFromSession(ctx context.Context, sessionID string,
// 	aggregationFunc AggregateFunc) (user *model.User, err error) {
// 	user = &model.User{}

// 	pipeline := mongo.Pipeline{
// 		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: 1}}}},
// 	}

// 	// err = findOneSession(ctx, bson.M{"_id": sessionID}).Decode(session)
// 	return
// }
