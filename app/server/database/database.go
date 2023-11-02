package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/lucas-kern/tower-of-babel_server/app/model"

	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

// database is used to manage database connections
// has important functions for connecting and using the database

// Database represents the current Database connection
type Database struct {
	client 	*mongo.Client
	databaseName string
	url     string
}

// Datastore represents a store for the data (Database session)
type Datastore struct {
	client *mongo.Client
	name    string
}

// Connect will connect to a mongodb client 
// returns a Database struct or an error
func Connect() (*Database, error){
	// located in the env file
	mongoHost := os.Getenv("MONGODB_URL")
	databaseName := os.Getenv("MONGODB_NAME")

	fmt.Print("Connecting to Host: ", mongoHost + "\n")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoHost))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &Database{client: client, databaseName: databaseName}, nil
}

// Close will disconnect the database client connection
func (d *Database) Close() {
	if d.client != nil {
		d.client.Disconnect(context.TODO())
	}
	d.client = nil
	log.Println("Database closed")
}

//	GetCollection gets a collection from the mongo database with name c
//	returns the collection
func (d *Database) GetCollection(c string) model.Collection{
	log.Println("From Datastore: Getting collection" + c)
	return GetMongoCollection(d.client.Database(d.databaseName).Collection(c))
}

//	GetUsers gets the user collection from the mongo database with name c
//	returns the users collection
func (d *Database) GetUsers() model.Collection{
	log.Println("Retrieving Users collection")
	return GetMongoCollection(d.client.Database(d.databaseName).Collection("users"))
}

//	GetBases gets the base collection from the mongo database with name c
//	returns the bases collection
func (d *Database) GetBases() model.Collection{
	log.Println("Retrieving Bases collection")
	return GetMongoCollection(d.client.Database(d.databaseName).Collection("bases"))
}
