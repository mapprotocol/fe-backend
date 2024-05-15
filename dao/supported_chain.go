package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameSupportedChain = "supported_chain"

type SupportedChain struct {
	ID        uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	ChainID   uint64    `gorm:"column:chain_id;type:bigint(20)" json:"chain_id"`
	ChainName string    `gorm:"column:chain_name;type:varchar(255)" json:"chain_name"`
	ChainIcon string    `gorm:"column:chain_icon;type:varchar(255)" json:"chain_icon"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func (*SupportedChain) TableName() string {
	return TableNameSupportedChain
}

func NewSupportedChain() *SupportedChain {
	return &SupportedChain{}
}

func (sc *SupportedChain) Find(ext *QueryExtra, pager Pager) (list []*SupportedChain, count int64, err error) {
	tx := db.GetDB().Where(sc)
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
		if err := tx.Model(sc).Count(&count).Error; err != nil {
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
