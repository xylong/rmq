package model

import "github.com/jinzhu/gorm"

// TransLog 转账日志
type TransLog struct {
	gorm.Model
	From   string `json:"from" gorm:"type:varchar(20)"`
	To     string `json:"to" gorm:"type:varchar(20)"`
	Money  int    `json:"money"`
	Status int    `json:"status"` // 0待处理 1成功 2失败
}

func NewTransLog() *TransLog {
	return &TransLog{}
}
