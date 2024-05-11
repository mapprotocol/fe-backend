package dao

import "gorm.io/gorm"

type QueryExtra struct {
	Conditions        map[string]interface{}
	OnlyKeyConditions map[string]struct{}
	OrderStr          string
	DistinctStr       []string
	SelectFields      []string
}

type Pager func(*gorm.DB) *gorm.DB

func Paginate(page, pageSize int) Pager {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 20
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Paginate2(offset, limit int) Pager {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}
