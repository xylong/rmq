package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
)

const (
	TransExchange = "trans_exchange" // 转账交换机
	TransrRouter  = "trans"          // 转账路由
	TransQueue    = "trans_a"        // a的转账队列
)

// Transfer 转账
func Transfer(trans *model.Trans) error {
	tx := transfer.GetDB().Begin()

	if row := tx.Model(&model.AMoney{}).Where("user_name=? and money>=?", trans.From, trans.Money).Update("money", gorm.Expr("money - ?", trans.Money)).RowsAffected; row == 0 {
		tx.Rollback()
		return fmt.Errorf("扣款失败")
	}

	log := &model.TransLog{
		From:  trans.From,
		To:    trans.To,
		Money: trans.Money,
	}
	if tx.Create(log).Error != nil {
		tx.Rollback()
		return fmt.Errorf("日志记录失败")
	}
	trans.ID = int(log.ID) // 交易号
	tx.Commit()
	return nil
}
