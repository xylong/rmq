package model

import (
	"github.com/jinzhu/gorm"
	"rmq/app"
)

// User 用户
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null;" json:"name" form:"name"`
	Avatar   string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);not null;unique" json:"email" form:"email"`
	Password string `gorm:"type:varchar(150);not null;" json:"password" form:"password"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Create(user *User) (uint, error) {
	if err := app.GetDB().Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}
