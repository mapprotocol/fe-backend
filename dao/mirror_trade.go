package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameMirrorTrade = "mirror_trade"

type MirrorTrade struct {
	ID            uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	DepositSwapID uint64    `gorm:"column:deposit_swap_id;type:bigint(20)" json:"deposit_swap_id"`
	OrderViewId   string    `gorm:"column:order_view_id;type:varchar(255)" json:"order_view_id"`
	SrcChain      uint64    `gorm:"column:src_chain;type:bigint(20)" json:"src_chain"`
	SrcToken      string    `gorm:"column:src_token;type:varchar(255)" json:"src_token"`
	InAmount      string    `gorm:"column:in_amount;type:varchar(255)" json:"in_amount"`
	DstChain      uint64    `gorm:"column:dst_chain;type:bigint(20)" json:"dst_chain"`
	DstToken      string    `gorm:"column:dst_token;type:varchar(255)" json:"dst_token"`
	OutAmount     string    `gorm:"column:out_amount;type:varchar(255)" json:"out_amount"`
	Status        uint8     `gorm:"column:status;type:tinyint(4)" json:"status"`
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewMirrorTrade() *MirrorTrade {
	return &MirrorTrade{}
}

func (mt *MirrorTrade) TableName() string {
	return TableNameMirrorTrade
}

func (mt *MirrorTrade) Create() error {
	return db.GetDB().Create(mt).Error
}

func (mt *MirrorTrade) Updates(update *MirrorTrade) error {
	return db.GetDB().Where(mt).Updates(update).Error
}

func (mt *MirrorTrade) First() (get *MirrorTrade, err error) {
	err = db.GetDB().Where(mt).First(&get).Error
	return get, err
}

func (mt *MirrorTrade) Last() (get *MirrorTrade, err error) {
	err = db.GetDB().Where(mt).Last(&get).Error
	return get, err
}

func (mt *MirrorTrade) Find(ext *QueryExtra, pager Pager) (list []*MirrorTrade, count int64, err error) {
	tx := db.GetDB().Where(mt)
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
		if err := tx.Model(mt).Count(&count).Error; err != nil {
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
