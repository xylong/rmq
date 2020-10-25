package lib

import (
	"log"
	"rmq/app"
	"strings"

	"github.com/streadway/amqp"
)

const (
	QueueRegister      = "register"
	QueueRegisterUnion = "register_union"
	ExchangeUser       = "user_exchange"
	RouterKeyUser      = "user_register"
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
func (mq *MQ) Send(key, exchange, msg string) (err error) {
	return mq.Channel.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}

// DeclareQueueAndBind 申明队列并且绑定
// queues 多个队列逗号分隔
func (mq *MQ) DeclareQueueAndBind(queues, key, exchange string) error {
	list := strings.Split(queues, ",")
	for _, queue := range list {
		q, err := mq.Channel.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			return err
		}

		if err = mq.Channel.QueueBind(q.Name, key, exchange, false, nil); err != nil {
			return err
		}
	}

	return nil
}
