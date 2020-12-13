package router

import (
	adapter "github.com/callicoder/go-docker/adapters"

	"github.com/callicoder/go-docker/helper/redis"

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
