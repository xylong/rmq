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
		//if c == "c1" {
		//	msg.Reject(true) // 重新入列
		//	continue         // 模拟c1故障
		//}
		fmt.Printf("%s receive msg:%s\n", c, string(msg.Body))
		go send(c, msg)
	}
}

// send 模拟邮件发送过程
func send(c string, msg amqp.Delivery) error {
	time.Sleep(time.Second * 3) // 模拟耗时
	fmt.Printf("%s send email to user:%s\n", c, string(msg.Body))
	msg.Ack(false)
	return nil
}

func main() {
	var c *string
	c = flag.String("c", "", "消费者")
	flag.Parse()
	if *c == "" {
		log.Fatal("c不能为空")
	}

	mq := lib.NewMQ()
	// 限流
	if err := mq.Channel.Qos(2, 0, false); err != nil {
		log.Fatal(err)
	}
	mq.Consume(lib.QueueRegister, *c, SendEmail)
	defer mq.Channel.Close()
}
