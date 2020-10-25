package lib

import (
	"log"
	"rmq/app"

	"github.com/streadway/amqp"
)

// MQ RabbitMQ
type MQ struct {
	Channel *amqp.Channel
}

// NewMQ 创建mq实例
func NewMQ() *MQ {
	var (
		err error
		c   *amqp.Channel
	)
	if c, err = app.GetConn().Channel(); err != nil {
		log.Println(err)
		return nil
	}
	return &MQ{
		Channel: c,
	}
}

// Send 发送消息
func (mq *MQ) Send(queue, msg string) (err error) {
	if _, err = mq.Channel.QueueDeclare(queue, false, false, false, false, nil); err != nil {
		return err
	}

	return mq.Channel.Publish("", queue, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}
