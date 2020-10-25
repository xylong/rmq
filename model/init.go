package model

import "rmq/app"

func init() {
	app.GetDB().AutoMigrate(&User{})
}
