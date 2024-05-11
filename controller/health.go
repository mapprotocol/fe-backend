package controller

import (
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {

	resp.SuccessList(c, page, count, list)
}
