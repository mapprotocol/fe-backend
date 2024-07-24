package logic

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mapprotocol/fe-backend/constants"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/tonrouter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
	"sync"
)

var BTCToken = entity.Token{
	ChainId:  constants.BTCChainID,
	Address:  constants.BTCTokenAddress,
	Name:     "BTC",
	Decimals: 8,
	Symbol:   "BTC",
}

//func GetBTCRoute(req *entity.RouteRequest) (ret []*entity.RouteResponse, code int) {
//	protocolFeeAmount := ""
//	protocolFeeSymbol := ""
//	tokenIn := entity.Token{}
//	tokenOut := entity.Token{}
//	path := entity.Path{}
//	tonTokenIn := entity.Token{}
//	tonTokenOut := entity.Token{}
//	tonPath := entity.Path{}
//	request := &butter.RouteRequest{
//		TokenInAddress:  req.TokenInAddress,
//		TokenOutAddress: req.TokenOutAddress,
//		Type:            req.Type,
//		Slippage:        req.Slippage,
//		FromChainID:     req.FromChainID,
//		ToChainID:       req.ToChainID,
//		Amount:          req.Amount,
//	}
//	if req.FromChainID == constants.TONChainID || req.ToChainID == constants.TONChainID {
//		tonRequest := &tonrouter.RouteRequest{
//			FromChainID:     req.FromChainID,
//			ToChainID:       req.ToChainID,
//			Amount:          req.Amount,
//			TokenInAddress:  req.TokenInAddress,
//			TokenOutAddress: constants.USDTOfTON,
//			Slippage:        req.Slippage,
//		}
//		route, err := tonrouter.Route(tonRequest)
//		if err != nil {
//			params := map[string]interface{}{
//				"request": utils.JSON(request),
//				"error":   err,
//			}
//			log.Logger().WithFields(params).Error("failed to request ton route")
//			return ret, resp.CodeInternalServerError
//		}
//		in := route.SrcChain.Route[0].Path[0].TokenIn
//		tonTokenIn = entity.Token{
//			ChainId:  route.SrcChain.ChainId,
//			Address:  in.Address,
//			Name:     in.Name,
//			Decimals: in.Decimals,
//			Symbol:   in.Symbol,
//			Icon:     in.Image,
//		}
//
//		out := route.SrcChain.Route[0].Path[0].TokenOut
//		tonTokenOut = entity.Token{
//			ChainId:  route.SrcChain.ChainId,
//			Address:  out.Address,
//			Name:     out.Name,
//			Decimals: out.Decimals,
//			Symbol:   out.Symbol,
//			Icon:     in.Image,
//		}
//
//		tonPath = entity.Path{
//			Name:      route.SrcChain.Route[0].DexName,
//			AmountIn:  route.SrcChain.TokenAmountIn,
//			AmountOut: route.SrcChain.TokenAmountOut,
//			TokenIn:   tonTokenIn,
//			TokenOut:  tonTokenOut,
//		}
//	}
//
//	if req.Action == dao.OrderActionToEVM {
//		request.FromChainID = constants.ChainPollChainID
//		if req.FromChainID == constants.BTCChainID {
//			request.TokenInAddress = constants.WBTCOfChainPoll
//
//			tokenIn = BTCToken
//			path = entity.Path{
//				Name:      constants.ExchangeNameFlushExchange,
//				AmountIn:  request.Amount,
//				AmountOut: request.Amount, // todo fee rate
//				TokenIn:   BTCToken,
//				TokenOut:  BTCToken,
//			}
//			protocolFeeAmount = "0" // todo todo fee rate
//			protocolFeeSymbol = "BTC"
//
//		} else if req.FromChainID == constants.TONChainID {
//			request.TokenInAddress = constants.USDTOfChainPoll
//
//			tokenIn = tonTokenIn
//			path = tonPath
//			protocolFeeAmount = "0"
//			protocolFeeSymbol = "USDT" // todo
//
//		}
//	} else if req.Action == dao.OrderActionFromEVM {
//		request.ToChainID = constants.ChainPollChainID
//		protocolFeeSymbol = request.TokenInAddress
//		if req.FromChainID == constants.BTCChainID {
//			request.TokenOutAddress = constants.WBTCOfChainPoll
//
//			tokenOut = BTCToken
//
//			path = entity.Path{
//				Name:      constants.ExchangeNameFlushExchange,
//				AmountIn:  request.Amount, // todo exchanged amount
//				AmountOut: request.Amount, // todo fee rate
//				TokenIn:   BTCToken,
//				TokenOut:  BTCToken,
//			}
//			protocolFeeAmount = "0" // todo exchanged amount
//			protocolFeeSymbol = "BTC"
//		} else if req.FromChainID == constants.TONChainID {
//			request.TokenOutAddress = constants.USDTOfChainPoll
//
//			tokenOut = tonTokenOut
//			path = tonPath
//			protocolFeeAmount = "0"
//			protocolFeeSymbol = "USDT" // todo
//		}
//	}
//	route, err := butter.Route(request)
//	if err != nil {
//		params := map[string]interface{}{
//			"request": utils.JSON(request),
//			"error":   err,
//		}
//		log.Logger().WithFields(params).Error("failed to request butter route")
//		return ret, resp.CodeInternalServerError
//	}
//
//	ret = make([]*entity.RouteResponse, 0, len(route))
//	for _, r := range route {
//		n := &entity.RouteResponse{
//			Hash: r.Hash,
//			TokenIn: entity.Token{
//				Address:  r.SrcChain.TokenIn.Address,
//				Name:     r.SrcChain.TokenIn.Name,
//				Decimals: r.SrcChain.TokenIn.Decimals,
//				Symbol:   r.SrcChain.TokenIn.Symbol,
//				Icon:     r.SrcChain.TokenIn.Icon,
//			},
//			TokenOut: entity.Token{
//				Address:  r.DstChain.TokenOut.Address,
//				Name:     r.DstChain.TokenOut.Name,
//				Decimals: r.DstChain.TokenOut.Decimals,
//				Symbol:   r.DstChain.TokenOut.Symbol,
//				Icon:     r.DstChain.TokenOut.Icon,
//			},
//			AmountIn:  r.SrcChain.TotalAmountIn,
//			AmountOut: r.DstChain.TotalAmountOut,
//			Path: []entity.Path{
//				{
//					Name:      r.SrcChain.Bridge,
//					AmountIn:  r.SrcChain.TotalAmountIn,
//					AmountOut: r.SrcChain.TotalAmountOut,
//					TokenIn: entity.Token{
//						ChainId:  r.SrcChain.ChainId,
//						Address:  r.SrcChain.TokenIn.Address,
//						Name:     r.SrcChain.TokenIn.Name,
//						Decimals: r.SrcChain.TokenIn.Decimals,
//						Symbol:   r.SrcChain.TokenIn.Symbol,
//						Icon:     r.SrcChain.TokenIn.Icon,
//					},
//					TokenOut: entity.Token{
//						ChainId:  r.SrcChain.ChainId,
//						Address:  r.SrcChain.TokenOut.Address,
//						Name:     r.SrcChain.TokenOut.Name,
//						Decimals: r.SrcChain.TokenOut.Decimals,
//						Symbol:   r.SrcChain.TokenOut.Symbol,
//						Icon:     r.SrcChain.TokenOut.Icon,
//					},
//				},
//				{
//					Name:      r.BridgeChain.Bridge,
//					AmountIn:  r.BridgeChain.TotalAmountIn,
//					AmountOut: r.BridgeChain.TotalAmountOut,
//					TokenIn: entity.Token{
//						ChainId:  r.BridgeChain.ChainId,
//						Address:  r.BridgeChain.TokenIn.Address,
//						Name:     r.BridgeChain.TokenIn.Name,
//						Decimals: r.BridgeChain.TokenIn.Decimals,
//						Symbol:   r.BridgeChain.TokenIn.Symbol,
//						Icon:     r.BridgeChain.TokenIn.Icon,
//					},
//					TokenOut: entity.Token{
//						ChainId:  r.BridgeChain.ChainId,
//						Address:  r.BridgeChain.TokenOut.Address,
//						Name:     r.BridgeChain.TokenOut.Name,
//						Decimals: r.BridgeChain.TokenOut.Decimals,
//						Symbol:   r.BridgeChain.TokenOut.Symbol,
//						Icon:     r.BridgeChain.TokenOut.Icon,
//					},
//				},
//				{
//					Name:      r.DstChain.Bridge,
//					AmountIn:  r.DstChain.TotalAmountIn,
//					AmountOut: r.DstChain.TotalAmountOut,
//					TokenIn: entity.Token{
//						ChainId:  r.DstChain.ChainId,
//						Address:  r.DstChain.TokenIn.Address,
//						Name:     r.DstChain.TokenIn.Name,
//						Decimals: r.DstChain.TokenIn.Decimals,
//						Symbol:   r.DstChain.TokenIn.Symbol,
//						Icon:     r.DstChain.TokenIn.Icon,
//					},
//					TokenOut: entity.Token{
//						ChainId:  r.SrcChain.ChainId,
//						Address:  r.SrcChain.TokenOut.Address,
//						Name:     r.SrcChain.TokenOut.Name,
//						Decimals: r.SrcChain.TokenOut.Decimals,
//						Symbol:   r.SrcChain.TokenOut.Symbol,
//						Icon:     r.SrcChain.TokenOut.Icon,
//					},
//				},
//			},
//			GasFee: entity.Fee{
//				Amount: r.GasFee.Amount,
//				Symbol: r.GasFee.Symbol,
//			},
//			BridgeFee: entity.Fee{
//				Amount: r.BridgeFee.Amount,
//				Symbol: r.BridgeFee.Symbol,
//			},
//			ProtocolFee: entity.Fee{
//				Amount: protocolFeeAmount,
//				Symbol: protocolFeeSymbol,
//			},
//		}
//		if req.Action == dao.OrderActionToEVM {
//			n.TokenIn = tokenIn
//			n.Path = append([]entity.Path{path}, n.Path...)
//		} else if req.Action == dao.OrderActionFromEVM {
//			n.TokenOut = tokenOut
//			n.Path = append(n.Path, path)
//		}
//
//		ret = append(ret, n)
//	}
//	return ret, resp.CodeSuccess
//}

//func GetTONRoute(req *entity.RouteRequest) (ret []*entity.RouteResponse, code int) {
//	protocolFeeAmount := ""
//	protocolFeeSymbol := ""
//	tonTokenIn := entity.Token{}
//	tonTokenOut := entity.Token{}
//
//	tonRequest := &tonrouter.RouteRequest{
//		FromChainID:     req.FromChainID,
//		ToChainID:       req.ToChainID,
//		Amount:          req.Amount,
//		TokenInAddress:  req.TokenInAddress,
//		TokenOutAddress: constants.USDTOfTON,
//		Slippage:        req.Slippage,
//	}
//	tonRoute, err := tonrouter.Route(tonRequest)
//	if err != nil {
//		params := map[string]interface{}{
//			"request": utils.JSON(tonRequest),
//			"error":   err,
//		}
//		log.Logger().WithFields(params).Error("failed to request ton route")
//		return ret, resp.CodeInternalServerError
//	}
//	in := tonRoute.SrcChain.Route[0].Path[0].TokenIn
//	tonTokenIn = entity.Token{
//		ChainId:  tonRoute.SrcChain.ChainId,
//		Address:  in.Address,
//		Name:     in.Name,
//		Decimals: in.Decimals,
//		Symbol:   in.Symbol,
//		Icon:     in.Image,
//	}
//
//	out := tonRoute.SrcChain.Route[0].Path[0].TokenOut
//	tonTokenOut = entity.Token{
//		ChainId:  tonRoute.SrcChain.ChainId,
//		Address:  out.Address,
//		Name:     out.Name,
//		Decimals: out.Decimals,
//		Symbol:   out.Symbol,
//		Icon:     in.Image,
//	}
//
//	//path = entity.Path{
//	//	Name:      route.SrcChain.Route[0].DexName,
//	//	AmountIn:  route.SrcChain.TokenAmountIn,
//	//	AmountOut: route.SrcChain.TokenAmountOut,
//	//	TokenIn:   tonTokenIn,
//	//	TokenOut:  tonTokenOut,
//	//}
//
//	request := &butter.RouteRequest{
//		TokenInAddress:  req.TokenInAddress,
//		TokenOutAddress: req.TokenOutAddress,
//		Type:            req.Type,
//		Slippage:        req.Slippage,
//		FromChainID:     req.FromChainID,
//		ToChainID:       req.ToChainID,
//		Amount:          tonRoute.SrcChain.TokenAmountOut,
//	}
//
//	switch req.Action {
//	case dao.OrderActionToEVM:
//		request.FromChainID = constants.ChainPollChainID
//		request.TokenInAddress = constants.USDTOfChainPoll
//
//		protocolFeeAmount = "0"
//		protocolFeeSymbol = "USDT" // todo
//	case dao.OrderActionFromEVM:
//		request.ToChainID = constants.ChainPollChainID
//		request.TokenOutAddress = constants.USDTOfChainPoll
//
//		protocolFeeAmount = "0"
//		protocolFeeSymbol = "USDT" // todo
//	}
//
//	butterRoute, err := butter.Route(request)
//	if err != nil {
//		params := map[string]interface{}{
//			"request": utils.JSON(request),
//			"error":   err,
//		}
//		log.Logger().WithFields(params).Error("failed to request butter route")
//		return ret, resp.CodeInternalServerError
//	}
//
//	ret = make([]*entity.RouteResponse, 0, len(butterRoute))
//	for _, r := range butterRoute {
//		var (
//			path      []entity.Path
//			tokenIn   entity.Token
//			tokenOut  entity.Token
//			amountIn  string
//			amountOut string
//		)
//		butterSrcChainTokenIn := entity.Token{
//			ChainId:  r.SrcChain.ChainId,
//			Address:  r.SrcChain.TokenIn.Address,
//			Name:     r.SrcChain.TokenIn.Name,
//			Decimals: r.SrcChain.TokenIn.Decimals,
//			Symbol:   r.SrcChain.TokenIn.Symbol,
//			Icon:     r.SrcChain.TokenIn.Icon,
//		}
//
//		switch req.Action {
//		case dao.OrderActionToEVM:
//			tokenIn = tonTokenIn
//			tokenOut = entity.Token{
//				ChainId:  r.DstChain.ChainId,
//				Address:  r.DstChain.TokenOut.Address,
//				Name:     r.DstChain.TokenOut.Name,
//				Decimals: r.DstChain.TokenOut.Decimals,
//				Symbol:   r.DstChain.TokenOut.Symbol,
//				Icon:     r.DstChain.TokenOut.Icon,
//			}
//			amountIn = tonRoute.SrcChain.TokenAmountIn
//			amountOut = r.DstChain.TotalAmountOut
//			path = []entity.Path{
//				{
//					Name:      tonRoute.SrcChain.Route[0].DexName,
//					AmountIn:  tonRoute.SrcChain.TokenAmountIn,
//					AmountOut: tonRoute.SrcChain.TokenAmountOut,
//					TokenIn:   tonTokenIn,
//					TokenOut:  tonTokenOut,
//				},
//				{
//					Name:      constants.ExchangeNameFlushExchange,
//					AmountIn:  tonRoute.SrcChain.TokenAmountOut,
//					AmountOut: tonRoute.SrcChain.TokenAmountOut,
//					TokenIn:   tonTokenOut,
//					TokenOut: entity.Token{
//						ChainId:  r.SrcChain.ChainId,
//						Address:  r.SrcChain.TokenIn.Address,
//						Name:     r.SrcChain.TokenIn.Name,
//						Decimals: r.SrcChain.TokenIn.Decimals,
//						Symbol:   r.SrcChain.TokenIn.Symbol,
//						Icon:     r.SrcChain.TokenIn.Icon,
//					},
//				},
//				{
//					Name:      constants.ExchangeNameButter,
//					AmountIn:  tonRoute.SrcChain.TokenAmountOut,
//					AmountOut: tonRoute.SrcChain.TokenAmountOut,
//					TokenIn: entity.Token{
//						ChainId:  r.SrcChain.ChainId,
//						Address:  r.SrcChain.TokenIn.Address,
//						Name:     r.SrcChain.TokenIn.Name,
//						Decimals: r.SrcChain.TokenIn.Decimals,
//						Symbol:   r.SrcChain.TokenIn.Symbol,
//						Icon:     r.SrcChain.TokenIn.Icon,
//					},
//					TokenOut: tokenOut,
//				},
//			}
//		case dao.OrderActionFromEVM:
//			tokenIn = entity.Token{
//				ChainId:  r.SrcChain.ChainId,
//				Address:  r.SrcChain.TokenIn.Address,
//				Name:     r.SrcChain.TokenIn.Name,
//				Decimals: r.SrcChain.TokenIn.Decimals,
//				Symbol:   r.SrcChain.TokenIn.Symbol,
//				Icon:     r.SrcChain.TokenIn.Icon,
//			}
//			tokenOut = tonTokenOut
//			amountIn = r.SrcChain.TotalAmountIn
//			amountOut = tonRoute.SrcChain.TokenAmountOut
//		}
//
//		//if req.Action == dao.OrderActionToEVM {
//		//	tokenIn = tonTokenIn
//		//	tokenOut = entity.Token{
//		//		ChainId:  r.DstChain.ChainId,
//		//		Address:  r.DstChain.TokenOut.Address,
//		//		Name:     r.DstChain.TokenOut.Name,
//		//		Decimals: r.DstChain.TokenOut.Decimals,
//		//		Symbol:   r.DstChain.TokenOut.Symbol,
//		//		Icon:     r.DstChain.TokenOut.Icon,
//		//	}
//		//	amountIn = tonRoute.SrcChain.TokenAmountIn
//		//	amountOut = r.DstChain.TotalAmountOut
//		//	path = []entity.Path{
//		//		{
//		//			Name:      tonRoute.SrcChain.Route[0].DexName,
//		//			AmountIn:  tonRoute.SrcChain.TokenAmountIn,
//		//			AmountOut: tonRoute.SrcChain.TokenAmountOut,
//		//			TokenIn:   tonTokenIn,
//		//			TokenOut:  tonTokenOut,
//		//		},
//		//		{
//		//			Name:      constants.ExchangeNameFlushExchange,
//		//			AmountIn:  tonRoute.SrcChain.TokenAmountOut,
//		//			AmountOut: tonRoute.SrcChain.TokenAmountOut,
//		//			TokenIn:   tonTokenOut,
//		//			TokenOut: entity.Token{
//		//				ChainId:  r.SrcChain.ChainId,
//		//				Address:  r.SrcChain.TokenIn.Address,
//		//				Name:     r.SrcChain.TokenIn.Name,
//		//				Decimals: r.SrcChain.TokenIn.Decimals,
//		//				Symbol:   r.SrcChain.TokenIn.Symbol,
//		//				Icon:     r.SrcChain.TokenIn.Icon,
//		//			},
//		//		},
//		//		{
//		//			Name:      constants.ExchangeNameButter,
//		//			AmountIn:  tonRoute.SrcChain.TokenAmountOut,
//		//			AmountOut: tonRoute.SrcChain.TokenAmountOut,
//		//			TokenIn: entity.Token{
//		//				ChainId:  r.SrcChain.ChainId,
//		//				Address:  r.SrcChain.TokenIn.Address,
//		//				Name:     r.SrcChain.TokenIn.Name,
//		//				Decimals: r.SrcChain.TokenIn.Decimals,
//		//				Symbol:   r.SrcChain.TokenIn.Symbol,
//		//				Icon:     r.SrcChain.TokenIn.Icon,
//		//			},
//		//			TokenOut: tokenOut,
//		//		},
//		//	}
//		//} else if req.Action == dao.OrderActionFromEVM {
//		//	//tokneOut := tonTokenOut
//		//}
//
//		n := &entity.RouteResponse{
//			Hash:      r.Hash,
//			TokenIn:   tokenIn,
//			TokenOut:  tokenOut,
//			AmountIn:  amountIn,
//			AmountOut: amountOut,
//			Path:      path,
//			GasFee: entity.Fee{
//				Amount: r.GasFee.Amount,
//				Symbol: r.GasFee.Symbol,
//			},
//			BridgeFee: entity.Fee{
//				Amount: r.BridgeFee.Amount,
//				Symbol: r.BridgeFee.Symbol,
//			},
//			ProtocolFee: entity.Fee{
//				Amount: protocolFeeAmount,
//				Symbol: protocolFeeSymbol,
//			},
//		}
//		//if req.Action == dao.OrderActionToEVM {
//		//	n.TokenIn = tokenIn
//		//	n.Path = append([]entity.Path{path}, n.Path...)
//		//} else if req.Action == dao.OrderActionFromEVM {
//		//	n.TokenOut = tokenOut
//		//	n.Path = append(n.Path, path)
//		//}
//
//		ret = append(ret, n)
//	}
//	return ret, resp.CodeSuccess
//}

//func GetTONToEVMRoute(req *entity.RouteRequest) (ret []*entity.RouteResponse, code int) {
//	var (
//		protocolFeeAmount string
//		protocolFeeSymbol string
//		tonTokenIn        entity.Token
//		tonTokenOut       entity.Token
//	)
//
//	tonRequest := &tonrouter.RouteRequest{
//		FromChainID:     req.FromChainID,
//		ToChainID:       req.ToChainID,
//		Amount:          req.Amount,
//		TokenInAddress:  req.TokenInAddress,
//		TokenOutAddress: constants.USDTOfTON,
//		Slippage:        req.Slippage,
//	}
//	tonRoute, err := tonrouter.Route(tonRequest)
//	if err != nil {
//		params := map[string]interface{}{
//			"request": utils.JSON(tonRequest),
//			"error":   err,
//		}
//		log.Logger().WithFields(params).Error("failed to request ton route")
//		return ret, resp.CodeInternalServerError
//	}
//	in := tonRoute.SrcChain.Route[0].Path[0].TokenIn
//	tonTokenIn = entity.Token{
//		ChainId:  tonRoute.SrcChain.ChainId,
//		Address:  in.Address,
//		Name:     in.Name,
//		Decimals: in.Decimals,
//		Symbol:   in.Symbol,
//		Icon:     in.Image,
//	}
//
//	out := tonRoute.SrcChain.Route[0].Path[0].TokenOut
//	tonTokenOut = entity.Token{
//		ChainId:  tonRoute.SrcChain.ChainId,
//		Address:  out.Address,
//		Name:     out.Name,
//		Decimals: out.Decimals,
//		Symbol:   out.Symbol,
//		Icon:     in.Image,
//	}
//
//	request := &butter.RouteRequest{
//		TokenInAddress:  req.TokenInAddress,
//		TokenOutAddress: req.TokenOutAddress,
//		Type:            req.Type,
//		Slippage:        req.Slippage,
//		FromChainID:     req.FromChainID,
//		ToChainID:       req.ToChainID,
//		Amount:          tonRoute.SrcChain.TokenAmountOut,
//	}
//
//	switch req.Action {
//	case dao.OrderActionToEVM:
//		request.FromChainID = constants.ChainPollChainID
//		request.TokenInAddress = constants.USDTOfChainPoll
//
//		protocolFeeAmount = "0"
//		protocolFeeSymbol = "USDT" // todo
//	case dao.OrderActionFromEVM:
//		request.ToChainID = constants.ChainPollChainID
//		request.TokenOutAddress = constants.USDTOfChainPoll
//
//		protocolFeeAmount = "0"
//		protocolFeeSymbol = "USDT" // todo
//	}
//
//	butterRoute, err := butter.Route(request)
//	if err != nil {
//		params := map[string]interface{}{
//			"request": utils.JSON(request),
//			"error":   err,
//		}
//		log.Logger().WithFields(params).Error("failed to request butter route")
//		return ret, resp.CodeInternalServerError
//	}
//
//	ret = make([]*entity.RouteResponse, 0, len(butterRoute))
//	for _, r := range butterRoute {
//		var (
//			amountIn  string
//			amountOut string
//			path      []entity.Path
//			tokenIn   entity.Token
//			tokenOut  entity.Token
//		)
//		butterSrcChainTokenIn := entity.Token{
//			ChainId:  r.SrcChain.ChainId,
//			Address:  r.SrcChain.TokenIn.Address,
//			Name:     r.SrcChain.TokenIn.Name,
//			Decimals: r.SrcChain.TokenIn.Decimals,
//			Symbol:   r.SrcChain.TokenIn.Symbol,
//			Icon:     r.SrcChain.TokenIn.Icon,
//		}
//		butterDstChainTokenOut := entity.Token{
//			ChainId:  r.DstChain.ChainId,
//			Address:  r.DstChain.TokenOut.Address,
//			Name:     r.DstChain.TokenOut.Name,
//			Decimals: r.DstChain.TokenOut.Decimals,
//			Symbol:   r.DstChain.TokenOut.Symbol,
//			Icon:     r.DstChain.TokenOut.Icon,
//		}
//
//		switch req.Action {
//		case dao.OrderActionToEVM:
//			tokenIn = tonTokenIn
//			tokenOut = butterDstChainTokenOut
//			amountIn = tonRoute.SrcChain.TokenAmountIn
//			amountOut = r.DstChain.TotalAmountOut
//			path = []entity.Path{
//				{
//					Name:      tonRoute.SrcChain.Route[0].DexName,
//					AmountIn:  tonRoute.SrcChain.TokenAmountIn,
//					AmountOut: tonRoute.SrcChain.TokenAmountOut,
//					TokenIn:   tonTokenIn,
//					TokenOut:  tonTokenOut,
//				},
//				{
//					Name:      constants.ExchangeNameFlushExchange,
//					AmountIn:  tonRoute.SrcChain.TokenAmountOut,
//					AmountOut: r.SrcChain.TotalAmountIn,
//					TokenIn:   tonTokenOut,
//					TokenOut:  butterSrcChainTokenIn,
//				},
//				{
//					Name:      constants.ExchangeNameButter,
//					AmountIn:  r.SrcChain.TotalAmountIn,
//					AmountOut: r.DstChain.TotalAmountOut,
//					TokenIn:   butterSrcChainTokenIn,
//					TokenOut:  butterDstChainTokenOut,
//				},
//			}
//		case dao.OrderActionFromEVM:
//			tokenIn = butterSrcChainTokenIn
//			tokenOut = tonTokenOut
//			amountIn = r.SrcChain.TotalAmountIn
//			amountOut = tonRoute.SrcChain.TokenAmountOut
//			path = []entity.Path{
//				{
//					Name:      constants.ExchangeNameButter,
//					AmountIn:  r.SrcChain.TotalAmountIn,
//					AmountOut: r.DstChain.TotalAmountOut,
//					TokenIn:   butterSrcChainTokenIn,
//					TokenOut:  butterDstChainTokenOut,
//				},
//				{
//					Name:      constants.ExchangeNameFlushExchange,
//					AmountIn:  r.DstChain.TotalAmountOut,
//					AmountOut: tonRoute.SrcChain.TokenAmountIn,
//					TokenIn:   butterDstChainTokenOut,
//					TokenOut:  tonTokenIn,
//				},
//				{
//					Name:      tonRoute.SrcChain.Route[0].DexName,
//					AmountIn:  tonRoute.SrcChain.TokenAmountIn,
//					AmountOut: tonRoute.SrcChain.TokenAmountOut,
//					TokenIn:   tonTokenIn,
//					TokenOut:  tonTokenOut,
//				},
//			}
//		}
//
//		n := &entity.RouteResponse{
//			Hash:      r.Hash,
//			TokenIn:   tokenIn,
//			TokenOut:  tokenOut,
//			AmountIn:  amountIn,
//			AmountOut: amountOut,
//			Path:      path,
//			GasFee: entity.Fee{
//				Amount: r.GasFee.Amount,
//				Symbol: r.GasFee.Symbol,
//			},
//			BridgeFee: entity.Fee{
//				Amount: r.BridgeFee.Amount,
//				Symbol: r.BridgeFee.Symbol,
//			},
//			ProtocolFee: entity.Fee{
//				Amount: protocolFeeAmount,
//				Symbol: protocolFeeSymbol,
//			},
//		}
//		ret = append(ret, n)
//	}
//	return ret, resp.CodeSuccess
//}

func GetTONToEVMRoute(req *entity.RouteRequest) (ret []*entity.RouteResponse, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	tonRequest := &tonrouter.RouteRequest{
		FromChainID:     req.FromChainID,
		ToChainID:       req.ToChainID,
		Amount:          req.Amount,
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: constants.USDTOfTON,
		Slippage:        req.Slippage,
	}
	tonRoute, err := tonrouter.Route(tonRequest)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(tonRequest),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request ton route")
		return ret, resp.CodeInternalServerError
	}
	in := tonRoute.SrcChain.Route[0].Path[0].TokenIn
	tonTokenIn = entity.Token{
		ChainId:  tonRoute.SrcChain.ChainId,
		Address:  in.Address,
		Name:     in.Name,
		Decimals: in.Decimals,
		Symbol:   in.Symbol,
		Icon:     in.Image,
	}

	out := tonRoute.SrcChain.Route[0].Path[0].TokenOut
	tonTokenOut = entity.Token{
		ChainId:  tonRoute.SrcChain.ChainId,
		Address:  out.Address,
		Name:     out.Name,
		Decimals: out.Decimals,
		Symbol:   out.Symbol,
		Icon:     in.Image,
	}

	request := &butter.RouteRequest{
		TokenInAddress:  constants.USDTOfChainPoll,
		TokenOutAddress: req.TokenOutAddress,
		Type:            req.Type,
		Slippage:        req.Slippage,
		FromChainID:     constants.ChainPollChainID,
		ToChainID:       req.ToChainID,
		Amount:          tonRoute.SrcChain.TokenAmountOut,
	}

	butterRoute, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		return ret, resp.CodeInternalServerError
	}

	ret = make([]*entity.RouteResponse, 0, len(butterRoute))
	for _, r := range butterRoute {
		butterSrcChainTokenIn := entity.Token{
			ChainId:  r.SrcChain.ChainId,
			Address:  r.SrcChain.TokenIn.Address,
			Name:     r.SrcChain.TokenIn.Name,
			Decimals: r.SrcChain.TokenIn.Decimals,
			Symbol:   r.SrcChain.TokenIn.Symbol,
			Icon:     r.SrcChain.TokenIn.Icon,
		}
		butterDstChainTokenOut := entity.Token{
			ChainId:  r.DstChain.ChainId,
			Address:  r.DstChain.TokenOut.Address,
			Name:     r.DstChain.TokenOut.Name,
			Decimals: r.DstChain.TokenOut.Decimals,
			Symbol:   r.DstChain.TokenOut.Symbol,
			Icon:     r.DstChain.TokenOut.Icon,
		}

		n := &entity.RouteResponse{
			Hash:      r.Hash,
			TokenIn:   tonTokenIn,
			TokenOut:  butterDstChainTokenOut,
			AmountIn:  tonRoute.SrcChain.TokenAmountIn,
			AmountOut: r.DstChain.TotalAmountOut,
			Path: []entity.Path{
				{
					Name:      tonRoute.SrcChain.Route[0].DexName,
					AmountIn:  tonRoute.SrcChain.TokenAmountIn,
					AmountOut: tonRoute.SrcChain.TokenAmountOut,
					TokenIn:   tonTokenIn,
					TokenOut:  tonTokenOut,
				},
				{
					Name:      constants.ExchangeNameFlushExchange,
					AmountIn:  tonRoute.SrcChain.TokenAmountOut,
					AmountOut: r.SrcChain.TotalAmountIn,
					TokenIn:   tonTokenOut,
					TokenOut:  butterSrcChainTokenIn,
				},
				{
					Name:      constants.ExchangeNameButter,
					AmountIn:  r.SrcChain.TotalAmountIn,
					AmountOut: r.DstChain.TotalAmountOut,
					TokenIn:   butterSrcChainTokenIn,
					TokenOut:  butterDstChainTokenOut,
				},
			},
			GasFee: entity.Fee{
				Amount: r.GasFee.Amount,
				Symbol: r.GasFee.Symbol,
			},
			BridgeFee: entity.Fee{
				Amount: r.BridgeFee.Amount,
				Symbol: r.BridgeFee.Symbol,
			},
			ProtocolFee: entity.Fee{
				Amount: "0",
				Symbol: "USDT",
			},
		}
		ret = append(ret, n)
	}
	return ret, resp.CodeSuccess
}

func GetEVMToTONRoute(req *entity.RouteRequest, slippage uint64) (ret []*entity.RouteResponse, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	request := &butter.RouteRequest{
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: constants.USDTOfChainPoll,
		Type:            req.Type,
		Slippage:        strconv.FormatUint(slippage/3*2, 10),
		FromChainID:     req.FromChainID,
		ToChainID:       constants.ChainPollChainID,
		Amount:          req.Amount,
	}
	butterRoutes, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		return ret, resp.CodeInternalServerError
	}
	if len(butterRoutes) == 0 {
		return ret, resp.CodeButterNotAvailableRoute
	}

	tonRequest := &tonrouter.RouteRequest{
		FromChainID:     req.FromChainID,
		ToChainID:       req.ToChainID,
		TokenInAddress:  constants.USDTOfTON,
		TokenOutAddress: req.TokenOutAddress,
		Slippage:        strconv.FormatUint(slippage/3, 10),
	}
	tonRoutes, err := getTONRoutes(tonRequest, butterRoutes) // todo skip error ?
	if err != nil {
		return ret, resp.CodeTONRouteServerError
	}
	if len(tonRoutes) != len(butterRoutes) {
		return ret, resp.CodeTONRouteServerError
	}

	ret = make([]*entity.RouteResponse, 0, len(butterRoutes))
	for _, r := range butterRoutes {
		tonRoute, ok := tonRoutes[r.Hash]
		if !ok {
			continue
		}

		in := tonRoute.SrcChain.Route[0].Path[0].TokenIn
		tonTokenIn = entity.Token{
			ChainId:  tonRoute.SrcChain.ChainId,
			Address:  in.Address,
			Name:     in.Name,
			Decimals: in.Decimals,
			Symbol:   in.Symbol,
			Icon:     in.Image,
		}

		out := tonRoute.SrcChain.Route[0].Path[0].TokenOut
		tonTokenOut = entity.Token{
			ChainId:  tonRoute.SrcChain.ChainId,
			Address:  out.Address,
			Name:     out.Name,
			Decimals: out.Decimals,
			Symbol:   out.Symbol,
			Icon:     in.Image,
		}

		butterSrcChainTokenIn := entity.Token{
			ChainId:  r.SrcChain.ChainId,
			Address:  r.SrcChain.TokenIn.Address,
			Name:     r.SrcChain.TokenIn.Name,
			Decimals: r.SrcChain.TokenIn.Decimals,
			Symbol:   r.SrcChain.TokenIn.Symbol,
			Icon:     r.SrcChain.TokenIn.Icon,
		}
		butterDstChainTokenOut := entity.Token{
			ChainId:  r.DstChain.ChainId,
			Address:  r.DstChain.TokenOut.Address,
			Name:     r.DstChain.TokenOut.Name,
			Decimals: r.DstChain.TokenOut.Decimals,
			Symbol:   r.DstChain.TokenOut.Symbol,
			Icon:     r.DstChain.TokenOut.Icon,
		}

		n := &entity.RouteResponse{
			Hash:      r.Hash,
			TokenIn:   butterSrcChainTokenIn,
			TokenOut:  tonTokenOut,
			AmountIn:  r.SrcChain.TotalAmountIn,
			AmountOut: tonRoute.SrcChain.TokenAmountOut,
			Path: []entity.Path{
				{
					Name:      constants.ExchangeNameButter,
					AmountIn:  r.SrcChain.TotalAmountIn,
					AmountOut: r.DstChain.TotalAmountOut,
					TokenIn:   butterSrcChainTokenIn,
					TokenOut:  butterDstChainTokenOut,
				},
				{
					Name:      constants.ExchangeNameFlushExchange,
					AmountIn:  r.DstChain.TotalAmountOut,
					AmountOut: tonRoute.SrcChain.TokenAmountIn,
					TokenIn:   butterDstChainTokenOut,
					TokenOut:  tonTokenIn,
				},
				{
					Name:      tonRoute.SrcChain.Route[0].DexName,
					AmountIn:  tonRoute.SrcChain.TokenAmountIn,
					AmountOut: tonRoute.SrcChain.TokenAmountOut,
					TokenIn:   tonTokenIn,
					TokenOut:  tonTokenOut,
				},
			},
			GasFee: entity.Fee{
				Amount: r.GasFee.Amount,
				Symbol: r.GasFee.Symbol,
			},
			BridgeFee: entity.Fee{
				Amount: r.BridgeFee.Amount,
				Symbol: r.BridgeFee.Symbol,
			},
			ProtocolFee: entity.Fee{
				Amount: "0",
				Symbol: "USDT",
			},
		}
		ret = append(ret, n)
	}
	return ret, resp.CodeSuccess
}

//func Swap(srcChain, srcToken, sender string, amount *big.Int, receiver, hash string, slippage uint64) (ret *entity.SwapResponse, code int) {
//	order := &dao.Order{
//		SrcChain: srcChain,
//		SrcToken: srcToken,
//		Sender:   sender,
//		DstChain: constants.BTCChainID,
//		DstToken: constants.BTCTokenAddress,
//		Receiver: receiver,
//		Action:   dao.OrderActionFromEVM,
//		Stage:    dao.OrderStag1,
//		Status:   dao.OrderStatusPending,
//		Slippage: slippage,
//	}
//	orderID, err := order.Create()
//	if err != nil {
//		log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
//		return nil, resp.CodeInternalServerError
//	}
//
//	orderIDByte32 := utils.Uint64ToByte32(orderID)
//	// PackOnReceived(amount *big.Int, orderId [32]byte, token common.Address, from common.Address, to []byte)
//	// todo amount
//	// todo src token
//	packed, err := PackOnReceived(amount, orderIDByte32, common.HexToAddress(srcToken), common.HexToAddress(sender), []byte(receiver))
//	if err != nil {
//		params := map[string]interface{}{
//			"amount":        amount,
//			"orderID":       order,
//			"orderIDByte32": orderIDByte32,
//			"srcToken":      srcToken, // src token
//			"sender":        sender,
//			"receiver":      receiver,
//			"error":         err,
//		}
//		log.Logger().WithFields(params).Error("failed to pack onReceived")
//		return ret, resp.CodeInternalServerError
//	}
//	encodedCallback, err := EncodeSwapCallbackParams(common.HexToAddress(viper.GetString("feRouterContract")), common.HexToAddress(sender), packed) // todo sender
//	if err != nil {
//		params := map[string]interface{}{
//			"feRouter": viper.GetString("feRouterContract"),
//			"sender":   sender,
//			"packed":   hex.EncodeToString(packed),
//			"error":    err,
//		}
//		log.Logger().WithFields(params).Error("failed to encode swap callback params")
//		return ret, resp.CodeInternalServerError
//	}
//
//	request := &butter.SwapRequest{
//		Hash:     hash,
//		Slippage: slippage,
//		From:     sender,
//		Receiver: viper.GetString("butterRouterContract"),
//		CallData: encodedCallback,
//	}
//	txData, err := butter.Swap(request)
//	if err != nil {
//		params := map[string]interface{}{
//			"request": utils.JSON(request),
//			"error":   err,
//		}
//		log.Logger().WithFields(params).Error("failed to request butter swap")
//		return ret, resp.CodeInternalServerError
//	}
//	ret = &entity.SwapResponse{
//		To:      txData.To,
//		Data:    txData.Data,
//		Value:   txData.Value,
//		ChainId: txData.ChainId,
//	}
//	return ret, resp.CodeSuccess
//}

func Swap(srcChain, srcToken, sender string, amount *big.Int, dstToken, receiver, hash string, slippage uint64) (ret *entity.SwapResponse, code int) {
	order := &dao.Order{
		SrcChain: srcChain,
		SrcToken: srcToken,
		Sender:   sender,
		DstChain: constants.TONChainID,
		DstToken: dstToken,
		Receiver: receiver,
		Action:   dao.OrderActionFromEVM,
		Stage:    dao.OrderStag1,
		Status:   dao.OrderStatusPending,
		Slippage: slippage,
	}
	orderID, err := order.Create()
	if err != nil {
		log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
		return nil, resp.CodeInternalServerError
	}

	orderIDByte32 := utils.Uint64ToByte32(orderID)
	// todo check on the amount in the contract?
	packed, err := PackOnReceived(amount, orderIDByte32, common.HexToAddress(constants.USDTOfChainPoll), common.HexToAddress(sender), []byte(receiver))
	if err != nil {
		params := map[string]interface{}{
			"amount":        amount,
			"orderID":       order,
			"orderIDByte32": orderIDByte32,
			"token":         constants.USDTOfChainPoll,
			"sender":        sender,
			"receiver":      receiver,
			"error":         err,
		}
		log.Logger().WithFields(params).Error("failed to pack onReceived")
		return ret, resp.CodeInternalServerError
	}
	encodedCallback, err := EncodeSwapCallbackParams(common.HexToAddress(viper.GetString("feRouterContract")), common.HexToAddress(sender), packed) // todo sender
	if err != nil {
		params := map[string]interface{}{
			"feRouter": viper.GetString("feRouterContract"),
			"sender":   sender,
			"packed":   hex.EncodeToString(packed),
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to encode swap callback params")
		return ret, resp.CodeInternalServerError
	}

	request := &butter.SwapRequest{
		Hash:     hash,
		Slippage: slippage / 3 * 2,
		From:     sender,
		Receiver: viper.GetString("butterRouterContract"),
		CallData: encodedCallback,
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

func getTONRoutes(tonRequest *tonrouter.RouteRequest, routes []*butter.RouteResponseData) (map[string]*tonrouter.RouteData, error) {
	if len(routes) == 0 {
		return make(map[string]*tonrouter.RouteData), nil
	}

	var wg sync.WaitGroup
	result := sync.Map{}
	errChan := make(chan error, len(routes))

	for _, r := range routes {
		if r == nil {
			continue
		}

		wg.Add(1)
		go func(hash, amount string, request *tonrouter.RouteRequest) {
			defer wg.Done()

			request.Amount = amount
			tonRoute, err := tonrouter.Route(request)
			if err != nil {
				params := map[string]interface{}{
					"request": utils.JSON(request),
					"error":   err,
				}
				log.Logger().WithFields(params).Error("failed to request ton route")
				errChan <- err
				return
			}

			result.Store(hash, tonRoute)
		}(r.Hash, r.DstChain.TotalAmountOut, tonRequest)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	finalResult := make(map[string]*tonrouter.RouteData, len(routes))
	result.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(*tonrouter.RouteData)
		finalResult[k] = v
		return true
	})

	return finalResult, nil
}
