package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"zalopay-api/autodebit"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	v1 := router.Group("/zalopay-api/api")
	{
		v1.POST("/create-binding", autodebit.CreateBinding)
		v1.POST("/create-order", autodebit.CreateOrder)
		v1.POST("/query-order", autodebit.QueryOrder)
		v1.GET("/get-binding/:binding_token", autodebit.GetBinding)
		v1.POST("/query-binding", autodebit.QueryBinding)
		v1.POST("/check-balance", autodebit.CheckBalance)
		v1.POST("/refund-order", autodebit.RefundOrder)
		v1.POST("/query-refund", autodebit.QueryRefund)
		v1.POST("/pay", autodebit.Pay)
		v1.POST("/un-bind", autodebit.Unbinding)

	}
	return router
}
