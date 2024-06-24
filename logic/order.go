package logic

import (
	"errors"
	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resp"
	"github.com/mapprotocol/fe-backend/utils"
	"gorm.io/gorm"
)

func choiceSubWalletID() int64 {
	return 0
}

func getChainNameByChainID(chainID uint64) string {
	return "ETH"
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
		//var depositAddress, err = ceffu.GetClient().GetDepositAddress(context.Background(), chainName, dstToken, subWalletID)
		//if err != nil {
		//	params := map[string]interface{}{
		//		"chainName":   chainName,
		//		"token":       dstToken,
		//		"subWalletID": subWalletID,
		//		"error":       err,
		//	}
		//	log.Logger().WithFields(params).Error("failed to get deposit address")
		//	return nil, resp.CodeInternalServerError
		//}

		//account = &dao.Account{
		//	SubWalletID: subWalletID,
		//	ChainID:     dstChain,
		//	ChainName:   chainName,
		//	//Address:     depositAddress,
		//}
		//if err := account.Create(); err != nil {
		//	log.Logger().WithField("account", utils.JSON(account)).WithField("error", err).Error("failed to create account")
		//	return nil, resp.CodeInternalServerError
		//}
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
		Stage:  dao.SwapStageDeposit,
		Status: dao.DepositStatusPending,
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

func OrderList(sender string, page, size int) (ret []*entity.OrderListResponse, count int64, code int) {
	list, count, err := dao.NewDepositSwapWithSender(sender).Find(nil, dao.Paginate(page, size))
	if err != nil {
		fields := map[string]interface{}{
			"page":  page,
			"size":  size,
			"error": err,
		}
		log.Logger().WithFields(fields).Error("failed to get deposit swap list")
		return nil, 0, resp.CodeInternalServerError
	}

	length := len(list)
	if length == 0 {
		return []*entity.OrderListResponse{}, count, resp.CodeSuccess
	}

	ret = make([]*entity.OrderListResponse, 0, length)
	for _, s := range list {
		ret = append(ret, &entity.OrderListResponse{
			OrderID:        s.ID,
			SrcChain:       s.SrcChain,
			SrcToken:       s.SrcToken,
			Sender:         s.Sender,
			Amount:         s.Amount,
			DstChain:       s.DstChain,
			DstToken:       s.DstToken,
			Receiver:       s.Receiver,
			DepositAddress: s.DepositAddress,
			TxHash:         s.TxHash,
			Action:         s.Action,
			Status:         s.Status,
			CreatedAt:      s.CreatedAt.Unix(),
		})
	}
	return ret, count, resp.CodeSuccess
}

func OrderDetail(orderID uint64) (ret *entity.OrderDetailResponse, code int) {
	order, err := dao.NewDepositSwapWithID(orderID).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Logger().WithField("orderID", orderID).WithField("error", err).Error("failed to get deposit swap")
		return nil, resp.CodeInternalServerError
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resp.CodeOrderNotFound
	}

	return &entity.OrderDetailResponse{
		OrderID:        order.ID,
		SrcChain:       order.SrcChain,
		SrcToken:       order.SrcToken,
		Sender:         order.Sender,
		Amount:         order.Amount,
		DstChain:       order.DstChain,
		DstToken:       order.DstToken,
		Receiver:       order.Receiver,
		DepositAddress: order.DepositAddress,
		TxHash:         order.TxHash,
		Action:         order.Action,
		Status:         order.Status,
		CreatedAt:      order.CreatedAt.Unix(),
	}, resp.CodeSuccess
}
