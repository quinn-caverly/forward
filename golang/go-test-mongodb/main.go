package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://mongo-service:27017"))
	if err != nil {
		log.Fatal("Was not able to connect to the mongodb via the service, ", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("products").Collection("stussy")

	document := bson.D{
		{Key: "name", Value: "John Doe"},
		{Key: "age", Value: 30},
	}

	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
        log.Fatal("Error when attempting write to db, ", err)
	}

    log.Println("Success?")
}
