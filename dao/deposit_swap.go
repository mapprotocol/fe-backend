package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameDepositSwap = "deposit_swap"

const (
	DepositSwapActionDeposit = iota + 1
	DepositSwapActionSwap
)

const (
	DepositStatusPending = iota + 1
	DepositStatusConfirmed
)

const (
	MirrorStatusPending    = 10
	MirrorStatusProcessing = 20
	MirrorStatusSent       = 30
	MirrorStatusConfirmed  = 40
	MirrorStatusFailed     = 99
)

const (
	SellStatusSent = iota + 1
	SellStatusConfirmed
	SellStatusFailed
)

const (
	ButterSwapStatusSent = iota + 1
	ButterSwapStatusConfirmed
	ButterSwapStatusFailed
)

const (
	SwapStageDeposit = iota + 1
	SwapStageMirror
	SwapStageSell
)

const OldestLimit = 10

type DepositSwap struct {
	ID              uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	SrcChain        uint64    `gorm:"column:src_chain;type:bigint(20)" json:"src_chain"`
	SrcToken        string    `gorm:"column:src_token;type:varchar(255)" json:"src_token"`
	Amount          string    `gorm:"column:amount;type:varchar(255)" json:"amount"`
	Sender          string    `gorm:"column:sender;type:varchar(255)" json:"sender"`
	DstChain        uint64    `gorm:"column:dst_chain;type:bigint(20)" json:"dst_chain"`
	DstToken        string    `gorm:"column:dst_token;type:varchar(255)" json:"dst_token"`
	Receiver        string    `gorm:"column:receiver;type:varchar(255)" json:"receiver"`
	DepositAddress  string    `gorm:"column:deposit_address;type:varchar(255)" json:"deposit_address"`
	Mask            uint32    `gorm:"column:mask;type:int(11)" json:"mask"`
	TxHash          string    `gorm:"column:tx_hash;type:varchar(255)" json:"tx_hash"`
	Action          uint8     `gorm:"column:action;type:tinyint(4)" json:"action"`
	Stage           uint8     `gorm:"column:stage;type:tinyint(4)" json:"stage"`
	Status          int32     `gorm:"column:status;type:int(11)" json:"status"`
	OrderViewID     string    `gorm:"column:order_view_id;type:varchar(255)" json:"order_view_id"`
	ExchangeOrderID int64     `gorm:"column:exchange_order_id;type:bigint(20)" json:"exchange_order_id"`
	OutAmount       string    `gorm:"column:out_amount;type:varchar(255)" json:"out_amount"`
	SwapTxHash      string    `gorm:"column:swap_tx_hash;type:varchar(255)" json:"swap_tx_hash"`
	SwapChainID     string    `gorm:"column:swap_chain_id;type:varchar(255)" json:"swap_chain_id"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewDepositSwap() *DepositSwap {
	return &DepositSwap{}
}

func NewDepositSwapWithID(id uint64) *DepositSwap {
	return &DepositSwap{ID: id}
}

func NewDepositSwapWithSender(sender string) *DepositSwap {
	return &DepositSwap{Sender: sender}
}

func (ds *DepositSwap) TableName() string {
	return TableNameDepositSwap
}

func (ds *DepositSwap) Create() error {
	return db.GetDB().Create(ds).Error
}

func (ds *DepositSwap) Updates(update *DepositSwap) error {
	return db.GetDB().Where(ds).Updates(update).Error
}

func (ds *DepositSwap) First() (get *DepositSwap, err error) {
	err = db.GetDB().Where(ds).First(&get).Error
	return get, err
}

func (ds *DepositSwap) Last() (get *DepositSwap, err error) {
	err = db.GetDB().Where(ds).Last(&get).Error
	return get, err
}

func (ds *DepositSwap) Find(ext *QueryExtra, pager Pager) (list []*DepositSwap, count int64, err error) {
	tx := db.GetDB().Where(ds)
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
		if err := tx.Model(ds).Count(&count).Error; err != nil {
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

func (ds *DepositSwap) GetOldest10ByStatus(id uint64, status uint8) (list []*DepositSwap, err error) { // todo add stage to query params
	err = db.GetDB().Where(ds).Where("id >= ?", id).Where("status = ?", status).Limit(OldestLimit).Find(&list).Error
	return list, err
}
