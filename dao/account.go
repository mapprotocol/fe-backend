package dao

import (
	"github.com/mapprotocol/ceffu-fe-backend/resource/db"
	"time"
)

const TableNameAccount = "account"

type Account struct {
	ID          uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	SubWalletID int64     `gorm:"column:sub_wallet_id;type:bigint(20)" json:"sub_wallet_id"`
	Network     string    `gorm:"column:wallet_name;type:varchar(255)" json:"wallet_name"`
	Address     string    `gorm:"column:address;type:varchar(255)" json:"address"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func NewAccount(subWalletID int64, network string) *Account {
	return &Account{
		SubWalletID: subWalletID,
		Network:     network,
	}
}

func (sw *Account) TableName() string {
	return TableNameAccount
}

func (sw *Account) Create() error {
	return db.GetDB().Create(sw).Error
}

func (sw *Account) Updates(update *Account) error {
	return db.GetDB().Where(sw).Updates(update).Error
}

func (sw *Account) First() (get *Account, err error) {
	err = db.GetDB().Where(sw).First(&get).Error
	return get, err
}
