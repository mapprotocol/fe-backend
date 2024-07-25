package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/constants"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/utils"
	"math/big"
	"strconv"
)

const (
	exactIn  = "exactIn"
	exactOut = "exactOut"
)

func Route(c *gin.Context) {
	req := &entity.RouteRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}
	//if utils.IsEmpty(req.FromChainID) || req.FromChainID == "0" {
	//	resp.ParameterErr(c, "missing fromChainId")
	//	return
	//}
	//if _, ok := new(big.Int).SetString(req.FromChainID, 10); !ok {
	//	resp.ParameterErr(c, "invalid fromChainId")
	//	return
	//}

	//if utils.IsEmpty(req.ToChainID) || req.ToChainID == "0" {
	//	resp.ParameterErr(c, "missing toChainId")
	//	return
	//}
	//if _, ok := new(big.Int).SetString(req.ToChainID, 10); !ok {
	//	resp.ParameterErr(c, "invalid toChainId")
	//	return
	//}
	if utils.IsEmpty(req.Amount) || req.Amount == "0" {
		resp.ParameterErr(c, "missing amount")
		return
	}
	if _, err := strconv.ParseFloat(req.Amount, 64); err != nil {
		resp.ParameterErr(c, "invalid amount")
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
	if utils.IsEmpty(req.Type) {
		resp.ParameterErr(c, "missing type")
		return
	}
	if req.Type != exactIn && req.Type != exactOut {
		resp.ParameterErr(c, "type must be exactIn or exactOut")
		return
	}
	if utils.IsEmpty(req.Slippage) || req.Slippage == "0" {
		resp.ParameterErr(c, "missing slippage")
	}
	slippage, err := strconv.ParseUint(req.Slippage, 10, 64)
	if err != nil {
		resp.ParameterErr(c, "invalid slippage")
		return
	}
	if slippage < constants.SlippageMin || slippage > constants.SlippageMax {
		resp.ParameterErr(c, "invalid slippage")
		return
	}
	if req.Action == dao.OrderActionToEVM {
		if req.FromChainID != constants.BTCChainID && req.FromChainID != constants.TONChainID {
			resp.ParameterErr(c, "invalid fromChainId")
			return
		}
	} else if req.Action == dao.OrderActionFromEVM {
		if req.ToChainID != constants.BTCChainID && req.ToChainID != constants.TONChainID {
			resp.ParameterErr(c, "invalid toChainId")
			return
		}
	} else {
		resp.ParameterErr(c, "invalid action")
		return
	}

	code := resp.CodeSuccess
	ret := make([]*entity.RouteResponse, 0)

	switch req.Action {
	case dao.OrderActionToEVM:
		if req.FromChainID == constants.TONChainID {
			ret, code = logic.GetTONToEVMRoute(req)
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		}
	case dao.OrderActionFromEVM:
		if req.ToChainID == constants.TONChainID {
			ret, code = logic.GetEVMToTONRoute(req, slippage)
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		}
	}

	resp.SuccessList(c, int64(len(ret)), ret)
}

func Swap(c *gin.Context) {
	req := &entity.SwapRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}
	if utils.IsEmpty(req.SrcChain) {
		resp.ParameterErr(c, "missing srcChain")
		return
	}
	if _, ok := new(big.Int).SetString(req.SrcChain, 10); !ok {
		resp.ParameterErr(c, "invalid srcChain")
		return
	}

	if !utils.IsValidEvmAddress(req.SrcToken) {
		resp.ParameterErr(c, "invalid srcToken")
		return
	}
	//if !utils.IsValidEvmAddress(req.DstToken) {
	//	resp.ParameterErr(c, "invalid dstToken")
	//	return
	//} // todo
	if utils.IsEmpty(req.Amount) {
		resp.ParameterErr(c, "missing amount")
		return
	}
	if _, err := strconv.ParseFloat(req.Amount, 64); err != nil {
		resp.ParameterErr(c, "invalid amount")
		return
	}
	if req.Decimal <= 0 {
		resp.ParameterErr(c, "missing decimal")
		return
	}

	if utils.IsEmpty(req.DstChain) {
		resp.ParameterErr(c, "missing dstChain")
		return
	}
	if _, ok := new(big.Int).SetString(req.DstChain, 10); !ok {
		resp.ParameterErr(c, "invalid dstChain")
		return
	}
	//if !utils.IsValidBitcoinAddress(req.Receiver, logic.NetParams) {
	//	resp.ParameterErr(c, "invalid receiver")
	//	return
	//}

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
	if slippage < constants.SlippageMin || slippage > constants.SlippageMax {
		resp.ParameterErr(c, "invalid slippage")
		return
	}

	amountBigFloat, ok := new(big.Float).SetString(req.Amount)
	if !ok {
		resp.ParameterErr(c, "invalid amount")
		return
	}
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(req.Decimal)), nil)
	amount := new(big.Float).Mul(amountBigFloat, new(big.Float).SetInt(exp))
	amountBigInt, _ := amount.Int(nil)

	ret, code := logic.Swap(req.SrcChain, req.SrcToken, req.Sender, amountBigInt, req.DstChain, req.Receiver, req.Hash, slippage)
	if code != resp.CodeSuccess {
		resp.Error(c, code)
		return
	}
	resp.Success(c, ret)
}
