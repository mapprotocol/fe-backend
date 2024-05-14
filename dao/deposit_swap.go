package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameDepositSwap = "deposit_swap"

type DepositSwap struct {
	ID        uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewDepositSwap() *DepositSwap {
	return &DepositSwap{}
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
