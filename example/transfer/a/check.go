package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
)

var cro *cron.Cron

// InitCron 初始化定时任务
func InitCron() error {
	cro = cron.New(cron.WithSeconds())
	_, err := cro.AddFunc("0/3 * * * * *", CancelTransaction)
	return err
}

// CancelTransaction 取消交易
// 定时执行
func CancelTransaction() {
	if err := transfer.GetDB().Model(&model.TransLog{}).
		Where("TIMESTAMPDIFF(SECOND,updated_at,now()) > 20 and status<>2").
		Update("status", model.TradeFail).Error; err != nil {
		log.Fatal(err)
	}
}

func main() {
	c := make(chan error)

	go func() {
		if err := transfer.InitDB("a"); err != nil {
			c <- err
		}
	}()

	go func() {
		if err := InitCron(); err != nil {
			c <- err
		} else {
			cro.Start()
		}
	}()

	<-c
}
