package logic

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"gorm.io/gorm"
	blog "log"

	"github.com/mapprotocol/fe-backend/dao"
	"github.com/mapprotocol/fe-backend/entity"
	"github.com/mapprotocol/fe-backend/resource/log"
	"github.com/mapprotocol/fe-backend/resp"
)

var NetParams = &chaincfg.Params{}

const (
	NetworkMainnet = "mainnet"
	NetworkTestnet = "testnet"
)

func InitNetworkParams(network string) {
	switch network {
	case NetworkMainnet, "":
		NetParams = &chaincfg.MainNetParams
		blog.Print("initialized network: ", NetworkMainnet)
	case NetworkTestnet:
		NetParams = &chaincfg.TestNet3Params
		blog.Print("initialized network: ", NetworkTestnet)
	default:
		panic("unknown network")
	}
}

//func CreateOrder(srcChain, srcToken, sender, amount string, dstChain, dstToken, receiver string, action uint8, slippage uint64) (ret *entity.CreateOrderResponse, code int) {
//	var (
//		addressStr    string
//		privateKeyStr string
//	)
//
//	if srcChain == constants.BTCChainID {
//		privateKey, err := generateKey()
//		if err != nil {
//			log.Logger().WithField("error", err).Error("failed to generate key")
//			return nil, resp.CodeInternalServerError
//		}
//		address, err := makeTaprootAddress(privateKey, NetParams)
//		if err != nil {
//			log.Logger().WithField("error", err).Error("failed to make address")
//			return nil, resp.CodeInternalServerError
//		}
//		addressStr = address.String()
//		privateKeyStr = string(privateKey.Serialize())
//	}
//
//	order := &dao.Order{
//		SrcChain:   srcChain,
//		SrcToken:   srcToken,
//		Sender:     sender,
//		InAmount:   amount,
//		Relayer:    addressStr,
//		RelayerKey: privateKeyStr,
//		DstChain:   dstChain,
//		DstToken:   dstToken,
//		Receiver:   receiver,
//		Action:     action,
//		Stage:      dao.OrderStag1,
//		Status:     dao.OrderStatusPending,
//		Slippage:   slippage,
//	}
//	orderID, err := order.Create()
//	if err != nil {
//		log.Logger().WithField("order", utils.JSON(order)).WithField("error", err).Error("failed to create order")
//		return nil, resp.CodeInternalServerError
//	}
//
//	return &entity.CreateOrderResponse{
//		OrderID: orderID,
//		Relayer: addressStr,
//	}, resp.CodeSuccess
//}

func UpdateOrder(orderID uint64, txHash string) int {
	order, err := dao.NewOrderWithID(orderID).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp.CodeOrderNotFound
		}
		log.Logger().WithField("order_id", orderID).WithField("error", err).Error("failed to get order")
		return resp.CodeInternalServerError
	}
	if order.Action == dao.OrderActionToEVM {
		// todo check tx hash
	} else if order.Action == dao.OrderActionFromEVM {
		// todo check tx hash
		//common.HexToHash()
	}

	update := dao.Order{
		ID:       orderID,
		InTxHash: txHash,
	}
	if err := dao.NewOrder().Updates(&update); err != nil {
		params := map[string]interface{}{
			"order_id": orderID,
			"txHash":   txHash,
			"error":    err,
		}
		log.Logger().WithFields(params).Error("failed to update order")
		return resp.CodeInternalServerError
	}
	return resp.CodeSuccess
}

func OrderList(sender string, page, size int) (ret []*entity.OrderListResponse, count int64, code int) {
	list, count, err := dao.NewOrderWithSender(sender).Find(nil, dao.Paginate(page, size))
	if err != nil {
		fields := map[string]interface{}{
			"page":  page,
			"size":  size,
			"error": err,
		}
		log.Logger().WithFields(fields).Error("failed to get order list")
		return nil, 0, resp.CodeInternalServerError
	}

	length := len(list)
	if length == 0 {
		return []*entity.OrderListResponse{}, count, resp.CodeSuccess
	}

	ret = make([]*entity.OrderListResponse, 0, length)
	for _, s := range list {
		ret = append(ret, &entity.OrderListResponse{
			OrderID:   s.ID,
			SrcChain:  s.SrcChain,
			SrcToken:  s.SrcToken,
			Sender:    s.Sender,
			InAmount:  s.InAmount,
			DstChain:  s.DstChain,
			DstToken:  s.DstToken,
			Receiver:  s.Receiver,
			OutAmount: s.OutAmount,
			Action:    s.Action,
			Status:    s.Status,
			CreatedAt: s.CreatedAt.Unix(),
		})
	}
	return ret, count, resp.CodeSuccess
}

func OrderDetail(orderID uint64) (ret *entity.OrderDetailResponse, code int) {
	order, err := dao.NewOrderWithID(orderID).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Logger().WithField("orderID", orderID).WithField("error", err).Error("failed to get order")
		return nil, resp.CodeInternalServerError
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, resp.CodeOrderNotFound
	}

	return &entity.OrderDetailResponse{
		OrderID:   order.ID,
		SrcChain:  order.SrcChain,
		SrcToken:  order.SrcToken,
		Sender:    order.Sender,
		InAmount:  order.InAmount,
		DstChain:  order.DstChain,
		DstToken:  order.DstToken,
		Receiver:  order.Receiver,
		OutAmount: order.OutAmount,
		Action:    order.Action,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.Unix(),
	}, resp.CodeSuccess
}

func generateKey() (*btcec.PrivateKey, error) {
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func makeTaprootAddress(privKey *btcec.PrivateKey, netParams *chaincfg.Params) (btcutil.Address, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(privKey.PubKey())

	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey),
		netParams,
	)
	if err != nil {
		return nil, err
	}
	return address, nil
}
