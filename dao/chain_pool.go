package dao

import (
	"github.com/mapprotocol/fe-backend/resource/db"
	"github.com/mapprotocol/fe-backend/utils"
	"time"
)

const TableNameChainPool = "chain_pool"

type ChainPool struct {
	ID                 uint64    `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true" json:"id"`
	ChainID            string    `gorm:"column:chain_id;type:varchar(255)" json:"chain_id"`
	ChainName          string    `gorm:"column:chain_name;type:varchar(255)" json:"chain_name"`
	ChainRPC           string    `gorm:"column:chain_rpc;type:varchar(255)" json:"chain_rpc"`
	USDTContract       string    `gorm:"column:usdt_contract;type:varchar(255)" json:"usdt_contract"`
	WBTCContract       string    `gorm:"column:wbtc_contract;type:varchar(255)" json:"wbtc_contract"`
	FeRouterContract   string    `gorm:"column:fe_router_contract;type:varchar(255)" json:"fe_router_contract"`
	ChainPoolContract  string    `gorm:"column:chain_pool_contract;type:varchar(255)" json:"chain_pool_contract"`
	GasLimitMultiplier string    `gorm:"column:gas_limit_multiplier;type:varchar(255)" json:"gas_limit_multiplier"`
	CreatedAt          time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

func (*ChainPool) TableName() string {
	return TableNameChainPool
}

func NewChainPool() *ChainPool {
	return &ChainPool{}
}

func NewChainPoolWithChainID(chainID string) *ChainPool {
	return &ChainPool{ChainID: chainID}
}

func (cp *ChainPool) First() (get *ChainPool, err error) {
	err = db.GetDB().Where(cp).First(&get).Error
	return get, err
}

func (cp *ChainPool) Find(ext *QueryExtra, pager Pager) (list []*ChainPool, count int64, err error) {
	tx := db.GetDB().Where(cp)
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
		if err := tx.Model(cp).Count(&count).Error; err != nil {
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
