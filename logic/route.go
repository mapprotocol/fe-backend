package logic

import (
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/utils"
)

func GetRoute(req *entity.RouteRequest) (ret *entity.RouteResponse, code int) {
	request := &butter.RouteRequest{
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: req.TokenOutAddress,
		Kind:            req.Kind,
		Slippage:        req.Slippage,
		FromChainID:     req.FromChainID,
		ToChainID:       req.ToChainID,
		Amount:          req.Amount,
	}
	route, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		return ret, resp.CodeInternalServerError
	}
	ret = &entity.RouteResponse{
		//Route: nil,
		ButterRoute: route,
	}
	return ret, resp.CodeSuccess
}

func Swap(hash, slippage, from, receiver string) (ret *entity.SwapResponse, code int) {
	request := &butter.SwapRequest{
		Hash:     hash,
		Slippage: slippage,
		From:     from,
		Receiver: receiver,
		CallData: "", // todo build call data
	}
	txData, err := butter.Swap(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter swap")
		return ret, resp.CodeInternalServerError
	}
	ret = &entity.SwapResponse{
		To:      txData.To,
		Data:    txData.Data,
		Value:   txData.Value,
		ChainId: txData.ChainId,
	}
	return ret, resp.CodeSuccess
}
