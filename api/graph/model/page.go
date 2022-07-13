package model

import (
	"context"
	"errors"
	"linkshare_api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Page struct {
	URL          string             `json:"URL" bson:"_id,omitempty"`
	Description  *string            `json:"description"`
	Title        *string            `json:"title"`
	Links        []*string          `json:"links"`
	OwningUserID primitive.ObjectID `json:"owningUserID" bson:"user_id"`
	Schema       int                `json:"-"`                   // omitted from graphql
	DateAdded    primitive.DateTime `json:"-" bson:"date_added"` // omitted from graphql
}

func (page *Page) LoadByURL(ctx context.Context, findOnepage utils.FindOneFunc) (err error) {
	if len(page.URL) == 0 {
		return errors.New("no URL on page")
	}
	err = findOnepage(ctx, bson.M{
		"_id": page.URL,
	}).Decode(page)

	return
}
