package logic

import (
	"errors"
	"github.com/mapprotocol/ceffu-fe-backend/dao"
	"github.com/mapprotocol/ceffu-fe-backend/entity"
	"github.com/mapprotocol/ceffu-fe-backend/resource/log"
	"github.com/mapprotocol/ceffu-fe-backend/resp"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"github.com/mapprotocol/ceffu-fe-backend/utils/ceffu"
	"gorm.io/gorm"
)

func SupportedTokens(chainID uint64, symbol string, page, size int) (ret []*entity.SupportedTokensResponse, count int64, code int) {
	list, count, err := dao.NewSupportedToken(chainID, symbol).Find(nil, dao.Paginate(page, size))
	if err != nil {
		fields := map[string]interface{}{
			"page":  page,
			"size":  size,
			"error": err,
		}
		log.Logger().WithFields(fields).Error("failed to get supported chain list")
		return nil, 0, resp.CodeInternalServerError
	}

	length := len(list)
	if length == 0 {
		return []*entity.SupportedTokensResponse{}, count, resp.CodeSuccess
	}

	ret = make([]*entity.SupportedTokensResponse, 0, length)
	for _, c := range list {
		ret = append(ret, &entity.SupportedTokensResponse{
			ChainID:  c.ChainID,
			Symbol:   c.Symbol,
			Name:     c.Name,
			Decimals: c.Decimal,
		})
	}
	return ret, count, resp.CodeSuccess
}

func CreateOrder(srcChain uint64, srcToken, sender, amount string, dstChain uint64, dstToken, receiver string) (ret *entity.CreateOrderResponse, code int) {
	account, err := dao.NewAccount(dstChain).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Logger().WithField("dstChain", dstChain).WithField("error", err).Error("failed to get account")
		return nil, resp.CodeInternalServerError
	}

	// create account if not exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		subWalletID := choiceSubWalletID()           // TODO: choice sub wallet id
		chainName := getChainNameByChainID(dstChain) // TODO: get chain name by chain id
		// create deposit address
		depositAddress, err := ceffu.GetDepositAddress(chainName, dstToken, subWalletID)
		if err != nil {
			params := map[string]interface{}{
				"chainName":   chainName,
				"token":       dstToken,
				"subWalletID": subWalletID,
				"error":       err,
			}
			log.Logger().WithFields(params).Error("failed to get deposit address")
			return nil, resp.CodeInternalServerError
		}

		account = &dao.Account{
			SubWalletID: subWalletID,
			ChainID:     dstChain,
			ChainName:   chainName,
			Address:     depositAddress,
		}
		if err := account.Create(); err != nil {
			log.Logger().WithField("account", utils.JSON(account)).WithField("error", err).Error("failed to create account")
			return nil, resp.CodeInternalServerError
		}
	}

	// create order
	order := &dao.DepositSwap{
		SrcChain:       srcChain,
		SrcToken:       srcToken,
		Sender:         sender,
		Amount:         amount,
		DstChain:       dstChain,
		DstToken:       dstToken,
		Receiver:       receiver,
		DepositAddress: account.Address,
		//Mask:     1, // TODO: set mask
		//Action:   1, // TODO: set action
		Status: dao.DepositSwapStatusPending,
	}
	if err := order.Create(); err != nil {
		log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
		return nil, resp.CodeInternalServerError
	}

	return &entity.CreateOrderResponse{
		OrderID:        order.ID,
		DepositAddress: account.Address,
	}, resp.CodeSuccess
}
