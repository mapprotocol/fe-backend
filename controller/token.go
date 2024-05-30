package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/ceffu-fe-backend/entity"
	"github.com/mapprotocol/ceffu-fe-backend/logic"
	"github.com/mapprotocol/ceffu-fe-backend/resp"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
)

func SupportedTokens(c *gin.Context) {
	req := &entity.SupportedTokensRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	page, size := utils.ValidatePage(req.Page, req.Size)

	list, count, code := logic.SupportedTokens(req.ChainID, req.Symbol, page, size)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.SuccessList(c, count, list)
}
