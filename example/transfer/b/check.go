package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
	"rmq/example/transfer/service"
	"rmq/lib"
	"strings"
)

var client *lib.MQ

// receiveFromA 从a接收消息
func receiveFromA(messages <-chan amqp.Delivery, c string) {
	for msg := range messages {
		ts := model.NewTrans()
		err := json.Unmarshal(msg.Body, ts)
		if err != nil {
			log.Println(err)
		} else {
			go func(t *model.Trans) {
				defer msg.Ack(false)
				saveLog(t)
			}(ts)
		}
	}
}

func saveLog(ts *model.Trans) {
	tx := transfer.GetDB().Begin()

	// 日志
	if err := tx.Create(&model.BTransLog{
		Tid:   ts.ID,
		From:  ts.From,
		To:    ts.To,
		Money: ts.Money,
	}).Error; err != nil {
		log.Println(err)
		tx.Rollback()
	}

	// 加钱
	db := tx.Model(&model.BMoney{}).Where("user_name=?", ts.To).Update("money", gorm.Expr("money + ?", ts.Money))
	if db.Error != nil || db.RowsAffected == 0 {
		log.Println(db.Error)
		tx.Rollback()
	}

	// 回调a
	if err := callBack(ts.ID); err != nil {
		log.Println(err)
		tx.Rollback()
	}

	tx.Commit()
}

// callBack 回调
func callBack(tid int) error {
	rep, err := http.Post("http://127.0.0.1:8080/callback", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("tid=%d", tid)))
	if err != nil {
		return err
	}
	defer rep.Body.Close()
	bytes, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return err
	}
	if string(bytes) == "success" {
		return nil
	} else {
		return fmt.Errorf("fail")
	}
}

func main() {
	var consumer *string
	consumer = flag.String("c", "", "consumer")
	flag.Parse()
	if *consumer == "" {
		log.Fatal("consumer can't be null")
	}

	//c := make(chan error)
	transfer.InitDB("b")
	//go func() {
	//	if err := transfer.InitDB("b"); err != nil {
	//		c <- err
	//	}
	//}()

	client = lib.NewMQ()
	defer client.Channel.Close()

	// 收满10条消息后，ack应答之后才可以再接收消息
	if err := client.Channel.Qos(10, 0, false); err != nil {
		log.Println(err)
	}
	client.Receive(service.TransQueue, *consumer, receiveFromA)

	//<-c
}
