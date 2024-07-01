package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/utils"
	"strconv"
)

const (
	exactIn  = "exactIn"
	exactOut = "exactOut"
)

func GetRoute(c *gin.Context) {
	req := &entity.RouteRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}
	if utils.IsEmpty(req.FromChainID) {
		resp.ParameterErr(c, "missing fromChainId")
		return
	}
	if utils.IsEmpty(req.ToChainID) {
		resp.ParameterErr(c, "missing toChainId")
		return
	}
	if utils.IsEmpty(req.Amount) {
		resp.ParameterErr(c, "missing amount")
		return
	}
	if utils.IsEmpty(req.TokenInAddress) {
		resp.ParameterErr(c, "missing tokenInAddress")
		return
	}
	if utils.IsEmpty(req.TokenOutAddress) {
		resp.ParameterErr(c, "missing tokenOutAddress")
		return
	}
	if utils.IsEmpty(req.Kind) {
		resp.ParameterErr(c, "missing type")
		return
	}
	if req.Kind != exactIn && req.Kind != exactOut {
		resp.ParameterErr(c, "type must be exactIn or exactOut")
		return
	}
	if utils.IsEmpty(req.Slippage) {
		resp.ParameterErr(c, "missing slippage")
	}
	slippage, err := strconv.ParseUint(req.Slippage, 10, 64)
	if err != nil {
		resp.ParameterErr(c, "invalid slippage")
		return
	}
	if slippage < 0 || slippage > 5000 {
		resp.ParameterErr(c, "invalid slippage")
		return
	}

	ret, code := logic.GetRoute(req)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}

func Swap(c *gin.Context) {
	req := &entity.SwapRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}
	if utils.IsEmpty(req.Hash) {
		resp.ParameterErr(c, "missing hash")
		return
	}
	if utils.IsEmpty(req.Slippage) {
		resp.ParameterErr(c, "missing slippage")
	}
	slippage, err := strconv.ParseUint(req.Slippage, 10, 64)
	if err != nil {
		resp.ParameterErr(c, "invalid slippage")
		return
	}
	if slippage < 0 || slippage > 5000 {
		resp.ParameterErr(c, "invalid slippage")
		return
	}
	if utils.IsEmpty(req.From) {
		resp.ParameterErr(c, "missing from")
		return
	}
	if utils.IsEmpty(req.Receiver) {
		resp.ParameterErr(c, "missing receiver")
		return
	}

	ret, code := logic.Swap(req.Hash, req.Slippage, req.From, req.Receiver)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}
