package dao

import (
	"time"

	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/utils"
)

const TableNameOrder = "bitcoin_order"

const OldestLimit = 10

const (
	OrderActionToEVM = iota + 1
	OrderActionFromEVM
)

const (
	OrderStag1 = iota + 1
	OrderStag2
)

const (
	OrderStatusTxPrepareSend = iota + 1
	OrderStatusTxSent
	OrderStatusTxFailed
	OrderStatusTxConfirmed
	OrderStatusCompleted
)

const (
	WithdrawStateInit = iota + 1
	WithdrawStateSend
	WithdrawStateFinish
)

type BitcoinOrder struct {
	ID             uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	SrcChain       uint64    `gorm:"column:src_chain;type:bigint(20)" json:"src_chain"`
	SrcToken       string    `gorm:"column:src_token;type:varchar(255)" json:"src_token"`
	Sender         string    `gorm:"column:sender;type:varchar(255)" json:"sender"`
	InAmount       string    `gorm:"column:in_amount;type:varchar(255)" json:"in_amount"`
	InTxHash       string    `gorm:"column:in_tx_hash;type:varchar(255)" json:"in_tx_hash"`
	DstChain       uint64    `gorm:"column:dst_chain;type:bigint(20)" json:"dst_chain"`
	DstToken       string    `gorm:"column:dst_token;type:varchar(255)" json:"dst_token"`
	Receiver       string    `gorm:"column:receiver;type:varchar(255)" json:"receiver"`
	OutAmount      string    `gorm:"column:out_amount;type:varchar(255)" json:"out_amount"`
	OutTxHash      string    `gorm:"column:out_tx_hash;type:varchar(255)" json:"out_tx_hash"`
	Action         uint8     `gorm:"column:action;type:tinyint(4)" json:"action"`
	Stage          uint8     `gorm:"column:stage;type:tinyint(4)" json:"stage"`
	Status         uint8     `gorm:"column:status;type:tinyint(4)" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	Relayer        string    `gorm:"column:relayer;type:varchar(255)" json:"relayer"`
	RelayerKey     string    `gorm:"column:relayer_key;type:varchar(255)" json:"relayer_key"`
	RelayToken     string    `gorm:"column:relay_token;type:varchar(255)" json:"relay_token"`
	RelayAmount    string    `gorm:"column:relay_amount;type:varchar(255)" json:"relay_amount"`
	RelayAmountInt uint64    `gorm:"column:relay_amount_int;type:bigint(20)" json:"relay_amount_int"`
	InAmountSat    string    `gorm:"column:in_amount_sat;type:varchar(255)" json:"in_amount_sat"`
	CollectStatus  uint8     `gorm:"column:collect_status;type:tinyint(4)" json:"collect_status"`
}

func NewOrder() *BitcoinOrder {
	return &BitcoinOrder{}
}

func NewBitcoinOrderWithID(id uint64) *BitcoinOrder {
	return &BitcoinOrder{ID: id}
}

func NewBitcoinOrderWithSender(sender string) *BitcoinOrder {
	return &BitcoinOrder{Sender: sender}
}

func (o *BitcoinOrder) TableName() string {
	return TableNameOrder
}

func (o *BitcoinOrder) Create() error {
	return db.GetDB().Create(o).Error
}

func (o *BitcoinOrder) Updates(update *BitcoinOrder) error {
	return db.GetDB().Where(o).Updates(update).Error
}

func (o *BitcoinOrder) UpdatesByIDs(ids []uint64, update *BitcoinOrder) error {
	return db.GetDB().Where("id in ?", ids).Updates(update).Error
}

func (o *BitcoinOrder) First() (get *BitcoinOrder, err error) {
	err = db.GetDB().Where(o).First(&get).Error
	return get, err
}

func (o *BitcoinOrder) Last() (get *BitcoinOrder, err error) {
	err = db.GetDB().Where(o).Last(&get).Error
	return get, err
}

func (o *BitcoinOrder) Find(ext *QueryExtra, pager Pager) (list []*BitcoinOrder, count int64, err error) {
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
		//if err := tx.Model(o).Count(&count).Error; err != nil {
		//	return list, count, err
		//}
		//if count == 0 {
		//	return list, 0, nil
		//}
		tx = tx.Scopes(pager)
	}
	err = tx.Find(&list).Error
	return list, count, err
}

func (o *BitcoinOrder) GetOldest10ByID(id uint64, limit int) (list []*BitcoinOrder, err error) {
	err = db.GetDB().Where(o).Where("id > ?", id).Limit(limit).Find(&list).Error
	return list, err
}
