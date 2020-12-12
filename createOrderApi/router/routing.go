package router

import (
	adapter "Api/adapters"

	"Api/helper/redis"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	client := router.Group("/api")
	{
		client.POST("/login", redis.Login)

		client.POST("/create", adapter.CreateOrder)
	}
	return router
}
