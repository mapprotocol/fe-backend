package logic

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	"github.com/spf13/viper"
	"math/big"
	"sync"
)

var isMultiChainPool = false
var feeRate = big.NewFloat(70) // 70/10000

func init() {
	isMultiChainPool = viper.GetBool("isMultiChainPool")
}

var BTCToken = entity.Token{
	ChainId:  constants.BTCChainID,
	Address:  constants.BTCTokenAddress,
	Name:     "BTC",
	Decimals: 8,
	Symbol:   "BTC",
}

func GetTONToEVMRoute(req *entity.RouteRequest, slippage uint64) (ret []*entity.RouteResponse, msg string, code int) {
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
	tokenAmountOut, ok := new(big.Float).SetString(tonRoute.SrcChain.TokenAmountOut)
	if !ok {
		params := map[string]interface{}{
			"request": utils.JSON(tonRequest),
			"amount":  tonRoute.SrcChain.TokenAmountOut,
		}
		log.Logger().WithFields(params).Error("failed to parse token amount out")
		return ret, "", resp.CodeTONRouteServerError
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

	gasFee := entity.Fee{
		Amount: tonRoute.GasFee.Amount,
		Symbol: tonRoute.GasFee.Symbol,
	}
	bridgeFee := entity.Fee{
		Amount: calcBridgeFee(tokenAmountOut, feeRate).Text('f', 6),
		Symbol: constants.BridgeFeeSymbol,
	}

	request := &butter.RouteRequest{
		TokenInAddress:  constants.USDTOfChainPoll,
		TokenOutAddress: req.TokenOutAddress,
		Type:            req.Type,
		Slippage:        slippage,
		FromChainID:     constants.ChainIDOfChainPool,
		ToChainID:       req.ToChainID,
		Amount:          tonRoute.SrcChain.TokenAmountOut,
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
	ret = make([]*entity.RouteResponse, 0, len(butterRoutes))
	for _, r := range butterRoutes {
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
			GasFee:    gasFee,
			BridgeFee: bridgeFee,
		}
		ret = append(ret, n)
	}
	return ret, "", resp.CodeSuccess
}

func GetEVMToTONRoute(req *entity.RouteRequest, slippage uint64) (ret []*entity.RouteResponse, msg string, code int) {
	var (
		tonTokenIn  entity.Token
		tonTokenOut entity.Token
	)

	request := &butter.RouteRequest{
		TokenInAddress:  req.TokenInAddress,
		TokenOutAddress: constants.USDTOfChainPoll,
		Type:            req.Type,
		Slippage:        slippage / 3 * 2,
		FromChainID:     req.FromChainID,
		ToChainID:       constants.ChainIDOfChainPool,
		Amount:          req.Amount,
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

	tonRequest := &tonrouter.RouteRequest{
		TokenInAddress:  constants.USDTOfTON,
		TokenOutAddress: req.TokenOutAddress,
		Slippage:        slippage,
	}
	tonRoutes, err := getTONRoutes(tonRequest, butterRoutes) // todo skip error ?
	if err != nil {
		return ret, "", resp.CodeTONRouteServerError
	}
	if len(tonRoutes) != len(butterRoutes) {
		return ret, "", resp.CodeTONRouteServerError
	}

	ret = make([]*entity.RouteResponse, 0, len(butterRoutes))
	for _, r := range butterRoutes {
		tonRoute, ok := tonRoutes[r.Hash]
		if !ok {
			continue
		}

		amountOut, ok := new(big.Float).SetString(r.DstChain.TotalAmountOut)
		if !ok {
			params := map[string]interface{}{
				"request": utils.JSON(tonRequest),
				"amount":  r.DstChain.TotalAmountOut,
			}
			log.Logger().WithFields(params).Error("failed to parse token amount out")
			return ret, "", resp.CodeTONRouteServerError
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
				Amount: calcBridgeFee(amountOut, feeRate).Text('f', 6),
				Symbol: constants.BridgeFeeSymbol,
			},
		}
		ret = append(ret, n)
	}
	return ret, "", resp.CodeSuccess
}

func GetSwapFromTON(sender string, dstChain, receiver, hash string) (ret *entity.SwapResponse, msg string, code int) {
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
	balance, err := getChainPoolBalance(dstChain)
	if err != nil {
		log.Logger().WithField("error", err).Error("failed to get ton router balance")
		return ret, "", resp.CodeInternalServerError
	}
	if amountOut.Cmp(balance) == 1 {
		params := map[string]interface{}{
			"amount":  amountOut.Text('f', -1),
			"balance": balance.Text('f', -1),
		}
		log.Logger().WithFields(params).Info("amount is greater than balance")
		return nil, "", resp.CodeInsufficientLiquidity // todo chain pool
	}

	request := &tonrouter.BridgeSwapRequest{
		Sender:   sender,
		Receiver: receiver,
		Hash:     hash,
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

func GetSwapFromEVM(srcChain *big.Int, srcToken, sender, amount string, dstChain *big.Int, dstToken, receiver, hash string, slippage uint64) (ret *entity.SwapResponse, msg string, code int) {
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

	chainPoolToken := constants.USDTOfChainPoll
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
	encodedCallback, err := EncodeSwapCallbackParams(common.HexToAddress(viper.GetString("feRouterContract")), common.HexToAddress(sender), packed) // todo sender
	if err != nil {
		params := map[string]interface{}{
			"feRouter": viper.GetString("feRouterContract"),
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

func getChainPoolBalance(dstChain string) (balance *big.Float, err error) {
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

func getUSDTDecimalOfChainPool(chain string) (decimal *big.Float) {
	switch chain {
	case constants.ChainIDOfEthereum:
		decimal = big.NewFloat(constants.USDTDecimalOfEthereum) // todo
	default:
		decimal = big.NewFloat(constants.USDTDecimalOfChainPool) // todo
	}
	return decimal
}

func calcBridgeFee(amount, feeRate *big.Float) (feeAmount *big.Float) {
	feeAmount = new(big.Float).Mul(amount, feeRate)
	feeAmount = new(big.Float).Quo(feeAmount, big.NewFloat(10000))
	log.Logger().WithField("amount", amount).WithField("feeRate", feeRate).WithField("feeAmount", feeAmount).Info("calc bridge fee")
	return feeAmount
}
