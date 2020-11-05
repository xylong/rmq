package model

import (
	"rmq/app"
	"time"
)

// UserRegisterNotify 用户通知
type UserRegisterNotify struct {
	Uid       int       `json:"uid" gorm:"type:int(11);primary_key"`
	Num       int       `json:"num" gorm:"type:tinyint(1);default:1;"` // 失败次数
	Status    bool      `json:"status" gorm:"type:bit(1);default:0;"`  // 是否结束 1结束
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime;"`
}

// NewUserRegisterNotify 创建注册记录
func NewUserRegisterNotify() *UserRegisterNotify {
	return &UserRegisterNotify{}
}

// Log 注册记录
// id 用户id
// max 最大重试次数
func (notify *UserRegisterNotify) Log(id string, max int) int64 {
	return app.GetDB().Exec("INSERT INTO `user_register_notifies`(`uid`,`updated_at`) VALUE(?,NOW()) ON DUPLICATE KEY"+
		" UPDATE `num`=IF(`status`=1,num,num+1),`status`=IF(`num`>=?,1,0),`updated_at`=IF(`status`=1,`updated_at`,NOW())",
		id, max).RowsAffected
}
