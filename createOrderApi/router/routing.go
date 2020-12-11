package router

import (
	"Api/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	client := router.Group("/api")
	{

		client.POST("/create", controller.CreateOrder)
	}
	return router
}
