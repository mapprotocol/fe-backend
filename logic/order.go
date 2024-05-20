package logic

import (
	"errors"
	"github.com/mapprotocol/ceffu-fe-backend/dao"
	"github.com/mapprotocol/ceffu-fe-backend/entity"
	"github.com/mapprotocol/ceffu-fe-backend/resource/log"
	"github.com/mapprotocol/ceffu-fe-backend/resp"
	"gorm.io/gorm"
)

func choiceSubWalletID() int64 {
	return 0
}

func getChainNameByChainID(chainID uint64) string {
	return "ETH"
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
	order, err := dao.NewDepositSwap(orderID).First()
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
