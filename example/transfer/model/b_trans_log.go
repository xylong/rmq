package model

import (
	"time"
)

// BTransLog 转账日志
type BTransLog struct {
	Tid       int    `json:"tid" gorm:"type:int(11);primary_key"` // 交易号
	From      string `json:"from" gorm:"type:varchar(20)"`
	To        string `json:"to" gorm:"type:varchar(20)"`
	Money     int    `json:"money"`
	Status    int    `json:"status"`                   // 0待处理 1成功 2失败
	ISBack    int    `json:"is_back" gorm:"default:0"` // 是否退款 0否 1是
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func NewBTransLog() *BTransLog {
	return &BTransLog{}
}
