package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"rmq/lib"
	"time"
)

// SendEmail 发送邮件
// 没有消息的时候就一直阻塞
func SendEmail(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		fmt.Printf("向user:%s发送邮件\n", string(msg.Body))
		time.Sleep(time.Second * 2) // 模拟发送邮件
		msg.Ack(false)
	}
}

func main() {
	mq := lib.NewMQ()
	mq.Consume(lib.QueueRegister, "c1", SendEmail)
	defer mq.Channel.Close()
}
