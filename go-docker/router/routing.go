package router

import (
	adapter "github.com/callicoder/go-docker/adapters"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api")
	{
		v1.POST("/create", adapter.CreateOrder)
		v1.GET("/query", adapter.QueryOrder)
		v1.POST("/refund", adapter.Refund)
		v1.GET("/queryrefund", adapter.QueryRefund)
		//v1.GET("/getlistbank",adapter.GetListBank)
		v1.POST("/callback", adapter.HandleCallback)
	}
	return router
}
