package transfer

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// InitDB 初始化数据库
func InitDB(dbName string) (err error) {
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		"root", "root", "127.0.0.1", 3306, dbName, "utf8"))
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(100)

	db.LogMode(true)
	return
}

func GetDB() *gorm.DB {
	return db
}
