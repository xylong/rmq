package model

import "github.com/jinzhu/gorm"

const (
	TradeProgressing = 0 // 交易进行中
	TradeSuccess     = 1 // 交易成功
	TradeFail        = 2 // 交易失败
	RefundYes        = 1 // 已退款
	RefundNo         = 0 // 未退款
)

// TransLog 转账日志
type TransLog struct {
	gorm.Model
	From   string `json:"from" gorm:"type:varchar(20)"`
	To     string `json:"to" gorm:"type:varchar(20)"`
	Money  int    `json:"money"`
	Status int    `json:"status"`                   // 0待处理 1成功 2失败
	ISBack int    `json:"is_back" gorm:"default:0"` // 是否退款 0否 1是
}

func NewTransLog() *TransLog {
	return &TransLog{}
}
