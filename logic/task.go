package logic

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mapprotocol/ceffu-fe-backend/dao"
	"github.com/mapprotocol/ceffu-fe-backend/resource/binance"
	"github.com/mapprotocol/ceffu-fe-backend/resource/ceffu"
	"github.com/mapprotocol/ceffu-fe-backend/resource/log"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"github.com/mapprotocol/ceffu-fe-backend/utils/tx"
	"time"
)

func TaskQueryTransferWithExchangeStatus() {
	order := dao.NewDepositSwap()
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByStatus(id, dao.MirrorStatusPending) //  todo add stage to query params
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

				// todo convert order status
				if transferDetail.Status == dao.MirrorStatusConfirmed || transferDetail.Status == dao.MirrorStatusFailed {
					update := &dao.DepositSwap{
						Status:    transferDetail.Status,
						OutAmount: transferDetail.Amount,
					}
					if err := dao.NewDepositSwapWithID(o.ID).Updates(update); err != nil {
						log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func TaskQuerySellOrderStatus() {
	order := dao.NewDepositSwap()
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByStatus(id, dao.SellStatusSent) //  todo add stage to query params
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get sell sent status order")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found sell sent status order", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				symbol := "BTCUSDT" // todo
				gotOrder, err := binance.GetClient().NewGetOrderService().Symbol(symbol).OrderId(o.ExchangeOrderID).Do(context.Background())
				if err != nil {
					params := map[string]interface{}{
						"orderID":         o.ID,
						"symbol":          symbol,
						"exchangeOrderID": o.ExchangeOrderID,
						"error":           err.Error(),
					}
					log.Logger().WithFields(params).Error("failed to get order")
					time.Sleep(5 * time.Second)
					continue
				}

				switch gotOrder.Status {
				case BinanceOrderStatusNew, BinanceOrderStatusPartiallyFilled:
					continue
				case BinanceOrderStatusFilled:
					trades, err := binance.GetClient().NewGetMyTradesService().Symbol(symbol).OrderId(o.ExchangeOrderID).Do(context.Background())
					if err != nil {
						params := map[string]interface{}{
							"orderID":         o.ID,
							"symbol":          symbol,
							"exchangeOrderID": o.ExchangeOrderID,
							"error":           err.Error(),
						}
						log.Logger().WithFields(params).Error("failed to get my trades")
						time.Sleep(5 * time.Second)
						continue
					}
					if trades == nil || len(trades) != 0 {
						// todo add alarm
						continue
					}
					// todo convert order status
					update := &dao.DepositSwap{
						Status:    dao.SellStatusConfirmed, // todo
						OutAmount: trades[0].QuoteQuantity, // todo
					}
					if err := dao.NewDepositSwapWithID(o.ID).Updates(update); err != nil {
						log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
						time.Sleep(5 * time.Second)
						continue
					}
				default:
					// todo add alarm
				}

			}
			time.Sleep(10 * time.Second)
		}
	}
}

func TaskButterSwapStatus() {
	order := dao.NewDepositSwap()
	for {
		for id := uint64(1); ; {
			orders, err := order.GetOldest10ByStatus(id, dao.ButterSwapStatusSent) //  todo add stage to query params
			if err != nil {
				log.Logger().WithField("error", err.Error()).Error("failed to get butter swap sent status order")
				time.Sleep(5 * time.Second)
				continue
			}

			length := len(orders)
			if length == 0 {
				log.Logger().Info("not found butter swap sent status order", "time", time.Now())
				time.Sleep(10 * time.Second)
				break
			}

			for i, o := range orders {
				if i == length-1 {
					id = o.ID + 1
				}

				chainInfo, err := dao.NewSupportedChainWithChainID(o.DstChain).First()
				if err != nil {
					log.Logger().WithField("chainID", o.DstChain).WithField("error", err.Error()).Error("failed to get chain info")
					time.Sleep(5 * time.Second)
					continue
				}

				transactor, err := tx.NewTransactor(chainInfo.ChainRPC)
				if err != nil {
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("error", err.Error()).Error("failed to create transactor")
					time.Sleep(5 * time.Second)
					continue
				}

				pending, err := transactor.TransactionIsPending(common.HexToHash(o.SwapTxHash))
				if err != nil {
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("error", err.Error()).Error("failed to create transactor")
					time.Sleep(5 * time.Second)
					continue
				}
				if pending {
					time.Sleep(5 * time.Second)
					continue
				}

				status, err := transactor.TransactionStatus(common.HexToHash(o.SwapTxHash))
				if err != nil {
					log.Logger().WithField("endpoint", chainInfo.ChainRPC).WithField("txHash", o.SwapTxHash).WithField("error", err.Error()).Error("get transaction status")
					time.Sleep(5 * time.Second)
					continue
				}

				update := &dao.DepositSwap{
					Status: dao.ButterSwapStatusConfirmed,
				}
				if status == types.ReceiptStatusFailed {
					update.Status = dao.ButterSwapStatusFailed
				}
				if err := dao.NewDepositSwapWithID(o.ID).Updates(update); err != nil {
					log.Logger().WithField("update", utils.JSON(update)).WithField("error", err.Error()).Error("failed to update order status")
					time.Sleep(5 * time.Second)
					continue
				}

			}
			time.Sleep(10 * time.Second)
		}
	}
}
