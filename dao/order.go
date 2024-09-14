package dao

import (
	"time"

	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/utils"
)

const TableNameOrder = "order"

const OldestLimit = 10

const (
	OrderActionToEVM = iota + 1
	OrderActionFromEVM
)

const (
	OrderStag1 = iota + 1
	OrderStag2
)

//const (
//	OrderStatusPending = iota + 1
//	OrderStatusConfirmed
//	OrderStatusFailed
//)

const Stage1StatusConfirmed = 1

const (
	OrderStatusTxPrepareSend = iota + 1
	OrderStatusTxSent
	OrderStatusTxFailed
	OrderStatusTxConfirmed
	OrderStatusCompleted
)

type Order struct {
	ID                  uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	OrderIDFromContract uint64    `gorm:"column:order_id_from_contract;type:bigint(20)" json:"order_id_from_contract"`
	SrcChain            string    `gorm:"column:src_chain;type:varchar(255)" json:"src_chain"`
	SrcToken            string    `gorm:"column:src_token;type:varchar(255)" json:"src_token"`
	Sender              string    `gorm:"column:sender;type:varchar(255)" json:"sender"`
	InAmount            string    `gorm:"column:in_amount;type:varchar(255)" json:"in_amount"`
	InTxHash            string    `gorm:"column:in_tx_hash;type:varchar(255)" json:"in_tx_hash"`
	Relayer             string    `gorm:"column:relayer;type:varchar(255)" json:"relayer"`
	RelayerKey          string    `gorm:"column:relayer_key;type:varchar(255)" json:"relayer_key"`
	RelayToken          string    `gorm:"column:relay_token;type:varchar(255)" json:"relay_token"`
	RelayAmount         string    `gorm:"column:relay_amount;type:varchar(255)" json:"relay_amount"`
	DstChain            string    `gorm:"column:dst_chain;type:varchar(255)" json:"dst_chain"`
	DstToken            string    `gorm:"column:dst_token;type:varchar(255)" json:"dst_token"`
	Receiver            string    `gorm:"column:receiver;type:varchar(255)" json:"receiver"`
	OutAmount           string    `gorm:"column:out_amount;type:varchar(255)" json:"out_amount"`
	OutTxHash           string    `gorm:"column:out_tx_hash;type:varchar(255)" json:"out_tx_hash"`
	Action              uint8     `gorm:"column:action;type:tinyint(4)" json:"action"`
	Stage               uint8     `gorm:"column:stage;type:tinyint(4)" json:"stage"`
	Status              uint8     `gorm:"column:status;type:int(11)" json:"status"`
	Slippage            uint64    `gorm:"column:slippage;type:bigint(20)" json:"slippage"`
	MinAmountOut        string    `gorm:"column:min_amount_out;type:bigint(20)" json:"min_amount_out"`
	CreatedAt           time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewOrder() *Order {
	return &Order{}
}

func NewOrderWithID(id uint64) *Order {
	return &Order{ID: id}
}

func NewOrderWithOrderIDFromContract(OrderIDFromContract uint64) *Order {
	return &Order{OrderIDFromContract: OrderIDFromContract}
}

func (o *Order) TableName() string {
	return TableNameOrder
}

func (o *Order) Create() error {
	return db.GetDB().Create(o).Error
}

func (o *Order) Updates(update *Order) error {
	return db.GetDB().Where(o).Updates(update).Error
}

func (o *Order) First() (get *Order, err error) {
	err = db.GetDB().Where(o).First(&get).Error
	return get, err
}

func (o *Order) Last() (get *Order, err error) {
	err = db.GetDB().Where(o).Last(&get).Error
	return get, err
}

func (o *Order) Find(ext *QueryExtra, pager Pager) (list []*Order, count int64, err error) {
	tx := db.GetDB().Where(o)
	if ext != nil {
		if ext.Conditions != nil {
			for k, v := range ext.Conditions {
				tx = tx.Where(k, v)
			}
		}
		if ext.OnlyKeyConditions != nil {
			for k := range ext.OnlyKeyConditions {
				tx = tx.Where(k)
			}
		}
		if !utils.IsEmpty(ext.OrderStr) {
			tx = tx.Order(ext.OrderStr)
		}
	}

	if pager != nil {
		if err := tx.Model(o).Count(&count).Error; err != nil {
			return list, count, err
		}
		if count == 0 {
			return list, 0, nil
		}
		tx = tx.Scopes(pager)
	}
	err = tx.Find(&list).Error
	return list, count, err
}

func (o *Order) GetOldest10ByStatus(id uint64, action, stage, status uint8) (list []*Order, err error) {
	err = db.GetDB().Where(o).Where("id >= ?", id).Where("action = ?", action).
		Where("stage = ?", stage).Where("status = ?", status).Limit(OldestLimit).Find(&list).Error
	return list, err
}

func (o *Order) GetOldest10ByID(id uint64) (list []*Order, err error) {
	err = db.GetDB().Where(o).Where("id >= ?", id).Limit(OldestLimit).Find(&list).Error
	return list, err
}
