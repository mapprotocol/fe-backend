package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/constants"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"strings"
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
	amountDecimal, err := decimal.NewFromString(req.Amount)
	if err != nil {
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
	feeRatio := uint64(0)
	if !utils.IsEmpty(req.FeeRatio) {
		var err error
		feeRatio, err = strconv.ParseUint(req.FeeRatio, 10, 64)
		if err != nil {
			resp.ParameterErr(c, "invalid feeRatio")
			return
		}
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

	msg := ""
	code := resp.CodeSuccess
	ret := make([]*entity.RouteResponse, 0)

	switch req.Action {
	case dao.OrderActionToEVM:
		if req.FromChainID == constants.TONChainID {
			ret, msg, code = logic.GetTONToEVMRoute(req, amountDecimal, feeRatio, slippage)
			if code == resp.CodeExternalServerError {
				resp.ExternalServerError(c, msg)
				return
			}
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		} else if req.FromChainID == constants.BTCChainID {
			ret, msg, code = logic.GetBitcoinToEVMRoute(req, amountDecimal, feeRatio, slippage)
			if code == resp.CodeExternalServerError {
				resp.ExternalServerError(c, msg)
				return
			}
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		} else {
			resp.ParameterErr(c, "invalid fromChainId")
			return
		}
	case dao.OrderActionFromEVM:
		if req.ToChainID == constants.TONChainID {
			ret, msg, code = logic.GetEVMToTONRoute(req, amountDecimal, feeRatio, slippage)
			if code == resp.CodeExternalServerError {
				resp.ExternalServerError(c, msg)
				return
			}
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		} else if req.ToChainID == constants.BTCChainID {
			ret, msg, code = logic.GetEVMToBitcoinRoute(req, amountDecimal, feeRatio, slippage)
			if code == resp.CodeExternalServerError {
				resp.ExternalServerError(c, msg)
				return
			}
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
	srcChain, ok := new(big.Int).SetString(req.SrcChain, 10)
	if !ok {
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
	amountBigFloat, ok := new(big.Float).SetString(req.Amount)
	if !ok {
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
	dstChain, ok := new(big.Int).SetString(req.DstChain, 10)
	if !ok {
		resp.ParameterErr(c, "invalid dstChain")
		return
	}
	//if !utils.IsValidBitcoinAddress(req.Receiver, logic.NetParams) {
	//	resp.ParameterErr(c, "invalid receiver")
	//	return
	//}

	if utils.IsEmpty(req.Receiver) {
		resp.ParameterErr(c, "missing receiver")
		return
	}
	feeRatio := uint64(0)
	if !utils.IsEmpty(req.FeeRatio) {
		var err error
		feeRatio, err = strconv.ParseUint(req.FeeRatio, 10, 64)
		if err != nil {
			resp.ParameterErr(c, "invalid feeRatio")
			return
		}
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
	if slippage < constants.SlippageMin || slippage > constants.SlippageMax {
		resp.ParameterErr(c, "invalid slippage")
		return
	}

	msg := ""
	code := resp.CodeSuccess
	ret := &entity.SwapResponse{}

	//switch req.SrcChain {
	//case constants.TONChainID:
	//	ret, msg, code = logic.GetSwapFromTON(req.Sender, req.DstChain, req.Receiver, req.FeeCollector, req.FeeRatio, req.Hash)
	//	if code == resp.CodeExternalServerError {
	//		resp.ExternalServerError(c, msg)
	//		return
	//	}
	//	if code != resp.CodeSuccess {
	//		resp.Error(c, code)
	//		return
	//	}
	//default:
	//	if strings.ToLower(req.Hash) == constants.LocalRouteHash {
	//		exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(req.Decimal)), nil)
	//		amount := new(big.Float).Mul(amountBigFloat, new(big.Float).SetInt(exp))
	//		amountBigInt, _ := amount.Int(nil)
	//
	//		ret, msg, code = logic.GetLocalRouteSwapFromEVM(srcChain, req.SrcToken, req.Sender, req.Amount, amountBigFloat, amountBigInt, dstChain, req.DstToken, req.Receiver, slippage)
	//		if code != resp.CodeSuccess {
	//			resp.Error(c, code)
	//			return
	//		}
	//	} else {
	//		ret, msg, code = logic.GetSwapFromEVM(srcChain, req.SrcToken, req.Sender, req.Amount, dstChain, req.DstToken, req.Receiver, req.Hash, slippage)
	//		if code == resp.CodeExternalServerError {
	//			resp.ExternalServerError(c, msg)
	//			return
	//		}
	//		if code != resp.CodeSuccess {
	//			resp.Error(c, code)
	//			return
	//		}
	//	}
	//}

	if strings.ToLower(req.Hash) == constants.LocalRouteHash {
		exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(req.Decimal)), nil)
		amount := new(big.Float).Mul(amountBigFloat, new(big.Float).SetInt(exp))
		amountBigInt, _ := amount.Int(nil)

		if req.DstChain == constants.TONChainID {
			ret, msg, code = logic.GetLocalRouteSwapFromEVMToTON(srcChain, req.SrcToken, req.Sender, req.Amount, amountBigFloat, amountBigInt, dstChain, req.DstToken, req.Receiver, slippage)
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		} else if req.DstChain == constants.BTCChainID {
			ret, msg, code = logic.GetLocalRouteSwapFromEVMToBitcoin(srcChain, req.SrcToken, req.Sender, req.Amount, amountBigFloat, amountBigInt, dstChain, req.DstToken, req.Receiver, slippage)
			if code != resp.CodeSuccess {
				resp.Error(c, code)
				return
			}
		}
	}

	if req.SrcChain == constants.TONChainID {
		ret, msg, code = logic.GetSwapFromTONToEVM(req.Sender, req.DstChain, req.Receiver, req.FeeCollector, req.FeeRatio, req.Hash)
		if code == resp.CodeExternalServerError {
			resp.ExternalServerError(c, msg)
			return
		}
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	} else if req.SrcChain == constants.BTCChainID {
		exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(req.Decimal)), nil)
		amount := new(big.Float).Mul(amountBigFloat, new(big.Float).SetInt(exp))
		amountBigInt, _ := amount.Int(nil)

		ret, msg, code = logic.GetSwapFromBitcoinToEVM(req.SrcChain, req.SrcToken, req.Sender, amountBigFloat, amountBigInt, req.DstChain, req.DstToken, req.Receiver, slippage, req.FeeCollector, feeRatio)
		if code == resp.CodeExternalServerError {
			resp.ExternalServerError(c, msg)
			return
		}
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}

	} else if req.DstChain == constants.TONChainID {
		ret, msg, code = logic.GetSwapFromEVMToTON(srcChain, req.SrcToken, req.Sender, req.Amount, dstChain, req.DstToken, req.Receiver, req.Hash, slippage)
		if code == resp.CodeExternalServerError {
			resp.ExternalServerError(c, msg)
			return
		}
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	} else if req.DstChain == constants.BTCChainID {
		ret, msg, code = logic.GetSwapFromEVMToBitcoin(srcChain, req.SrcToken, req.Sender, req.Amount, dstChain, req.DstToken, req.Receiver, req.Hash, slippage)
		if code == resp.CodeExternalServerError {
			resp.ExternalServerError(c, msg)
			return
		}
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	}
	resp.Success(c, ret)
}
