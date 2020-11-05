package model

import "github.com/jinzhu/gorm"

type BMoney struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"type:varchar(20)"`
	Money    int    `json:"money"`
}
