package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
	"log"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
	"rmq/example/transfer/service"
	"rmq/lib"
)

var (
	cro  *cron.Cron
	lock bool
)

// InitCron 初始化定时任务
func InitCron() error {
	cro = cron.New(cron.WithSeconds())
	if _, err := cro.AddFunc("0/3 * * * * *", CancelTransaction); err != nil {
		return err
	}
	if _, err := cro.AddFunc("0/4 * * * * *", Refund); err != nil {
		return err
	}
	if _, err := cro.AddFunc("0/2 * * * * *", Retransmission); err != nil {
		return err
	}
	return nil
}

// CancelTransaction 取消交易
// 定时执行
func CancelTransaction() {
	if err := transfer.GetDB().Model(&model.TransLog{}).
		Where("TIMESTAMPDIFF(SECOND,updated_at,now()) > 20 and status=0").
		Update("status", model.TradeFail).Error; err != nil {
		log.Fatal(err)
	}
}

// Refund 退款
func Refund() {
	if lock {
		fmt.Println("locked...")
		return
	}
	var (
		trans []model.Trans
		logs  []model.TransLog
	)
	tx := transfer.GetDB().Begin()
	lock = true // 加锁
	//time.Sleep(time.Second * 10)  // 测试下锁
	if err := tx.Where("status=2 and is_back=0").Select("id,`from`,money").Limit(10).Find(&logs).Scan(&trans).Error; err != nil {
		tx.Rollback()
	}

	for _, tran := range trans {
		// 退款
		if res := tx.Model(&model.AMoney{}).Where("user_name=?", tran.From).Update("money", gorm.Expr("money + ?", tran.Money)); res.Error != nil || res.RowsAffected == 0 {
			tx.Rollback()
		}
		// 退款状态
		if err := tx.Model(&model.TransLog{}).Where("id=?", tran.ID).Update("is_back", model.RefundYes).Error; err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	lock = false // 解锁
}

// Retransmission 重发
func Retransmission() {
	var (
		trans []model.Trans
		logs  []model.TransLog
	)
	if err := transfer.GetDB().Where("status=0 and TIMESTAMPDIFF(SECOND,updated_at,now()) <= 8").Select("*").Find(&logs).Scan(&trans).Error; err != nil {
		log.Println(err)
	} else {
		mq := lib.NewMQ()
		for _, t := range trans {
			go func(t *model.Trans) {
				str, _ := json.Marshal(t)
				if err := mq.Send(service.TransExchange, service.TransrRouter, string(str)); err != nil {
					log.Println(err)
				} else {
					log.Println("重发成功")
				}
			}(&t)
		}
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
