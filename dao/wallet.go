package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"github.com/mapprotocol/ceffu-fe-backend/utils"
	"time"
)

const (
	WalletType = 10
)

const TableNameWallet = "wallet"

type Wallet struct {
	ID         uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	WalletID   int64     `gorm:"column:wallet_id;type:bigint(20)" json:"wallet_id"`
	WalletName string    `gorm:"column:wallet_name;type:varchar(255)" json:"wallet_name"`
	WalletType int32     `gorm:"column:wallet_type;type:int(11)" json:"wallet_type"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewWallet() *Wallet {
	return &Wallet{}
}

func (w *Wallet) TableName() string {
	return TableNameWallet
}

func (w *Wallet) Create() error {
	return db.GetDB().Create(w).Error
}

func (w *Wallet) Updates(update *Wallet) error {
	return db.GetDB().Where(w).Updates(update).Error
}

func (w *Wallet) First() (get *Wallet, err error) {
	err = db.GetDB().Where(w).First(&get).Error
	return get, err
}

func (w *Wallet) Last() (get *Wallet, err error) {
	err = db.GetDB().Where(w).Last(&get).Error
	return get, err
}

func (w *Wallet) Find(ext *QueryExtra, pager Pager) (list []*Wallet, count int64, err error) {
	tx := db.GetDB().Where(w)
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
		if err := tx.Model(w).Count(&count).Error; err != nil {
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
