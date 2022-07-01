package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertOneFunc func(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
type FindOneAndUpdateFunc func(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult
type FindOneFunc func(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult
type AggregateFunc func(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*mongo.Cursor, error)
type DeleteOneFunc func(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
type UpdateByIDFunc func(ctx context.Context, id interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
