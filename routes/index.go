package routes

import (
	"github.com/20-VIGNESH-K/crud_operations/services"
	"github.com/gin-gonic/gin"
)

func ProfileRoute(router *gin.Engine) {
	router.Static("/createProfile", "./frontend/createProfile/")
	router.Static("/getAllProfile", "./frontend/getAllProfile/")
	router.Static("/getProfileByName", "./frontend/getProfileByName/")
	router.Static("/updateProfileByName", "./frontend/updateProfileByName/")
	router.Static("/deleteProfile", "./frontend/deleteProfile/")
	router.Static("/home", "./frontend/index/")

	router.POST("/create", services.Create)
	router.POST("/createMany", services.CreateMany)
	router.POST("/update/:name", services.Update)
	router.DELETE("/delete/:name", services.Delete)
	router.GET("/getAll", services.GetAll)
	router.GET("/getUserByName/:name", services.GetUser)
}
