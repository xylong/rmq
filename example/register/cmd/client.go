package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rmq/example/register/model"
	"rmq/example/register/service"
	"rmq/lib"
	"time"
)

func email(consumer string, message *amqp.Delivery) error {
	time.Sleep(time.Second * 2)

	delay := message.Headers["x-delay"] //  原有的延迟时间
	flag := true
	if flag {
		row := model.NewUserRegisterNotify().Log(string(message.Body), 3) // 最大重试3次
		if row > 0 {
			newDelay := int(delay.(int32) * 2) // 重发延迟时间延长
			fmt.Printf("%s sent email to user:%s failed,will try again in %d seconds \n",
				consumer, string(message.Body), newDelay)
			if err := client.SendDelay(service.UserDelayExchange, service.RegisterRouter, string(message.Body), newDelay); err != nil {
				log.Println(err)
			}
		} else {
			log.Println("达到最大次数，停止发送")
		}

		return message.Reject(false) // 丢弃消息
	} else {
		fmt.Printf("%s sent email to user:%s successfully \n", consumer, string(message.Body))
		return message.Ack(false) // 发送成功
	}
}

func sendEmail(messages <-chan amqp.Delivery, consumer string) {
	for message := range messages {
		// 模拟c1故障
		/*if consumer == "c1" {
			_ = message.Reject(true) // 拒接消息，true会重新入列
			continue
		}*/
		fmt.Printf("received:%s\n", message.Body)
		go email(consumer, &message)
	}
}

var client *lib.MQ

func main() {
	var c *string
	c = flag.String("c", "", "consumer")
	flag.Parse()

	if *c == "" {
		log.Fatal("consumer can't be null")
	}

	client = lib.NewMQ()
	defer client.Channel.Close()
	_ = client.Channel.Qos(10, 0, false) // 收满10条消息后，ack应答之后才可以再接收消息
	client.Receive(service.RegisterQueue, "c1", sendEmail)
}
