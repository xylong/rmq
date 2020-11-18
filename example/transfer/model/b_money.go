package model

import "github.com/jinzhu/gorm"

type BMoney struct {
	gorm.Model
	Tid      int    `json:"tid"` // a的交易号
	UserName string `json:"user_name" gorm:"type:varchar(20)"`
	Money    int    `json:"money"`
}
