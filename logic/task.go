package logic

import (
	"context"
	"github.com/mapprotocol/ceffu-fe-backend/dao"
	"github.com/mapprotocol/ceffu-fe-backend/resource/ceffu"
	"github.com/mapprotocol/ceffu-fe-backend/resource/log"
	"time"
)

func TaskQueryTransferWithExchangeStatus() {
	order := dao.NewDepositSwap()
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByStatus(id, dao.MirrorAndSellStatusPending)
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get mirror pending status order")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found mirror pending status order", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				var walletID int64 = 0 // todo validate wallet id
				transferDetail, err := ceffu.GetClient().TransferDetailWithExchange(context.Background(), o.OrderViewID, walletID)
				if err != nil {
					params := map[string]interface{}{
						"orderID":     o.ID,
						"orderViewID": o.OrderViewID,
						"walletID":    walletID,
						"error":       err.Error(),
					}
					log.Logger().WithFields(params).Error("failed to get transfer detail")
					time.Sleep(5 * time.Second)
					continue
				}

				if transferDetail.Status == dao.MirrorAndSellStatusConfirmed || transferDetail.Status == dao.MirrorAndSellStatusFailed {
					update := &dao.DepositSwap{
						Status:    transferDetail.Status,
						OutAmount: transferDetail.Amount,
					}
					if err := dao.NewDepositSwapWithID(o.ID).Updates(update); err != nil {
						log.Logger().WithField("error", err.Error()).Error("failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}
}
