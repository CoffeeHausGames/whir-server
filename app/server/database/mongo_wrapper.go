package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This file is used to wrap MongoDB models 
// it allows to keep all Mongo objects in one place

// MongoCollection represents a mongogb implementation of [model.Collection]
type MongoCollection struct {
	Collection *mongo.Collection
}

// Makes a call to mongodb to count documents
func (c MongoCollection) CountDocuments(ctx context.Context, docs interface{}) (int64, error){
	return c.Collection.CountDocuments(ctx, docs)
}

// Makes a call to mongodb to update one document
func (c MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions)(*mongo.UpdateResult, error){
	return c.Collection.UpdateOne(ctx, filter, update, opts...)
}

// Makes a call to mongodb to insert one document
func (c MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error){
	result, err := c.Collection.InsertOne(ctx, document, opts...)
	// TODO check error and return if there is one here
	newID := result.InsertedID
	// If this ever changes in mongo to not be an objectID this will break
	return newID.(primitive.ObjectID), err
}

// Makes a call to mongodb to find one document
func (c MongoCollection) FindOne(doc interface{}, ctx context.Context, filter interface{},opts ...*options.FindOneOptions) error{
	err := c.Collection.FindOne(ctx, filter, opts...).Decode(doc)
	return err
}

// GetMongoCollection returns a mongo collection with the [name]
// Wraps a Mongo Collection with our implementation of it
func GetMongoCollection(mC *mongo.Collection) MongoCollection {
	c := MongoCollection{mC}

	return c
}
