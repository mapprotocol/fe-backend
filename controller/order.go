package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/ceffu-fe-backend/entity"
	"github.com/mapprotocol/ceffu-fe-backend/logic"
	"github.com/mapprotocol/ceffu-fe-backend/resp"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
)

func OrderList(c *gin.Context) {
	req := &entity.OrderListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	page, size := utils.ValidatePage(req.Page, req.Size)

	list, count, code := logic.OrderList(req.Sender, page, size)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.SuccessList(c, count, list)
}

func OrderDetail(c *gin.Context) {
	req := &entity.OrderDetailRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	ret, code := logic.OrderDetail(req.OrderID)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}
