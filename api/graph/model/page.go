package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Page struct {
	URL          string             `json:"URL" bson:"_id,omitempty"`
	Description  *string            `json:"description"`
	Title        *string            `json:"title"`
	Links        []*string          `json:"links"`
	OwningUserID primitive.ObjectID `json:"owningUserID" bson:"user_id"`
	Schema       int                `json:"-"`                   // omitted from graphql
	DateAdded    primitive.DateTime `json:"-" bson:"date_added"` // omitted from graphql
}
