package main

import (
	"context"
	"fmt"
	"log"

	"github.com/20-VIGNESH-K/crud_operations/config"
	"github.com/20-VIGNESH-K/crud_operations/routes"
	"github.com/20-VIGNESH-K/crud_operations/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoclient *mongo.Client
	ctx    context.Context
	server *gin.Engine
)


func initApp(mongoClient *mongo.Client) {
	routes.ProfileRoute(server)

}
func main() {
	server = gin.Default()
	mongoclient, err := config.ConnectDataBase()
	services.MongoClient = mongoclient
	services.Ctx = context.TODO()
	defer mongoclient.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	initApp(mongoclient)
	fmt.Println("server running on port", config.Port)
	log.Fatal(server.Run(config.Port))

}
