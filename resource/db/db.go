package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func Init(user, password, host, port, name string) {
	log.Println("connecting MySQL ... ", host)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	mdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
				Colorful: false,
			},
		),
		TranslateError: true,
	})
	if err != nil {
		panic(err)
	}
	if mdb == nil {
		panic("failed to connect database")
	}
	log.Println("connected database")
	db = mdb
	return
}

func GetDB() *gorm.DB {
	return db
}
