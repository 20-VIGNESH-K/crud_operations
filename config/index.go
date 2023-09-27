package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const ConnectionString = "mongodb://localhost:27017"
const Port = ":8080"
const DatabaseName = "Profile"
var ProfileCollection *mongo.Collection



func ConnectDataBase() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoConnection := options.Client().ApplyURI(ConnectionString)
	mongoClient, err := mongo.Connect(ctx, mongoConnection)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	ProfileCollection = mongoClient.Database(DatabaseName).Collection("ProfileCollections")
	fmt.Println("Database Connected")
	return mongoClient, nil
}
