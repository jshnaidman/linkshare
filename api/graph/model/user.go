package model

import (
	"context"
	"errors"
	"linkshare_api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  *string            `json:"username" bson:"username,omitempty"`
	FirstName *string            `json:"firstName" bson:"first_name,omitempty"`
	LastName  *string            `json:"lastName" bson:"last_name,omitempty"`
	Email     *string            `json:"email" bson:"email,omitempty"`
	GoogleID  *string            `json:"googleID" bson:"google_id,omitempty"`
	PageURLs  []string           `json:"pageURLs" bson:"pages,omitempty"`
	Schema    int                `json:"-" bson:"schema"` // omitted from graphql
}

func (user *User) UpsertUserByGoogleID(ctx context.Context,
	findOneUserAndUpdate utils.FindOneAndUpdateFunc) (updatedUser *User, err error) {

	updateOption := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	updateData := bson.M{"$set": *user}

	filter := bson.M{"google_id": user.GoogleID}

	updatedUser = &User{}

	// updated user will have Id
	err = findOneUserAndUpdate(context.TODO(), filter, updateData, updateOption).Decode(updatedUser)

	if err == mongo.ErrNoDocuments {
		err = nil
	}

	return updatedUser, err
}

func (user *User) Update(ctx context.Context, updateUserByID utils.UpdateByIDFunc) (err error) {
	if user.ID.IsZero() {
		return errors.New("no userID on user")
	}
	_, err = updateUserByID(ctx, user.ID, bson.M{"$set": *user})
	return
}

func (user *User) LoadByUsername(ctx context.Context, findOneUser utils.FindOneFunc) (err error) {
	if len(*user.Username) == 0 {
		return errors.New("no username on user")
	}
	err = findOneUser(ctx, bson.M{
		"username": user.Username,
	}).Decode(user)

	return
}

func (user *User) LoadByID(ctx context.Context, findOneUser utils.FindOneFunc) (err error) {
	if user.ID.IsZero() {
		return errors.New("no userID on user")
	}
	err = findOneUser(ctx, bson.M{
		"_id": user.ID,
	}).Decode(user)

	return
}

func (user *User) PushPage(ctx context.Context, page string, updateUserByID utils.UpdateByIDFunc) (err error) {
	if user.ID.IsZero() {
		return errors.New("no userID on user")
	}
	_, err = updateUserByID(ctx, user.ID, bson.M{"$push": bson.M{
		"pages": page,
	}})
	return
}
func (user *User) DeleteByUsername(ctx context.Context, deleteOneUser utils.DeleteOneFunc) (err error) {
	if len(*user.Username) == 0 {
		return errors.New("no username on user")
	}
	_, err = deleteOneUser(ctx, bson.M{
		"username": user.Username,
	})
	return
}
