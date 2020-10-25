package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"rmq/app"
)

func main() {
	var (
		err  error
		c    *amqp.Channel        // 频道
		msgs <-chan amqp.Delivery // 消息
	)

	conn := app.GetConn()
	defer conn.Close()

	if c, err = conn.Channel(); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if msgs, err = c.Consume("test", "c1", false, false, false, false, nil); err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		fmt.Println(msg.DeliveryTag, string(msg.Body))
	}
}
