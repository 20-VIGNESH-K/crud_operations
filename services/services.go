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

func CheckValidation(user models.Profile) bool {
	validate := validator.New()
	validate.RegisterValidation("customValidator", CustomValidator)

	err := validate.Struct(user)
	if err != nil {
		// Handle validation errors
		for _, validationErr := range err.(validator.ValidationErrors) {
			log.Printf("Validation Error in field %s: %s\n", validationErr.Field(), validationErr.Tag())
			return false
		}
	}

	log.Println("Validation successful")
	return true

}

func Create(context *gin.Context) {
	var ProfileCollection = MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")
	var user models.Profile
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(user)
	err := CheckValidation(user)
	if err == false {
		context.JSON(http.StatusNotAcceptable, gin.H{"message": "Validation error"})
	} else {

		create, err := ProfileCollection.InsertOne(Ctx, &user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(create.InsertedID)
		context.JSON(http.StatusOK, gin.H{"message": "success"})
	}

}

func CreateMany(context *gin.Context) {
	ProfileCollection := MongoClient.Database(config.DatabaseName).Collection("ProfileCollections")
	var profiles []*models.Profile
	if err := context.ShouldBindJSON(&profiles); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(profiles)

	// Create an array of interface to store the profile documents
	var profilesToInsert []interface{}

	// Convert each profile to an interface and add to the array
	for _, profile := range profiles {
		err := CheckValidation(*profile)
		if err != false {
			profilesToInsert = append(profilesToInsert, profile)
			context.JSON(http.StatusOK, gin.H{"message": "inserted successfully"})
		} else {
			context.JSON(http.StatusNotAcceptable, gin.H{"message": "Validation error"})
			continue
		}
	}

	// Insert the profiles into the collection
	_, err := ProfileCollection.InsertMany(Ctx, profilesToInsert)
	if err != nil {
		log.Println("Failed to insert profiles:", err)
		return
	}

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
	var user models.Profile
	err := ProfileCollection.FindOne(Ctx, filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// User with the specified username does not exist
		context.JSON(http.StatusNotFound, gin.H{"message": "Name not found"})
	} else {

		_, err := ProfileCollection.DeleteOne(Ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		context.JSON(http.StatusAccepted, gin.H{"message": "success"})
	}
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
