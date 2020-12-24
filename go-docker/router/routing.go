package router

import (
	adapter "github.com/callicoder/go-docker/adapters"

	"github.com/callicoder/go-docker/helper/redis"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api")
	{
	
		v1.POST("/login", redis.Login)
		v1.POST("/create", adapter.CreateOrder)
		v1.POST("/query",adapter.QueryOrder)
		v1.POST("/refund",adapter.Refund)
		v1.POST("/queryrefund",adapter.QueryRefund)
		v1.GET("/getlistbank",adapter.GetListBank)
	}
	return router
}
