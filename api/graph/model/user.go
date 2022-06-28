package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  *string            `json:"username" bson:"username,omitempty"`
	FirstName *string            `json:"firstName" bson:"first_name,omitempty"`
	LastName  *string            `json:"lastName" bson:"last_name,omitempty"`
	Email     *string            `json:"email" bson:"email,omitempty"`
	GoogleID  *string            `json:"googleID" bson:"google_id,omitempty"`
	PageURLs  []string           `json:"pageURLs" bson:"pages"`
	Schema    int                `json:"-" bson:"schema"` // omitted from graphql
}
