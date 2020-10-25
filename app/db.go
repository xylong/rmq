package app

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(dbConf.Driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbConf.User, dbConf.Password, dbConf.IP, dbConf.Port, dbConf.Name, dbConf.Charset))
	if err != nil {
		log.Fatal(err)
	}
	db.DB().SetMaxIdleConns(dbConf.MaxIdleConn)
	db.DB().SetMaxOpenConns(dbConf.MaxOpenConn)
	db.DB().SetConnMaxLifetime(dbConf.MaxConnLifeTime)

	db.LogMode(dbConf.LogEnable)
}

func GetDB() *gorm.DB {
	return db
}
