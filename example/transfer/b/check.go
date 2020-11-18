package main

import (
	"encoding/json"
	"flag"
	"github.com/streadway/amqp"
	"log"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
	"rmq/example/transfer/service"
	"rmq/lib"
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
	if err := transfer.GetDB().Create(&model.BTransLog{
		Tid:   ts.ID,
		From:  ts.From,
		To:    ts.To,
		Money: ts.Money,
	}).Error; err != nil {
		log.Println(err)
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
