package logic

import (
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/mapprotocol/fe-backend/bindings/erc20"
	"github.com/mapprotocol/fe-backend/constants"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/third-party/butter"
	"github.com/mapprotocol/fe-backend/third-party/tonrouter"
	"github.com/mapprotocol/fe-backend/utils"
	"github.com/mapprotocol/fe-backend/utils/reqerror"
	"github.com/shopspring/decimal"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

var isMultiChainPool bool
var feRouterContract string
var feRouterAddress common.Address

var BridgeFeeRate = decimal.NewFromFloat(30.0 / 10000.0) // 30/10000

const BaseTxFeeMultiplier = 1.5

var ToTONBaseTxFee = decimal.NewFromFloat(1 * BaseTxFeeMultiplier)               // unit: USDT
var TONToEVMBaseTxFee = decimal.NewFromFloat(1)                                  // unit: USDT
var BitcoinToEVMBaseTxFee = decimal.NewFromFloat(0.000007 * BaseTxFeeMultiplier) // unit: WBTC
var BitcoinTxBytes = decimal.NewFromFloat(200)

func Init() {
	isMultiChainPool = viper.GetBool("isMultiChainPool")

	feRouterContract = viper.GetString("feRouterContract")
	if utils.IsEmpty(feRouterContract) {
		panic("feRouterContract is empty")
	}
	feRouterAddress = common.HexToAddress(feRouterContract)
}

var (
	USDTLimit = big.NewFloat(5)
	WBTCLimit = big.NewFloat(0.0005)
)

func GetTONToEVMRoute(req *entity.RouteRequest, amount decimal.Decimal, feeRatio, slippage uint64) (ret []*entity.RouteResponse, msg string, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	tonRequest := &tonrouter.BridgeRouteRequest{
		ToChainID:       req.ToChainID,
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: req.TokenOutAddress,
		Amount:          req.Amount,
		TonSlippage:     slippage / 3,
		Slippage:        slippage,
		FeeCollector:    req.FeeCollector,
		FeeRatio:        req.FeeRatio,
	}
	tonRoute, err := tonrouter.BridgeRoute(tonRequest)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(tonRequest),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request ton route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeTONRouteServerError
	}

	bridgeAmountIn, err := decimal.NewFromString(tonRoute.SrcChain.TokenAmountOut)
	if err != nil {
		params := map[string]interface{}{
			"tokenAmountOut": tonRoute.SrcChain.TokenAmountOut,
			"error":          err,
		}
		log.Logger().WithFields(params).Error("failed to parse token amount out to decimal")
		return ret, "", resp.CodeInternalServerError
	}

	bridgeFees := calcTONToEVMBridgeFees(bridgeAmountIn, BridgeFeeRate)
	if bridgeAmountIn.Compare(bridgeFees) != 1 {
		return ret, "", resp.CodeAmountTooFew
	}
	//bridgeFees := calcToEVMBridgeFees(bridgeAmountIn, BridgeFeeRate)

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
		Icon:     out.Image,
	}

	gasFee := entity.Fee{
		Amount: tonRoute.GasFee.Amount,
		Symbol: tonRoute.GasFee.Symbol,
	}
	bridgeFee := entity.Fee{
		Amount: bridgeFees.StringFixedBank(6),
		Symbol: constants.USDTSymbol,
	}
	protocolFee := entity.Fee{
		Amount: tonRoute.ProtocolFee.Amount,
		Symbol: tonRoute.ProtocolFee.Symbol,
	}

	request := &butter.RouteRequest{
		TokenInAddress:  constants.USDTOfChainPool,
		TokenOutAddress: req.TokenOutAddress,
		Type:            req.Type,
		Slippage:        slippage,
		FromChainID:     constants.ChainIDOfChainPool,
		ToChainID:       req.ToChainID,
		Amount:          bridgeAmountIn.Sub(bridgeFees).String(),
	}

	butterRoutes, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}

	r := butterRoutes[0]
	ret = make([]*entity.RouteResponse, 0, len(butterRoutes))
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
		Hash:      tonRoute.Hash,
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
		GasFee:      gasFee,
		BridgeFee:   bridgeFee,
		ProtocolFee: protocolFee,
	}
	ret = append(ret, n)
	return ret, "", resp.CodeSuccess
}

func GetEVMToTONRoute(req *entity.RouteRequest, amount decimal.Decimal, feeRatio, slippage uint64) (ret []*entity.RouteResponse, msg string, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	request := &butter.RouteRequest{
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: constants.USDTOfChainPool,
		Type:            req.Type,
		Slippage:        slippage / 3 * 2,
		FromChainID:     req.FromChainID,
		ToChainID:       constants.ChainIDOfChainPool,
		Amount:          req.Amount,
		Referrer:        req.FeeCollector,
		RateOrNativeFee: req.FeeRatio,
	}
	butterRoutes, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}
	if len(butterRoutes) == 0 {
		return ret, "", resp.CodeButterNotAvailableRoute
	}

	r := butterRoutes[0]
	bridgeAmountIn, err := decimal.NewFromString(r.DstChain.TotalAmountOut)
	if err != nil {
		params := map[string]interface{}{
			"tokenAmountOut": r.SrcChain.TotalAmountOut,
			"error":          err,
		}
		log.Logger().WithFields(params).Error("failed to parse token amount out to decimal")
		return ret, "", resp.CodeInternalServerError
	}
	bridgeFees := calcToTONBridgeFees(bridgeAmountIn, BridgeFeeRate)
	if bridgeAmountIn.Compare(bridgeFees) != 1 {
		return ret, "", resp.CodeAmountTooFew
	}

	protocolFees := "0"
	if feeRatio > 0 {
		protocolFees = calcProtocolFees(amount, decimal.NewFromFloat(float64(feeRatio)/10000.0)).StringFixedBank(int32(r.SrcChain.TokenIn.Decimals))
	}

	tonRequest := &tonrouter.RouteRequest{
		TokenInAddress:  constants.USDTOfTON,
		TokenOutAddress: req.TokenOutAddress,
		Amount:          bridgeAmountIn.Sub(bridgeFees).String(),
		Slippage:        slippage,
	}
	tonRoute, err := tonrouter.Route(tonRequest)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request ton route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
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
		Icon:     out.Image,
	}

	ret = make([]*entity.RouteResponse, 0, 1)
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
				Name:      r.SrcChain.Bridge,
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
			Amount: bridgeFees.StringFixedBank(6),
			Symbol: constants.USDTSymbol,
		},
		ProtocolFee: entity.Fee{
			Amount: protocolFees,
			Symbol: r.SrcChain.TokenIn.Symbol,
		},
	}
	ret = append(ret, n)

	return ret, "", resp.CodeSuccess
}

func GetBitcoinToEVMRoute(req *entity.RouteRequest, amount decimal.Decimal, feeRatio, slippage uint64) (ret []*entity.RouteResponse, msg string, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	protocolFees := calcProtocolFees(amount, decimal.NewFromFloat(float64(feeRatio)/10000.0))
	//bitcoinRoute := GetBitcoinLocalRoutes(amount.String())
	bitcoinRoute := GetBitcoinLocalRoutes(amount.Sub(protocolFees).String())

	bridgeAmountIn, err := decimal.NewFromString(bitcoinRoute.SrcChain.TokenAmountOut)
	if err != nil {
		params := map[string]interface{}{
			"tokenAmountOut": bitcoinRoute.SrcChain.TokenAmountOut,
			"error":          err,
		}
		log.Logger().WithFields(params).Error("failed to parse token amount out to decimal")
		return ret, "", resp.CodeInternalServerError
	}

	bridgeFees := calcBitcoinToEVMBridgeFees(bridgeAmountIn, BridgeFeeRate)
	if bridgeAmountIn.Compare(bridgeFees) != 1 {
		return ret, "", resp.CodeAmountTooFew
	}
	//bridgeFees := calcToEVMBridgeFees(bridgeAmountIn, BridgeFeeRate)

	in := bitcoinRoute.SrcChain.Route[0].Path[0].TokenIn
	tonTokenIn = entity.Token{
		ChainId:  bitcoinRoute.SrcChain.ChainId,
		Address:  in.Address,
		Name:     in.Name,
		Decimals: in.Decimals,
		Symbol:   in.Symbol,
		Icon:     in.Image,
	}

	out := bitcoinRoute.SrcChain.Route[0].Path[0].TokenOut
	tonTokenOut = entity.Token{
		ChainId:  bitcoinRoute.SrcChain.ChainId,
		Address:  out.Address,
		Name:     out.Name,
		Decimals: out.Decimals,
		Symbol:   out.Symbol,
		Icon:     out.Image,
	}

	gasFee := entity.Fee{
		Amount: bitcoinRoute.GasFee.Amount,
		Symbol: bitcoinRoute.GasFee.Symbol,
	}
	bridgeFee := entity.Fee{
		Amount: bridgeFees.StringFixedBank(8),
		Symbol: constants.WBTCSymbol,
	}
	protocolFee := entity.Fee{
		Amount: protocolFees.StringFixedBank(8),
		Symbol: constants.WBTCSymbol,
	}

	request := &butter.RouteRequest{
		TokenInAddress:  constants.WBTCOfChainPool,
		TokenOutAddress: req.TokenOutAddress,
		Type:            req.Type,
		Slippage:        slippage,
		FromChainID:     constants.ChainIDOfChainPool,
		ToChainID:       req.ToChainID,
		Amount:          bridgeAmountIn.Sub(bridgeFees).String(),
	}

	butterRoutes, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}

	ret = make([]*entity.RouteResponse, 0, 1)
	r := butterRoutes[0]
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
		Hash:      bitcoinRoute.Hash,
		TokenIn:   tonTokenIn,
		TokenOut:  butterDstChainTokenOut,
		AmountIn:  amount.String(), // todo
		AmountOut: r.DstChain.TotalAmountOut,
		Path: []entity.Path{
			{
				Name:      bitcoinRoute.SrcChain.Route[0].DexName,
				AmountIn:  amount.String(),
				AmountOut: bitcoinRoute.SrcChain.TokenAmountOut,
				TokenIn:   tonTokenIn,
				TokenOut:  tonTokenOut,
			},
			{
				Name:      constants.ExchangeNameFlushExchange,
				AmountIn:  bitcoinRoute.SrcChain.TokenAmountOut,
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
		GasFee:      gasFee,
		BridgeFee:   bridgeFee,
		ProtocolFee: protocolFee,
	}
	ret = append(ret, n)

	return ret, "", resp.CodeSuccess
}

func GetEVMToBitcoinRoute(req *entity.RouteRequest, amount decimal.Decimal, feeRatio, slippage uint64) (ret []*entity.RouteResponse, msg string, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	request := &butter.RouteRequest{
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: constants.WBTCOfChainPool,
		Type:            req.Type,
		Slippage:        slippage / 3 * 2,
		FromChainID:     req.FromChainID,
		ToChainID:       constants.ChainIDOfChainPool,
		Amount:          req.Amount,
		Referrer:        req.FeeCollector,
		RateOrNativeFee: req.FeeRatio,
	}
	butterRoutes, err := butter.Route(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}
	if len(butterRoutes) == 0 {
		return ret, "", resp.CodeButterNotAvailableRoute
	}

	r := butterRoutes[0]
	amountDecimal, err := decimal.NewFromString(r.DstChain.TotalAmountOut)
	if err != nil {
		params := map[string]interface{}{
			"amount": r.DstChain.TotalAmountOut,
			"error":  err,
		}
		log.Logger().WithFields(params).Error("failed to parse total amount out to decimal")
		return ret, "", resp.CodeInternalServerError
	}
	bridgeFees := calcToBitcoinBridgeFees(amountDecimal, BridgeFeeRate)
	if amountDecimal.Compare(bridgeFees) != 1 {
		return ret, "", resp.CodeAmountTooFew
	}

	protocolFees := "0"
	if feeRatio > 0 {
		protocolFees = calcProtocolFees(amount, decimal.NewFromFloat(float64(feeRatio)/10000.0)).StringFixedBank(8)
	}

	bitcoinRoute := GetBitcoinLocalRoutes(amountDecimal.StringFixedBank(8))
	ret = make([]*entity.RouteResponse, 0, len(butterRoutes))
	in := bitcoinRoute.SrcChain.Route[0].Path[0].TokenIn
	tonTokenIn = entity.Token{
		ChainId:  bitcoinRoute.SrcChain.ChainId,
		Address:  in.Address,
		Name:     in.Name,
		Decimals: in.Decimals,
		Symbol:   in.Symbol,
		Icon:     in.Image,
	}

	out := bitcoinRoute.SrcChain.Route[0].Path[0].TokenOut
	tonTokenOut = entity.Token{
		ChainId:  bitcoinRoute.SrcChain.ChainId,
		Address:  out.Address,
		Name:     out.Name,
		Decimals: out.Decimals,
		Symbol:   out.Symbol,
		Icon:     out.Image,
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
		AmountIn:  amount.String(),
		AmountOut: bitcoinRoute.SrcChain.TokenAmountOut,
		Path: []entity.Path{
			{
				Name:      r.SrcChain.Bridge,
				AmountIn:  amount.String(),
				AmountOut: r.SrcChain.TotalAmountOut,
				TokenIn:   butterSrcChainTokenIn,
				TokenOut:  butterDstChainTokenOut,
			},
			{
				Name:      constants.ExchangeNameFlushExchange,
				AmountIn:  r.DstChain.TotalAmountOut,
				AmountOut: bitcoinRoute.SrcChain.TokenAmountIn,
				TokenIn:   butterDstChainTokenOut,
				TokenOut:  tonTokenIn,
			},
			{
				Name:      bitcoinRoute.SrcChain.Route[0].DexName,
				AmountIn:  bitcoinRoute.SrcChain.TokenAmountIn,
				AmountOut: bitcoinRoute.SrcChain.TokenAmountOut,
				TokenIn:   tonTokenIn,
				TokenOut:  tonTokenOut,
			},
		},
		GasFee: entity.Fee{
			Amount: r.GasFee.Amount,
			Symbol: r.GasFee.Symbol,
		},
		BridgeFee: entity.Fee{
			Amount: bridgeFees.StringFixedBank(8),
			Symbol: constants.WBTCSymbol,
		},
		ProtocolFee: entity.Fee{
			Amount: protocolFees,
			Symbol: r.SrcChain.TokenIn.Symbol,
		},
	}
	ret = append(ret, n)

	return ret, "", resp.CodeSuccess
}

func GetSwapFromTONToEVM(sender string, dstChain, receiver, feeCollector, feeRatio, hash string) (ret *entity.SwapResponse, msg string, code int) {
	amountOut, err := tonrouter.GetRouteAmountOut(hash)
	if err != nil {
		params := map[string]interface{}{
			"hash":  hash,
			"error": err,
		}
		log.Logger().WithFields(params).Error("failed to request ton get route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}
	if amountOut.Cmp(USDTLimit) == -1 {
		return ret, "", resp.CodeAmountTooFew
	}

	balance, err := getChainPoolUSDTBalance(dstChain)
	if err != nil {
		params := map[string]interface{}{
			"dstChain": dstChain,
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to get ton router balance")
		return ret, "", resp.CodeInternalServerError
	}
	if amountOut.Cmp(balance) == 1 {
		params := map[string]interface{}{
			"amount":  amountOut.Text('f', -1),
			"balance": balance.Text('f', -1),
		}
		log.Logger().WithFields(params).Info("amount is greater than ton router balance")
		return nil, "", resp.CodeInsufficientLiquidity // todo chain pool
	}

	request := &tonrouter.BridgeSwapRequest{
		Sender:       sender,
		Receiver:     receiver,
		FeeCollector: feeCollector,
		FeeRatio:     feeRatio,
		Hash:         hash,
	}
	txData, err := tonrouter.BridgeSwap(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request ton bridge swap")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}
	ret = &entity.SwapResponse{
		To:      txData.To,
		Data:    txData.Data,
		Value:   txData.Value,
		ChainId: constants.TONChainID,
	}
	return ret, "", resp.CodeSuccess
}
func GetSwapFromEVMToTON(srcChain *big.Int, srcToken, sender, amount string, dstChain *big.Int, dstToken, receiver, hash string, slippage uint64) (ret *entity.SwapResponse, msg string, code int) {
	amountOut, err := butter.GetRouteAmountOut(hash)
	if err != nil {
		params := map[string]interface{}{
			"hash":  hash,
			"error": err,
		}
		log.Logger().WithFields(params).Error("failed to request butter get route")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return nil, "", resp.CodeInternalServerError
	}
	if amountOut.Cmp(USDTLimit) == -1 {
		return ret, "", resp.CodeAmountTooFew
	}
	balance, err := tonrouter.Balance()
	if err != nil {
		log.Logger().WithField("error", err).Error("failed to get ton router balance")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return nil, "", resp.CodeInternalServerError
	}
	if amountOut.Cmp(balance) == 1 {
		params := map[string]interface{}{
			"amount":  amountOut.Text('f', -1),
			"balance": balance.Text('f', -1),
		}
		log.Logger().WithFields(params).Info("amount is greater than balance")
		return nil, "", resp.CodeInsufficientLiquidity // todo chain pool
	}

	chainPoolToken := constants.USDTOfChainPool
	if isMultiChainPool && dstChain.String() == constants.ChainIDOfEthereum {
		chainPoolToken = constants.USDTOfEthereum
	}
	params := ReceiverParam{
		OrderId:        [32]byte{},
		SrcChain:       srcChain,
		SrcToken:       []byte(srcToken),
		Sender:         []byte(sender),
		InAmount:       amount,
		ChainPoolToken: common.HexToAddress(chainPoolToken),
		DstChain:       dstChain,
		DstToken:       []byte(dstToken),
		Receiver:       []byte(receiver),
		Slippage:       slippage,
	}
	packed, err := PackOnReceived(big.NewInt(0), params)
	if err != nil {
		params := map[string]interface{}{
			"params": utils.JSON(params),
			"error":  err,
		}
		log.Logger().WithFields(params).Error("failed to pack onReceived")
		return ret, "", resp.CodeInternalServerError
	}
	encodedCallback, err := EncodeSwapCallbackParams(feRouterAddress, common.HexToAddress(sender), packed) // todo sender
	if err != nil {
		params := map[string]interface{}{
			"feRouter": feRouterContract,
			"sender":   sender,
			"packed":   hex.EncodeToString(packed),
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to encode swap callback params")
		return ret, "", resp.CodeInternalServerError
	}

	request := &butter.SwapRequest{
		Hash:     hash,
		Slippage: slippage / 3 * 2,
		From:     sender,
		Receiver: sender, // todo sender
		CallData: encodedCallback,
	}
	txData, err := butter.Swap(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter swap")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}
	ret = &entity.SwapResponse{
		To:      txData.To,
		Data:    txData.Data,
		Value:   txData.Value,
		ChainId: txData.ChainId,
	}
	return ret, "", resp.CodeSuccess
}

func GetSwapFromBitcoinToEVM(srcChain, srcToken, sender string, amount *big.Float, amountBigInt *big.Int, dstChain, dstToken, receiver string, slippage uint64, feeCollector string, feeRatio uint64) (ret *entity.SwapResponse, msg string, code int) {
	if amount.Cmp(WBTCLimit) == -1 {
		return ret, "", resp.CodeAmountTooFew
	}

	balance, err := getChainPoolWBTCBalance(dstChain)
	if err != nil {
		params := map[string]interface{}{
			"dstChain": dstChain,
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to get wbtc balance")
		return ret, "", resp.CodeInternalServerError
	}
	if amount.Cmp(balance) == 1 {
		params := map[string]interface{}{
			"dstChain": dstChain,
			"amount":   amount.Text('f', -1),
			"balance":  balance.Text('f', -1),
		}
		log.Logger().WithFields(params).Info("amount is greater than chain pool balance")
		return nil, "", resp.CodeInsufficientLiquidity // todo chain pool
	}

	privateKey, err := generateKey()
	if err != nil {
		log.Logger().WithField("error", err).Error("failed to generate key")
		return ret, "", resp.CodeInternalServerError
	}
	address, err := makeTaprootAddress(privateKey, NetParams)
	if err != nil {
		log.Logger().WithField("error", err).Error("failed to make address")
		return ret, "", resp.CodeInternalServerError
	}

	order := &dao.BitcoinOrder{
		SrcChain:     srcChain,
		SrcToken:     srcToken,
		Sender:       sender,
		InAmount:     amount.Text('f', 8),
		InAmountSat:  amountBigInt.String(),
		Relayer:      address.String(),
		RelayerKey:   privateKey.Key.String(),
		DstChain:     dstChain,
		DstToken:     dstToken,
		Receiver:     receiver,
		Action:       dao.OrderActionToEVM,
		Stage:        dao.OrderStag1,
		Status:       dao.OrderStatusTxPrepareSend,
		Slippage:     slippage,
		FeeRatio:     feeRatio,
		FeeCollector: feeCollector,
	}
	if err := order.Create(); err != nil {
		log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
		return ret, "", resp.CodeInternalServerError
	}

	ret = &entity.SwapResponse{
		To:      address.String(),
		Data:    "",
		Value:   "0x" + amountBigInt.Text(16),
		ChainId: constants.BTCChainID,
	}
	return ret, "", resp.CodeSuccess
}

func GetSwapFromEVMToBitcoin(srcChain *big.Int, srcToken, sender, amount string, dstChain *big.Int, dstToken, receiver, hash string, slippage uint64) (ret *entity.SwapResponse, msg string, code int) {
	amountOut, err := butter.GetRouteAmountOut(hash)
	if err != nil {
		params := map[string]interface{}{
			"hash":  hash,
			"error": err,
		}
		log.Logger().WithFields(params).Error("failed to request butter get route")
		//ext, ok := err.(*reqerror.ExternalRequestError)
		var ext *reqerror.ExternalRequestError
		if ok := errors.As(err, &ext); ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return nil, "", resp.CodeInternalServerError
	}
	if amountOut.Cmp(WBTCLimit) == -1 {
		return ret, "", resp.CodeAmountTooFew
	}

	balanceSats, err := btcApiClient.Balance(btcVaultAddress)
	if err != nil {
		log.Logger().WithField("error", err).Error("failed to get bitcoin vault address balance")
		return nil, "", resp.CodeInternalServerError
	}
	balance := big.NewFloat(btcutil.Amount(balanceSats).ToBTC())
	if amountOut.Cmp(balance) == 1 {
		params := map[string]interface{}{
			"amount":  amountOut.Text('f', -1),
			"balance": balance.Text('f', -1),
		}
		log.Logger().WithFields(params).Info("amount is greater than bitcoin vault address balance")
		return nil, "", resp.CodeInsufficientLiquidity
	}

	chainPoolToken := constants.WBTCOfChainPool
	if isMultiChainPool && dstChain.String() == constants.ChainIDOfEthereum {
		chainPoolToken = constants.WBTCOfEthereum
	}
	params := ReceiverParam{
		OrderId:        [32]byte{},
		SrcChain:       srcChain,
		SrcToken:       []byte(srcToken),
		Sender:         []byte(sender),
		InAmount:       amount,
		ChainPoolToken: common.HexToAddress(chainPoolToken),
		DstChain:       dstChain,
		DstToken:       []byte(dstToken),
		Receiver:       []byte(receiver),
		Slippage:       slippage,
	}
	packed, err := PackOnReceived(big.NewInt(0), params)
	if err != nil {
		params := map[string]interface{}{
			"params": utils.JSON(params),
			"error":  err,
		}
		log.Logger().WithFields(params).Error("failed to pack onReceived")
		return ret, "", resp.CodeInternalServerError
	}
	encodedCallback, err := EncodeSwapCallbackParams(feRouterAddress, common.HexToAddress(sender), packed) // todo sender
	if err != nil {
		params := map[string]interface{}{
			"feRouter": feRouterContract,
			"sender":   sender,
			"packed":   hex.EncodeToString(packed),
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to encode swap callback params")
		return ret, "", resp.CodeInternalServerError
	}

	request := &butter.SwapRequest{
		Hash:     hash,
		Slippage: slippage / 3 * 2,
		From:     sender,
		Receiver: sender,
		CallData: encodedCallback,
	}
	txData, err := butter.Swap(request)
	if err != nil {
		params := map[string]interface{}{
			"request": utils.JSON(request),
			"error":   err,
		}
		log.Logger().WithFields(params).Error("failed to request butter swap")
		ext, ok := err.(*reqerror.ExternalRequestError)
		if ok && ext.HasPublicError() {
			return nil, ext.PublicError(), resp.CodeExternalServerError
		}
		return ret, "", resp.CodeInternalServerError
	}
	ret = &entity.SwapResponse{
		To:      txData.To,
		Data:    txData.Data,
		Value:   txData.Value,
		ChainId: txData.ChainId,
	}
	return ret, "", resp.CodeSuccess
}

//func GetLocalRouteSwapFromEVMToTON(srcChain *big.Int, srcToken, sender, amount string, amountBigFloat *big.Float, amountBigInt, dstChain *big.Int, dstToken, receiver string, slippage uint64) (ret *entity.SwapResponse, msg string, code int) {
//	if amountBigFloat.Cmp(USDTLimit) == -1 {
//		return ret, "", resp.CodeAmountTooFew
//	}
//
//	chainPoolToken := constants.USDTOfChainPool
//	if isMultiChainPool && dstChain.String() == constants.ChainIDOfEthereum {
//		chainPoolToken = constants.USDTOfEthereum
//	}
//	params := ReceiverParam{
//		OrderId:        [32]byte{},
//		SrcChain:       srcChain,
//		SrcToken:       []byte(srcToken),
//		Sender:         []byte(sender),
//		InAmount:       amount,
//		ChainPoolToken: common.HexToAddress(chainPoolToken),
//		DstChain:       dstChain,
//		DstToken:       []byte(dstToken),
//		Receiver:       []byte(receiver),
//		Slippage:       slippage,
//	}
//	packed, err := PackOnReceived(amountBigInt, params)
//	if err != nil {
//		params := map[string]interface{}{
//			"params": utils.JSON(params),
//			"error":  err,
//		}
//		log.Logger().WithFields(params).Error("failed to pack onReceived")
//		return ret, "", resp.CodeInternalServerError
//	}
//	ret = &entity.SwapResponse{
//		To:      feRouterContract,
//		Data:    "0x" + hex.EncodeToString(packed),
//		Value:   "0x" + amountBigInt.Text(16),
//		ChainId: constants.ChainIDOfChainPool,
//	}
//	return ret, "", resp.CodeSuccess
//}
//
//func GetLocalRouteSwapFromEVMToBitcoin(srcChain *big.Int, srcToken, sender, amount string, amountBigFloat *big.Float, amountBigInt, dstChain *big.Int, dstToken, receiver string, slippage uint64) (ret *entity.SwapResponse, msg string, code int) {
//	if amountBigFloat.Cmp(WBTCLimit) == -1 {
//		return ret, "", resp.CodeAmountTooFew
//	}
//
//	chainPoolToken := constants.WBTCOfChainPool
//	if isMultiChainPool && dstChain.String() == constants.ChainIDOfEthereum {
//		chainPoolToken = constants.WBTCOfEthereum
//	}
//	params := ReceiverParam{
//		OrderId:        [32]byte{},
//		SrcChain:       srcChain,
//		SrcToken:       []byte(srcToken),
//		Sender:         []byte(sender),
//		InAmount:       amount,
//		ChainPoolToken: common.HexToAddress(chainPoolToken),
//		DstChain:       dstChain,
//		DstToken:       []byte(dstToken),
//		Receiver:       []byte(receiver),
//		Slippage:       slippage,
//	}
//	packed, err := PackOnReceived(amountBigInt, params)
//	if err != nil {
//		params := map[string]interface{}{
//			"params": utils.JSON(params),
//			"error":  err,
//		}
//		log.Logger().WithFields(params).Error("failed to pack onReceived")
//		return ret, "", resp.CodeInternalServerError
//	}
//	ret = &entity.SwapResponse{
//		To:      feRouterContract,
//		Data:    "0x" + hex.EncodeToString(packed),
//		Value:   "0x" + amountBigInt.Text(16),
//		ChainId: constants.ChainIDOfChainPool,
//	}
//	return ret, "", resp.CodeSuccess
//}
//
//func getToTONRoutes(tonRequest *tonrouter.RouteRequest, routes []*butter.RouteResponseData) (map[string]*tonrouter.RouteData, error) {
//	if len(routes) == 0 {
//		return make(map[string]*tonrouter.RouteData), nil
//	}
//
//	var wg sync.WaitGroup
//	result := sync.Map{}
//	errChan := make(chan error, len(routes))
//
//	for _, r := range routes {
//		if r == nil {
//			continue
//		}
//
//		wg.Add(1)
//		go func(hash, amount string, request *tonrouter.RouteRequest) {
//			defer wg.Done()
//
//			amountDecimal, err := decimal.NewFromString(amount)
//			if err != nil {
//				params := map[string]interface{}{
//					"amount": amount,
//					"error":  err,
//				}
//				log.Logger().WithFields(params).Error("failed to parse amount to decimal")
//				errChan <- err
//				return
//			}
//
//			bridgeFees := calcToTONBridgeFees(amountDecimal, BridgeFeeRate)
//			request.Amount = amountDecimal.Sub(bridgeFees).String()
//			tonRoute, err := tonrouter.Route(request)
//			if err != nil {
//				params := map[string]interface{}{
//					"request": utils.JSON(request),
//					"error":   err,
//				}
//				log.Logger().WithFields(params).Error("failed to request ton route")
//				errChan <- err
//				return
//			}
//
//			result.Store(hash, tonRoute)
//		}(r.Hash, r.DstChain.TotalAmountOut, tonRequest)
//	}
//
//	wg.Wait()
//	close(errChan)
//
//	if len(errChan) > 0 {
//		return nil, <-errChan
//	}
//
//	finalResult := make(map[string]*tonrouter.RouteData, len(routes))
//	result.Range(func(key, value interface{}) bool {
//		k := key.(string)
//		v := value.(*tonrouter.RouteData)
//		finalResult[k] = v
//		return true
//	})
//
//	return finalResult, nil
//}
//
//func getToBitcoinRoutes(routes []*butter.RouteResponseData) (map[string]*Route, error) {
//	if len(routes) == 0 {
//		return make(map[string]*Route), nil
//	}
//	finalResult := make(map[string]*Route, len(routes))
//
//	for _, r := range routes {
//		if r == nil {
//			continue
//		}
//
//		amountDecimal, err := decimal.NewFromString(r.DstChain.TotalAmountOut)
//		if err != nil {
//			params := map[string]interface{}{
//				"amount": r.DstChain.TotalAmountOut,
//				"error":  err,
//			}
//			log.Logger().WithFields(params).Error("failed to parse amount to decimal")
//			return nil, err
//		}
//
//		bridgeFees := calcToBitcoinBridgeFees(amountDecimal, BridgeFeeRate)
//		route := GetBitcoinLocalRoutes(amountDecimal.Sub(bridgeFees).String())
//		finalResult[r.Hash] = route
//	}
//
//	return finalResult, nil
//}

func getChainPoolUSDTBalance(dstChain string) (balance *big.Float, err error) {
	chainInfo := &dao.ChainPool{}
	if isMultiChainPool && dstChain == constants.ChainIDOfEthereum {
		chainInfo, err = dao.NewChainPoolWithChainID(constants.ChainIDOfEthereum).First()
		if err != nil {
			log.Logger().WithField("chainID", constants.ChainIDOfEthereum).WithField("error", err.Error()).Error("failed to get chain info")
			return balance, err

		}
	} else {
		chainInfo, err = dao.NewChainPoolWithChainID(constants.ChainIDOfChainPool).First()
		if err != nil {
			log.Logger().WithField("chainID", constants.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
			return balance, err
		}
	}

	cli, err := ethclient.Dial(chainInfo.ChainRPC)
	if err != nil {
		params := map[string]interface{}{
			"chainID":  chainInfo.ChainID,
			"chainRPC": chainInfo.ChainRPC,
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to dial chain rpc")
		return balance, err
	}

	caller, err := erc20.NewErc20Caller(common.HexToAddress(chainInfo.USDTContract), cli)
	if err != nil {
		params := map[string]interface{}{
			"chainID":      chainInfo.ChainID,
			"USDTContract": chainInfo.USDTContract,
			"error":        err,
		}
		log.Logger().WithFields(params).Error("failed to new erc20 caller")
		return balance, err
	}

	bal, err := caller.BalanceOf(nil, common.HexToAddress(chainInfo.ChainPoolContract))
	if err != nil {
		params := map[string]interface{}{
			"chainID":      chainInfo.ChainID,
			"USDTContract": chainInfo.USDTContract,
			"account":      chainInfo.ChainPoolContract,
			"error":        err,
		}
		log.Logger().WithFields(params).Error("failed to get balance of chain pool contract")
		return balance, err
	}
	balance = new(big.Float).Quo(new(big.Float).SetInt(bal), getUSDTDecimalOfChainPool(dstChain))
	return balance, err
}
func getChainPoolWBTCBalance(dstChain string) (balance *big.Float, err error) {
	chainInfo := &dao.ChainPool{}
	if isMultiChainPool && dstChain == constants.ChainIDOfEthereum {
		chainInfo, err = dao.NewChainPoolWithChainID(constants.ChainIDOfEthereum).First()
		if err != nil {
			log.Logger().WithField("chainID", constants.ChainIDOfEthereum).WithField("error", err.Error()).Error("failed to get chain info")
			return balance, err

		}
	} else {
		chainInfo, err = dao.NewChainPoolWithChainID(constants.ChainIDOfChainPool).First()
		if err != nil {
			log.Logger().WithField("chainID", constants.ChainIDOfChainPool).WithField("error", err.Error()).Error("failed to get chain info")
			return balance, err
		}
	}

	cli, err := ethclient.Dial(chainInfo.ChainRPC)
	if err != nil {
		params := map[string]interface{}{
			"chainID":  chainInfo.ChainID,
			"chainRPC": chainInfo.ChainRPC,
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to dial chain rpc")
		return balance, err
	}

	caller, err := erc20.NewErc20Caller(common.HexToAddress(chainInfo.WBTCContract), cli)
	if err != nil {
		params := map[string]interface{}{
			"chainID":      chainInfo.ChainID,
			"WBTCContract": chainInfo.WBTCContract,
			"error":        err,
		}
		log.Logger().WithFields(params).Error("failed to new erc20 caller")
		return balance, err
	}

	bal, err := caller.BalanceOf(nil, common.HexToAddress(chainInfo.ChainPoolContract))
	if err != nil {
		params := map[string]interface{}{
			"chainID":      chainInfo.ChainID,
			"WBTCContract": chainInfo.WBTCContract,
			"account":      chainInfo.ChainPoolContract,
			"error":        err,
		}
		log.Logger().WithFields(params).Error("failed to get balance of chain pool contract")
		return balance, err
	}
	balance = new(big.Float).Quo(new(big.Float).SetInt(bal), getWBTCDecimalOfChainPool(dstChain))
	return balance, err
}

func getUSDTDecimalOfChainPool(chain string) (decimal *big.Float) {
	switch chain {
	case constants.ChainIDOfEthereum:
		decimal = big.NewFloat(constants.USDTDecimalOfEthereum) // todo
	default:
		decimal = big.NewFloat(constants.USDTDecimalOfChainPool) // todo
	}
	return decimal
}

func getWBTCDecimalOfChainPool(chain string) (decimal *big.Float) {
	switch chain {
	case constants.ChainIDOfEthereum:
		decimal = big.NewFloat(constants.WBTCDecimalOfEthereum) // todo
	default:
		decimal = big.NewFloat(constants.WBTCDecimalOfChainPool) // todo
	}
	return decimal
}

//func calcBridgeFee(amount, feeRate *big.Float) (feeAmount *big.Float) {
//	feeAmount = new(big.Float).Mul(amount, feeRate)
//	feeAmount = new(big.Float).Quo(feeAmount, big.NewFloat(10000))
//	log.Logger().WithField("amount", amount).WithField("BridgeFeeRate", feeRate).WithField("feeAmount", feeAmount).Info("calc bridge fee")
//	return feeAmount
//}

func calcBridgeAndProtocolFees(amount, bridgeFeeRate, protocolFeeRate decimal.Decimal) (bridgeFeesStr, protocolFeesStr, afterAmountStr string) {
	//bridgeFees := amount.Mul(BridgeFeeRate).Div(decimal.New(1, 4))
	//protocolFees := amount.Mul(protocolFeeRate).Div(decimal.New(1, 4))
	//afterAmount := amount.Sub(bridgeFees).Sub(protocolFees)

	bridgeFees := amount.Mul(bridgeFeeRate)
	protocolFees := amount.Mul(protocolFeeRate)
	afterAmount := amount.Sub(bridgeFees).Sub(protocolFees)
	fields := map[string]interface{}{
		"amount":          amount,
		"bridgeFees":      bridgeFees,
		"protocolFees":    protocolFees,
		"afterAmount":     afterAmount,
		"BridgeFeeRate":   bridgeFeeRate,
		"protocolFeeRate": protocolFeeRate,
	}

	log.Logger().WithFields(fields).Info("calc bridge and protocol fees")
	return bridgeFees.String(), protocolFees.String(), afterAmount.String()
}

func calcToTONBridgeFees(amount, bridgeFeeRate decimal.Decimal) (bridgeFees decimal.Decimal) {
	bridgeFees = amount.Mul(bridgeFeeRate).Add(ToTONBaseTxFee)

	fields := map[string]interface{}{
		"amount":        amount,
		"BridgeFeeRate": bridgeFeeRate,
		"bridgeFees":    bridgeFees,
		"baseTxFee":     ToTONBaseTxFee,
	}
	log.Logger().WithFields(fields).Info("complete the calc to ton bridge fees")
	return bridgeFees
}

func calcToBitcoinBridgeFees(amount, bridgeFeeRate decimal.Decimal) (bridgeFees decimal.Decimal) {
	feeRate := GetGlobalFeeRate()
	baseTxFee := BitcoinTxBytes.Mul(decimal.NewFromInt(feeRate)).
		Mul(decimal.NewFromFloat(BaseTxFeeMultiplier)).
		Div(decimal.NewFromUint64(constants.BTCDecimal))
	bridgeFees = amount.Mul(bridgeFeeRate).Add(baseTxFee)

	fields := map[string]interface{}{
		"amount":        amount,
		"BridgeFeeRate": bridgeFeeRate,
		"bridgeFees":    bridgeFees,
		"btcFeeRate":    feeRate,
		"baseTxFee":     baseTxFee,
	}
	log.Logger().WithFields(fields).Info("completed the calc to bitcoin bridge fees")
	return bridgeFees
}

//func calcToEVMBridgeFees(amount, bridgeFeeRate decimal.Decimal) (bridgeFees decimal.Decimal) {
//	bridgeFees = amount.Mul(bridgeFeeRate).Add(EVMBaseTxFee)
//
//	fields := map[string]interface{}{
//		"amount":        amount,
//		"BridgeFeeRate": bridgeFeeRate,
//		"bridgeFees":    bridgeFees,
//	}
//	log.Logger().WithFields(fields).Info("complete the calc ton to evm bridge fees")
//	return bridgeFees
//}

func calcTONToEVMBridgeFees(amount, bridgeFeeRate decimal.Decimal) (bridgeFees decimal.Decimal) {
	bridgeFees = amount.Mul(bridgeFeeRate).Add(TONToEVMBaseTxFee)

	fields := map[string]interface{}{
		"amount":        amount,
		"BridgeFeeRate": bridgeFeeRate,
		"bridgeFees":    bridgeFees,
		"baseTxFee":     TONToEVMBaseTxFee,
	}
	log.Logger().WithFields(fields).Info("complete the calc ton to evm bridge fees")
	return bridgeFees
}

func calcBitcoinToEVMBridgeFees(amount, bridgeFeeRate decimal.Decimal) (bridgeFees decimal.Decimal) {
	bridgeFees = amount.Mul(bridgeFeeRate).Add(BitcoinToEVMBaseTxFee)

	fields := map[string]interface{}{
		"amount":        amount,
		"BridgeFeeRate": bridgeFeeRate,
		"bridgeFees":    bridgeFees,
		"baseTxFee":     BitcoinToEVMBaseTxFee,
	}
	log.Logger().WithFields(fields).Info("complete the calc bitcoin to evm bridge fees")
	return bridgeFees
}

func calcProtocolFees(amount, protocolFeeRate decimal.Decimal) (protocolFees decimal.Decimal) {
	protocolFees = amount.Mul(protocolFeeRate)

	fields := map[string]interface{}{
		"amount":          amount,
		"protocolFeeRate": protocolFeeRate,
		"protocolFees":    protocolFees,
	}
	log.Logger().WithFields(fields).Info("complete the calc protocol fees")
	return protocolFees
}

//func calcToTONBridgeAndProtocolFees(amount, BridgeFeeRate, protocolFeeRate decimal.Decimal) (bridgeFees, protocolFees decimal.Decimal) {
//	bridgeFees = amount.Mul(BridgeFeeRate).Add(ToTONBaseTxFee)
//	protocolFees = amount.Mul(protocolFeeRate)
//
//	fields := map[string]interface{}{
//		"amount":          amount,
//		"BridgeFeeRate":   BridgeFeeRate,
//		"bridgeFees":      bridgeFees,
//		"baseTxFee":       ToTONBaseTxFee,
//		"protocolFees":    protocolFees,
//		"protocolFeeRate": protocolFeeRate,
//	}
//	log.Logger().WithFields(fields).Info("complete the calc to ton bridge and protocol fees")
//	return bridgeFees, protocolFees
//}
//
//func calcToBitcoinBridgeAndProtocolFees(amount, BridgeFeeRate, protocolFeeRate decimal.Decimal) (bridgeFees, protocolFees decimal.Decimal) {
//	feeRate := GetGlobalFeeRate()
//	baseTxFee := BitcoinTxBytes.Mul(decimal.NewFromInt(feeRate)).
//		Mul(decimal.NewFromFloat(BaseTxFeeMultiplier)).
//		Div(decimal.NewFromUint64(constants.BTCDecimal))
//	bridgeFees = amount.Mul(BridgeFeeRate).Sub(baseTxFee)
//	protocolFees = amount.Mul(protocolFeeRate)
//
//	fields := map[string]interface{}{
//		"amount":          amount,
//		"BridgeFeeRate":   BridgeFeeRate,
//		"bridgeFees":      bridgeFees,
//		"btcFeeRate":      feeRate,
//		"baseTxFee":       baseTxFee,
//		"protocolFees":    protocolFees,
//		"protocolFeeRate": protocolFeeRate,
//	}
//	log.Logger().WithFields(fields).Info("completed the calc to bitcoin bridge and protocol fees")
//	return bridgeFees, protocolFees
//}
//
//func calcTONToEVMBridgeAndProtocolFees(amount, BridgeFeeRate, protocolFeeRate decimal.Decimal) (bridgeFees, protocolFees decimal.Decimal) {
//	bridgeFees = amount.Mul(BridgeFeeRate).Add(TONToEVMBaseTxFee)
//	protocolFees = amount.Mul(protocolFeeRate)
//
//	fields := map[string]interface{}{
//		"amount":          amount,
//		"BridgeFeeRate":   BridgeFeeRate,
//		"bridgeFees":      bridgeFees,
//		"baseTxFee":       TONToEVMBaseTxFee,
//		"protocolFees":    protocolFees,
//		"protocolFeeRate": protocolFeeRate,
//	}
//	log.Logger().WithFields(fields).Info("complete the calc ton to evm bridge and protocol fees")
//	return bridgeFees, protocolFees
//}
//
//func calcBitcoinToEVMBridgeAndProtocolFees(amount, BridgeFeeRate, protocolFeeRate decimal.Decimal) (bridgeFees, protocolFees decimal.Decimal) {
//	bridgeFees = amount.Mul(BridgeFeeRate).Add(BitcoinToEVMBaseTxFee)
//	protocolFees = amount.Mul(protocolFeeRate)
//
//	fields := map[string]interface{}{
//		"amount":          amount,
//		"BridgeFeeRate":   BridgeFeeRate,
//		"bridgeFees":      bridgeFees,
//		"baseTxFee":       BitcoinToEVMBaseTxFee,
//		"protocolFees":    protocolFees,
//		"protocolFeeRate": protocolFeeRate,
//	}
//	log.Logger().WithFields(fields).Info("complete the calc bitcoin to evm bridge and protocol fees")
//	return bridgeFees, protocolFees
//}

func calcBridgeAndProtocolFees1(amount, bridgeFeeRate, protocolFeeRate *big.Float) (bridgeFees, protocolFees, afterAmount *big.Float) {
	bridgeFees = new(big.Float).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Float).Quo(bridgeFees, big.NewFloat(10000))
	afterAmount = new(big.Float).Sub(amount, bridgeFees)

	protocolFees = new(big.Float).Mul(amount, protocolFeeRate)
	protocolFees = new(big.Float).Quo(protocolFees, big.NewFloat(10000))
	afterAmount = new(big.Float).Sub(afterAmount, protocolFees)

	fields := map[string]interface{}{
		"amount":          amount,
		"bridgeFees":      bridgeFees,
		"protocolFees":    protocolFees,
		"afterAmount":     afterAmount,
		"BridgeFeeRate":   bridgeFeeRate,
		"protocolFeeRate": protocolFeeRate,
	}

	log.Logger().WithFields(fields).Info("calc bridge and protocol fees")
	return bridgeFees, protocolFees, afterAmount
}

func calcBridgeAndProtocolFees2(amount, bridgeFeeRate, protocolFeeRate *big.Rat) (bridgeFeesStr, protocolFeesStr, afterAmountStr string) {
	bridgeFees := new(big.Rat).Mul(amount, bridgeFeeRate)
	bridgeFees = new(big.Rat).Quo(bridgeFees, new(big.Rat).SetUint64(10000))
	afterAmount := new(big.Rat).Sub(amount, bridgeFees)

	protocolFees := new(big.Rat).Mul(amount, protocolFeeRate)
	protocolFees = new(big.Rat).Quo(protocolFees, new(big.Rat).SetUint64(10000))
	afterAmount = new(big.Rat).Sub(afterAmount, protocolFees)

	fields := map[string]interface{}{
		"amount":          amount,
		"bridgeFees":      bridgeFees,
		"protocolFees":    protocolFees,
		"afterAmount":     afterAmount,
		"BridgeFeeRate":   bridgeFeeRate,
		"protocolFeeRate": protocolFeeRate,
	}

	log.Logger().WithFields(fields).Info("calc bridge and protocol fees")
	return bridgeFees.FloatString(-1), protocolFees.FloatString(-1), afterAmount.FloatString(-1)
}
