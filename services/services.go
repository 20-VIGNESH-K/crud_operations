package services

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/20-VIGNESH-K/crud_operations/config"
	"github.com/20-VIGNESH-K/crud_operations/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Ctx         context.Context
	MongoClient *mongo.Client
)

func Create(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")

	var user models.Profile
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}
	create, err := ProfileCollection.InsertOne(Ctx, &user)
	if err != nil {
		log.Fatal(err)
	}
	context.JSON(http.StatusCreated, gin.H{"status": "success", "data": create})

}

func Update(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")

	name := context.Param("name")

	filter := bson.M{"name": name}
	var user models.Profile
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}
	_, err := ProfileCollection.UpdateOne(Ctx, filter, bson.M{"$set": &user})

	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}
	context.JSON(http.StatusCreated, gin.H{"message": "success"})

}

func Delete(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")
	name := context.Param("name")
	filter := bson.M{"name": name}
	_, err := ProfileCollection.DeleteOne(Ctx, filter)
	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}
	context.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func GetAll(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")
	var users []*models.Profile
	cursor, err := ProfileCollection.Find(Ctx, bson.D{{}})
	if err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	for cursor.Next(Ctx) {
		var user models.Profile
		err := cursor.Decode(&user)
		if err != nil {
			context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	cursor.Close(Ctx)
	if len(users) == 0 {
		context.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	context.JSON(http.StatusOK, users)

}
