package model

import (
	"github.com/jinzhu/gorm"
)

// Order 订单
type Order struct {
	gorm.Model
	Uid   int    `json:"uid"`   // 用户id
	No    string `json:"no"`    // 订单号
	Money int    `json:"money"` // 金额
	State int    `json:"state"`
}

func NewOrder() *Order {
	return &Order{}
}
