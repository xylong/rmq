package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rmq/lib"
	"time"
)

// SendEmail 发送邮件
// 没有消息的时候就一直阻塞
func SendEmail(msgs <-chan amqp.Delivery, c string) {
	for msg := range msgs {
		fmt.Printf("%s向user:%s发送邮件\n", c, string(msg.Body))
		time.Sleep(time.Second * 2) // 模拟发送邮件
		if c == "c1" {
			msg.Reject(true) // 重新入列
			continue         // 模拟c1故障
		}
		msg.Ack(false)
	}
}

func main() {
	var c *string
	c = flag.String("c", "", "消费者")
	flag.Parse()
	if *c == "" {
		log.Fatal("c不能为空")
	}

	mq := lib.NewMQ()
	mq.Consume(lib.QueueRegister, *c, SendEmail)
	defer mq.Channel.Close()
}
