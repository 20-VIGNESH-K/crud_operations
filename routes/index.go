package routes

import (
	"github.com/20-VIGNESH-K/crud_operations/services"
	"github.com/gin-gonic/gin"
	
)

func ProfileRoute(router *gin.Engine) {
	router.POST("/create",services.Create)

}
