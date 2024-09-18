package dao

import (
	"time"

	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/utils"
)

const TableNameFilterLog = "sol_filter_log"

type FilterLog struct {
	ID          uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	ChainID     string    `gorm:"column:chain_id;type:varchar(255)" json:"chain_id"`
	Topic       string    `gorm:"column:topic;type:varchar(255)" json:"topic"`
	LatestLogID uint64    `gorm:"column:latest_log_id;type:bigint(20)" json:"latest_log_id"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewFilterLog(chainID, topic string) *FilterLog {
	return &FilterLog{
		ChainID: chainID,
		Topic:   topic,
	}
}

func (fl *FilterLog) TableName() string {
	return TableNameFilterLog
}

func (fl *FilterLog) Create() error {
	return db.GetDB().Create(fl).Error
}

func (fl *FilterLog) Updates(update *FilterLog) error {
	return db.GetDB().Where(fl).Updates(update).Error
}

func (fl *FilterLog) UpdateLatestLogID(latestLogID uint64) error {
	return db.GetDB().Model(&FilterLog{}).Where(fl).Updates(map[string]interface{}{
		"latest_log_id": latestLogID,
	}).Error
}

func (fl *FilterLog) First() (get *FilterLog, err error) {
	err = db.GetDB().Where(fl).First(&get).Error
	return get, err
}

func (fl *FilterLog) Last() (get *FilterLog, err error) {
	err = db.GetDB().Where(fl).Last(&get).Error
	return get, err
}

func (fl *FilterLog) Find(ext *QueryExtra, pager Pager) (list []*FilterLog, count int64, err error) {
	tx := db.GetDB().Where(fl)
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
		if err := tx.Model(fl).Count(&count).Error; err != nil {
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

func (fl *FilterLog) GetOldest10ByStatus(id uint64, action, stage, status uint8) (list []*FilterLog, err error) {
	err = db.GetDB().Where(fl).Where("id >= ?", id).Where("action = ?", action).
		Where("stage = ?", stage).Where("status = ?", status).Limit(OldestLimit).Find(&list).Error
	return list, err
}
