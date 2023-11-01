package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/20-VIGNESH-K/crud_operations/config"
	"github.com/20-VIGNESH-K/crud_operations/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Ctx         context.Context
	MongoClient *mongo.Client
)

// CustomValidator is a custom validation function
//
//	func CustomValidator(fl validator.FieldLevel) bool {
//		// Use a regular expression to check if the field contains only letters
//		name := fl.Field().String()
//		match, _ := regexp.MatchString("^[A-Za-z]+$", name)
//		return match
//	}
func CustomValidator(fl validator.FieldLevel) bool {
	// Use a regular expression to check if the field contains only letters and spaces
	name := fl.Field().String()
	// Remove leading and trailing spaces before validation
	trimmedName := strings.TrimSpace(name)
	match, _ := regexp.MatchString("^[A-Za-z ]+$", trimmedName)
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

	var user models.Profile
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}
	filter := bson.M{"name": user.Name}
	existingProfile := config.ProfileCollection.FindOne(Ctx, filter)
	if existingProfile.Err() == nil {
		context.JSON(http.StatusConflict, gin.H{"message": "Profile with the same name already exists"})
		return // Return early if the profile already exists
	}
	fmt.Println(user)
	err := CheckValidation(user)
	if err == false {
		context.JSON(http.StatusNotAcceptable, gin.H{"message": "Validation error"})
	} else {

		create, err := config.ProfileCollection.InsertOne(Ctx, &user)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(create.InsertedID)
		context.JSON(http.StatusOK, gin.H{"message": "success"})
	}

}

func CreateMany(context *gin.Context) {

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
	_, err := config.ProfileCollection.InsertMany(Ctx, profilesToInsert)
	if err != nil {
		log.Println("Failed to insert profiles:", err)
		return
	}

}

func Update(context *gin.Context) {

	name := context.Param("name")

	filter := bson.M{"name": name}
	var user models.Profile

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, err.Error())
	}

	existingProfile := config.ProfileCollection.FindOne(Ctx, filter)
	if existingProfile.Err() != nil {
		context.JSON(http.StatusConflict, gin.H{"message": "Profile name not exists"})
		return // Return early if the profile already exists
	}

	result, err := config.ProfileCollection.UpdateOne(Ctx, filter, bson.M{"$set": &user})
	if result.ModifiedCount == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"message": "You Updated nothing!!!..."})
		return
	}

	if err != nil {
		fmt.Println("Error")
		log.Fatal(err)
	}
	context.JSON(http.StatusCreated, gin.H{"message": "success"})

}

func Delete(context *gin.Context) {

	name := context.Param("name")
	filter := bson.M{"name": name}
	var user models.Profile
	err := config.ProfileCollection.FindOne(Ctx, filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		// User with the specified username does not exist
		context.JSON(http.StatusNotFound, gin.H{"message": "Name not found!!!. Enter correct Profile name to delete."})
	} else {

		_, err := config.ProfileCollection.DeleteOne(Ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		context.JSON(http.StatusAccepted, gin.H{"message": "success"})
	}
}

func GetAll(context *gin.Context) {

	var users []*models.Profile
	cursor, err := config.ProfileCollection.Find(Ctx, bson.D{{}})
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

	var users []*models.Profile
	name := context.Param("name")
	filter := bson.M{"name": name}
	cursor, err := config.ProfileCollection.Find(Ctx, filter)
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

func GetAllProfilesSortedByName(c *gin.Context) {

	// Set the sorting options by name.
	options := options.Find().SetSort(bson.D{{Key: "name", Value: 1}}) // 1 for ascending, -1 for descending

	// Perform the database query.
	cursor, err := config.ProfileCollection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching items"})
		return
	}

	// Store the results in a slice.
	var items []models.Profile
	if err := cursor.All(context.TODO(), &items); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding items"})
		return
	}

	// Optionally, you can sort the items in Go for display purposes.
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	// Return the sorted items.
	c.JSON(http.StatusOK, items)
}
