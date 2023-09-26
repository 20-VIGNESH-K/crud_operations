package routes

import (
	"github.com/20-VIGNESH-K/crud_operations/services"
	"github.com/gin-gonic/gin"
)

func ProfileRoute(router *gin.Engine) {
	router.POST("/create", services.Create)
	router.POST("/update/:name", services.Update)
	router.DELETE("/delete/:name", services.Delete)
	router.GET("/getAll", services.GetAll)
}
