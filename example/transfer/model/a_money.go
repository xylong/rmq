package model

import "github.com/jinzhu/gorm"

type AMoney struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"type:varchar(20)"`
	Money    int    `json:"money"`
}

func NewAMoney() *AMoney {
	return &AMoney{}
}
