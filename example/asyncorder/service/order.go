package service

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rmq/app"
	"rmq/example/asyncorder/model"
	"rmq/lib"
	"time"
)

const (
	OrderExchange = "order_exchange"
	OrderQueue    = "order"
	OrderRouter   = "order"
)

// OrderInit 初始化订单队列
func OrderInit() error {
	mq := lib.NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	// 申明交换机
	if err := mq.Channel.ExchangeDeclare(OrderExchange, "direct", false, false, false, false, nil); err != nil {
		return fmt.Errorf("exchange error:%s", err)
	}
	// 绑定队列
	if err := mq.DeclareAndBind(OrderRouter, OrderExchange, OrderQueue); err != nil {
		return fmt.Errorf("queue bind error:%s", err)
	}
	return nil
}

// CreateOrder 订单入库
func CreateOrder(messages <-chan amqp.Delivery, consumer string) {
	for message := range messages {
		fmt.Printf("%s received:%s\n", consumer, message.Body)
		order := model.NewOrder()
		if err := json.Unmarshal(message.Body, order); err != nil {
			log.Fatal(err)
			message.Reject(false)
		} else {
			go func() {
				time.Sleep(time.Second * 5)
				if err := app.GetDB().Create(order).Error; err != nil {
					log.Fatal(err)
				}
			}()
			_ = message.Ack(false)
		}
	}
}
