package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameSwapProcess = "swap_process"

type SwapProcess struct {
	ID        uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewSwapProcess() *SwapProcess {
	return &SwapProcess{}
}

func (sp *SwapProcess) TableName() string {
	return TableNameSwapProcess
}

func (sp *SwapProcess) Create() error {
	return db.GetDB().Create(sp).Error
}

func (sp *SwapProcess) Updates(update *SwapProcess) error {
	return db.GetDB().Where(sp).Updates(update).Error
}

func (sp *SwapProcess) First() (get *SwapProcess, err error) {
	err = db.GetDB().Where(sp).First(&get).Error
	return get, err
}

func (sp *SwapProcess) Last() (get *SwapProcess, err error) {
	err = db.GetDB().Where(sp).Last(&get).Error
	return get, err
}

func (sp *SwapProcess) Find(ext *QueryExtra, pager Pager) (list []*SwapProcess, count int64, err error) {
	tx := db.GetDB().Where(sp)
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
		if err := tx.Model(sp).Count(&count).Error; err != nil {
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
