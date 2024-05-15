package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameSupportedToken = "supported_token"

type SupportedToken struct {
	ID        uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	ChainID   uint64    `gorm:"column:chain_id;type:bigint(20)" json:"chain_id"`
	Name      string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Symbol    string    `gorm:"column:symbol;type:varchar(255)" json:"symbol"`
	Decimal   uint32    `gorm:"column:decimal;type:int(11)" json:"decimal"`
	Icon      string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func (*SupportedToken) TableName() string {
	return TableNameSupportedToken
}

func NewSupportedToken(chainID uint64, symbol string) *SupportedToken {
	return &SupportedToken{
		ChainID: chainID,
		Symbol:  symbol,
	}
}

func (st *SupportedToken) Find(ext *QueryExtra, pager Pager) (list []*SupportedToken, count int64, err error) {
	tx := db.GetDB().Where(st)
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
		if err := tx.Model(st).Count(&count).Error; err != nil {
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
