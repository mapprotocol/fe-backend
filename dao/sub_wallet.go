package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const TableNameSubWallet = "sub_wallet"

type SubWallet struct {
	ID             uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	ParentWalletID int64     `gorm:"column:parent_wallet_id;type:bigint(20)" json:"parent_wallet_id"`
	WalletID       int64     `gorm:"column:wallet_id;type:bigint(20)" json:"wallet_id"`
	WalletName     string    `gorm:"column:wallet_name;type:varchar(255)" json:"wallet_name"`
	WalletType     int32     `gorm:"column:wallet_type;type:int(11)" json:"wallet_type"`
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewSubWallet() *SubWallet {
	return &SubWallet{}
}

func (sw *SubWallet) TableName() string {
	return TableNameSubWallet
}

func (sw *SubWallet) Create() error {
	return db.GetDB().Create(sw).Error
}

func (sw *SubWallet) Updates(update *SubWallet) error {
	return db.GetDB().Where(sw).Updates(update).Error
}

func (sw *SubWallet) First() (get *SubWallet, err error) {
	err = db.GetDB().Where(sw).First(&get).Error
	return get, err
}

func (sw *SubWallet) Last() (get *SubWallet, err error) {
	err = db.GetDB().Where(sw).Last(&get).Error
	return get, err
}

func (sw *SubWallet) Find(ext *QueryExtra, pager Pager) (list []*SubWallet, count int64, err error) {
	tx := db.GetDB().Where(sw)
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
		if err := tx.Model(sw).Count(&count).Error; err != nil {
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
