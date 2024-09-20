package dao

import (
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/utils"
	"time"
)

const TableNameSolOrder = "sol_order"

type SolOrder struct {
	ID             uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	SrcHash        string    `gorm:"column:src_hash;type:varchar(255)" json:"src_hash"`
	SrcChain       string    `gorm:"column:src_chain;type:varchar(255)" json:"src_chain"`
	SrcToken       string    `gorm:"column:src_token;type:varchar(255)" json:"src_token"`
	Sender         string    `gorm:"column:sender;type:varchar(255)" json:"sender"`
	InAmount       string    `gorm:"column:in_amount;type:varchar(255)" json:"in_amount"`
	Relayer        string    `gorm:"column:relayer;type:varchar(255)" json:"relayer"`
	RelayerKey     string    `gorm:"column:relayer_key;type:varchar(255)" json:"relayer_key"`
	RelayToken     string    `gorm:"column:relay_token;type:varchar(255)" json:"relay_token"`
	RelayAmount    string    `gorm:"column:relay_amount;type:varchar(255)" json:"relay_amount"`
	DstChain       string    `gorm:"column:dst_chain;type:varchar(255)" json:"dst_chain"`
	DstToken       string    `gorm:"column:dst_token;type:varchar(255)" json:"dst_token"`
	Receiver       string    `gorm:"column:receiver;type:varchar(255)" json:"receiver"`
	OutAmount      string    `gorm:"column:out_amount;type:varchar(255)" json:"out_amount"`
	OutTxHash      string    `gorm:"column:out_tx_hash;type:varchar(255)" json:"out_tx_hash"`
	Action         uint8     `gorm:"column:action;type:tinyint(4)" json:"action"`
	Stage          uint8     `gorm:"column:stage;type:tinyint(4)" json:"stage"`
	Status         uint8     `gorm:"column:status;type:int(11)" json:"status"`
	Slippage       uint64    `gorm:"column:slippage;type:bigint(20)" json:"slippage"`
	BridgeId       uint64    `gorm:"column:bridge_id;type:bigint(20)" json:"bridge_id"`
	OrderId        string    `gorm:"column:order_id;type:varchar(255)" json:"order_id"`
	ChainPoolToken string    `gorm:"column:chain_pool_token;type:varchar(255)" json:"chain_pool_token"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewSolOrder() *SolOrder {
	return &SolOrder{}
}

func NewSolOrderWithID(id uint64) *SolOrder {
	return &SolOrder{ID: id}
}

func (o *SolOrder) TableName() string {
	return TableNameSolOrder
}

func (o *SolOrder) Create() error {
	return db.GetDB().Create(o).Error
}

func (o *SolOrder) Updates(update *SolOrder) error {
	return db.GetDB().Where(o).Updates(update).Error
}

func (o *SolOrder) First() (get *SolOrder, err error) {
	err = db.GetDB().Where(o).First(&get).Error
	return get, err
}

func (o *SolOrder) Last() (get *SolOrder, err error) {
	err = db.GetDB().Where(o).Last(&get).Error
	return get, err
}

func (o *SolOrder) Find(ext *QueryExtra, pager Pager) (list []*SolOrder, count int64, err error) {
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

func (o *SolOrder) GetOldest10ByStatus(id uint64, action, stage, status uint8) (list []*SolOrder, err error) {
	err = db.GetDB().Where(o).Where("id >= ?", id).Where("action = ?", action).
		Where("stage = ?", stage).Where("status = ?", status).Limit(OldestLimit).Find(&list).Error
	return list, err
}

func (o *SolOrder) GetOldest10ByID(id uint64) (list []*SolOrder, err error) {
	err = db.GetDB().Where(o).Where("id >= ?", id).Limit(OldestLimit).Find(&list).Error
	return list, err
}
