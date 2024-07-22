package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/controller"
)

func Register(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	v1.GET("/health", controller.Health)
	v1.GET("/route", controller.Route)
	v1.GET("/swap", controller.Swap)
	v1.GET("/order/create", controller.CreateOrder)
	v1.GET("/order/update", controller.UpdateOrder)
	v1.GET("/order/list", controller.OrderList)
	v1.GET("/order/detail", controller.OrderDetail)
}
