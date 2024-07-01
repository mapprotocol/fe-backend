package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/utils"
)

func CreateOrder(c *gin.Context) {
	req := &entity.CreateOrderRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	// sender address check
	if req.Action == dao.OrderActionToEVM {
		if !utils.IsValidBitcoinAddress(req.Sender, logic.NetParams) {
			resp.ParameterErr(c, "invalid sender")
			return
		}
		if !utils.IsValidEvmAddress(req.Receiver) {
			resp.ParameterErr(c, "invalid receiver")
			return
		}
	} else if req.Action == dao.OrderActionFromEVM {
		if !utils.IsValidEvmAddress(req.Sender) {
			resp.ParameterErr(c, "invalid sender")
			return
		}
		if !utils.IsValidBitcoinAddress(req.Receiver, logic.NetParams) {
			resp.ParameterErr(c, "invalid receiver")
			return
		}
	} else {
		resp.ParameterErr(c, "invalid action")
	}

	ret, code := logic.CreateOrder(req.SrcChain, req.SrcToken, req.Sender, req.Amount, req.DstChain, req.DstToken, req.Receiver, req.Action)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}

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
