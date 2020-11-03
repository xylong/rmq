package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rmq/example/register/service"
	"rmq/lib"
	"time"
)

func email(consumer string, message *amqp.Delivery) error {
	time.Sleep(time.Second * 3)
	fmt.Printf("%s send email to user:%s\n", consumer, string(message.Body))
	return message.Ack(false)
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

func main() {
	var c *string
	c = flag.String("c", "", "consumer")
	flag.Parse()

	if *c == "" {
		log.Fatal("consumer can't be null")
	}

	mq := lib.NewMQ()
	defer mq.Channel.Close()
	_ = mq.Channel.Qos(10, 0, false) // 收满10条消息后，ack应答之后才可以再接收消息
	mq.Receive(service.RegisterQueue, "c1", sendEmail)
}
