package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
)

// Transfer 转账
func Transfer(trans *model.Trans) error {
	tx := transfer.GetDB().Begin()

	if row := tx.Model(&model.AMoney{}).Where("user_name=? and money>=?", trans.From, trans.Money).Update("money", gorm.Expr("money - ?", trans.Money)).RowsAffected; row == 0 {
		tx.Rollback()
		return fmt.Errorf("扣款失败")
	}

	if err := tx.Create(&model.TransLog{
		From:  trans.From,
		To:    trans.To,
		Money: trans.Money,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("日志记录失败")
	}

	tx.Commit()
	return nil
}
