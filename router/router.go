package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/ceffu-fe-backend/controller"
)

func Register(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")
	v1.GET("/health", controller.Health)
}
