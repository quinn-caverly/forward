package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

var stussyCollection *mongo.Collection

func CreateConnToMongo() (error) {
    c, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://mongo-service:27017"))
	if err != nil {
		return fmt.Errorf("Was not able to connect to mongo via the service %w", err)
	}
    client = c
    return nil
}

func CloseConnToMongo() (error) {
    return client.Disconnect(context.TODO())
}

func GetStussyColl() (*mongo.Collection) {
    if stussyCollection != nil {
        return stussyCollection
    }
    stussyCollection = client.Database("products").Collection("stussy")
    return stussyCollection
}
