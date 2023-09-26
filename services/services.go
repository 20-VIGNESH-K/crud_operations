package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/20-VIGNESH-K/crud_operations/config"
	"github.com/20-VIGNESH-K/crud_operations/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Ctx         context.Context
	MongoClient *mongo.Client
)

// CustomValidator is a custom validation function
func CustomValidator(fl validator.FieldLevel) bool {
	// Use a regular expression to check if the field contains only letters
	name := fl.Field().String()
	match, _ := regexp.MatchString("^[A-Za-z]+$", name)
	return match
}

func Create(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")
	var user models.Profile
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	validate.RegisterValidation("customValidator", CustomValidator)

	err := validate.Struct(user)
	if err != nil {
		// Handle validation errors
		for _, validationErr := range err.(validator.ValidationErrors) {
			log.Printf("Validation Error in field %s: %s\n", validationErr.Field(), validationErr.Tag())
			context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Validation Error!! Enter Properly..."})
			return
		}
	} else {
		log.Println("Validation successful")
		context.JSON(http.StatusOK, gin.H{"status": "Success", "message": "Validation Successfull"})
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
	context.JSON(http.StatusAccepted, gin.H{"message": "success"})
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

func GetUser(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")
	var users []*models.Profile
	name := context.Param("name")
	filter := bson.M{"name": name}
	cursor, err := ProfileCollection.Find(Ctx, filter)
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
	context.JSON(http.StatusOK, users)

}
