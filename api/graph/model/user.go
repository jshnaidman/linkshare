package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  *string            `json:"username"`
	FirstName *string            `json:"firstName"`
	LastName  *string            `json:"LastName"`
	Email     *string            `json:"email"`
	GoogleID  *string            `json:"googleID"`
	PageURLs  []string           `json:"PageURLs"`
	Schema    int                `json:"-"` // omitted from graphql
}
