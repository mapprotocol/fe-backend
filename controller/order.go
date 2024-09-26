package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mapprotocol/fe-backend/constants"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/logic"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/utils"
)

//func CreateOrder(c *gin.Context) {
//	req := &entity.CreateOrderRequest{}
//	if err := c.ShouldBindJSON(req); err != nil {
//		resp.ParameterErr(c, "")
//		return
//	}
//	if req.Action != dao.OrderActionToEVM {
//		resp.ParameterErr(c, "invalid action")
//		return
//	}
//	if !utils.IsValidEvmAddress(req.Receiver) {
//		resp.ParameterErr(c, "invalid receiver")
//		return
//	}
//
//	if req.SrcChain == constants.BTCChainID {
//		if req.SrcToken != constants.BTCTokenAddress {
//			resp.ParameterErr(c, "invalid srcToken")
//			return
//		}
//		if !utils.IsValidBitcoinAddress(req.Sender, logic.NetParams) {
//			resp.ParameterErr(c, "invalid sender")
//			return
//		}
//
//	} else if req.SrcChain == constants.TONChainID {
//		// todo check src token
//		// todo check sender
//
//	} else {
//		resp.ParameterErr(c, "invalid srcChain")
//		return
//	}
//
//	if utils.IsEmpty(req.Amount) {
//		resp.ParameterErr(c, "missing amount")
//		return
//	}
//	if _, err := strconv.ParseFloat(req.Amount, 64); err != nil {
//		resp.ParameterErr(c, "invalid amount")
//		return
//	}
//	if utils.IsEmpty(req.DstChain) {
//		resp.ParameterErr(c, "missing dstChain")
//		return
//	}
//	if !utils.IsValidEvmAddress(req.DstToken) {
//		resp.ParameterErr(c, "invalid dstToken")
//		return
//	}
//	if utils.IsEmpty(req.Hash) {
//		resp.ParameterErr(c, "missing hash")
//		return
//	}
//	if req.Slippage < constants.SlippageMin || req.Slippage > constants.SlippageMax {
//		resp.ParameterErr(c, "invalid slippage")
//		return
//	}
//
//	ret, code := logic.CreateOrder(req.SrcChain, req.SrcToken, req.Sender, req.Amount, req.DstChain, req.DstToken, req.Receiver, req.Action, req.Slippage)
//	if code != resp.CodeSuccess {
//		resp.Error(c, code)
//		return
//	}
//	resp.Success(c, ret)
//}

//func UpdateOrder(c *gin.Context) {
//	req := &entity.UpdateOrderRequest{}
//	if err := c.ShouldBindJSON(req); err != nil {
//		resp.ParameterErr(c, "")
//		return
//	}
//
//	if code := logic.UpdateOrder(req.OrderID, req.InTxHash); code != resp.CodeSuccess {
//		resp.Error(c, code)
//		return
//	}
//	resp.SuccessNil(c)
//}

func OrderList(c *gin.Context) {
	req := &entity.OrderListRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	page, size := utils.ValidatePage(req.Page, req.Size)

	count := int64(0)
	code := resp.CodeSuccess
	var list []*entity.OrderListResponse

	if utils.IsValidBitcoinAddress(req.Sender, logic.NetParams) {
		list, count, code = logic.BitcoinOrderList(req.Sender, page, size)
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	} else {
		list, count, code = logic.OrderList(req.Sender, page, size)
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	}

	resp.SuccessList(c, count, list)
}

func OrderDetail(c *gin.Context) {
	req := &entity.OrderDetailRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	code := resp.CodeSuccess
	ret := &entity.OrderDetailResponse{}

	if utils.IsValidBitcoinAddress(req.Sender, logic.NetParams) {
		ret, code = logic.BitcoinOrderDetail(req.OrderID, req.Sender)
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	} else {
		ret, code = logic.OrderDetail(req.OrderID, req.Sender)
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	}

	resp.Success(c, ret)
}

func OrderDetailByOrderID(c *gin.Context) {
	req := &entity.OrderDetailByOrderIDRequest{}
	if err := c.ShouldBindQuery(req); err != nil {
		resp.ParameterErr(c, "")
		return
	}

	code := resp.CodeSuccess
	ret := &entity.OrderDetailByOrderIDResponse{}
	if req.ChainID == constants.BTCChainID {
		ret, code = logic.BitcoinOrderDetailByOrderID(req.OrderID)
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	} else if req.ChainID == constants.TONChainID || req.ChainID == constants.SOLChainID {
		ret, code = logic.OrderDetailByOrderID(req.OrderID)
		if code != resp.CodeSuccess {
			resp.Error(c, code)
			return
		}
	} else {
		resp.ParameterErr(c, "unsupported chain id")
		return
	}

	resp.Success(c, ret)
}
