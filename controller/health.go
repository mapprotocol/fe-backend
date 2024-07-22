package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/resp"
)

func Health(c *gin.Context) {
	resp.SuccessNil(c)
}
