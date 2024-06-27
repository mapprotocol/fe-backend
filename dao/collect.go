package dao

import (
	"time"

	"github.com/mapprotocol/fe-backend/resource/db"
)

const TableNameCollect = "collect"

const (
	CollectStatusPending = iota + 1
	CollectStatusConfirmed
)

type Collect struct {
	ID        uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	OrderID   uint64    `gorm:"column:order_id;type:bigint(20)" json:"order_id"`
	TxHash    string    `gorm:"column:tx_hash;type:varchar(255)" json:"tx_hash"`
	Status    int32     `gorm:"column:status;type:int(11)" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewCollect() *Collect {
	return &Collect{}
}

func (c *Collect) TableName() string {
	return TableNameCollect
}

func (c *Collect) Create() error {
	return db.GetDB().Create(c).Error
}

func (c *Collect) Updates(update *Collect) error {
	return db.GetDB().Where(c).Updates(update).Error
}

func (c *Collect) UpdatesByIDs(ids []uint64, update *Collect) error {
	return db.GetDB().Where("id in ?", ids).Updates(update).Error
}

func (c *Collect) First() (get *Collect, err error) {
	err = db.GetDB().Where(c).First(&get).Error
	return get, err
}
