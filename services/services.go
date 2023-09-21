package services

import (
	"context"
	"log"
	"net/http"

	"github.com/20-VIGNESH-K/crud_operations/config"
	"github.com/20-VIGNESH-K/crud_operations/models"
	"github.com/gin-gonic/gin"
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


