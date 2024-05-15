package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/ceffu-fe-backend/entity"
	"github.com/mapprotocol/ceffu-fe-backend/logic"
	"github.com/mapprotocol/ceffu-fe-backend/resp"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
)

func SupportedChains(c *gin.Context) {
	req := &entity.SupportedChainsRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	page, size := utils.ValidatePage(req.Page, req.Size)

	list, count, code := logic.SupportedChains(page, size)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.SuccessList(c, page, count, list)
}
