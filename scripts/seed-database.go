package main

import (
	"context"
	"fmt"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"github.com/lucas-kern/tower-of-babel_server/app/model"
)


func main() {
	localhost := "localhost:27017"
	mongohost := "mongodb://" + localhost

	fmt.Print("Connecting to Localhost: ", localhost+ "\n")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongohost))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Print(cancel)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
			log.Fatal(err)
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{"name": "tower-of-babel"})
	if err != nil {
			log.Fatal(err)
	}
	fmt.Println(databases)
	var collection *mongo.Collection
	if len(databases) == 1 {
		log.Print("database is already created")
	} else {
		log.Print("Creating database")
	}

	sample_bases := []model.Base{
		{
			Owner: "Kern",
			Tower: model.Building{
				PosX: 3,
				PosY: 5,
				PosZ: 8 },
			ArmyCamp: model.Building{
				PosX: 2,
				PosY: 5,
				PosZ: 1 },
			Barracks: model.Building{
				PosX: 1,
				PosY: 7,
				PosZ: 5 }},
		{
				Owner: "Coffee",
				Tower: model.Building{
					PosX: 3,
					PosY: 5,
					PosZ: 8 },
				ArmyCamp: model.Building{
					PosX: 2,
					PosY: 5,
					PosZ: 1 },
				Barracks: model.Building{
					PosX: 1,
					PosY: 7,
					PosZ: 5 }},
			{
				Owner: "Diego",
				Tower: model.Building{
					PosX: 3,
					PosY: 5,
					PosZ: 8 },
				ArmyCamp: model.Building{
					PosX: 2,
					PosY: 5,
					PosZ: 1 },
				Barracks: model.Building{
					PosX: 1,
					PosY: 7,
					PosZ: 5 }}}

	collection = client.Database("tower-of-babel").Collection("users")
	client.ListDatabaseNames(ctx, bson.M{"name": "tower-of-babel"})

	b := make([]interface{}, len(sample_bases))
	for i := range sample_bases {
			b[i] = sample_bases[i]
	}

	_, err = collection.InsertMany(context.Background(), b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(collection)
	defer cancel()
	defer client.Disconnect(ctx) 
}