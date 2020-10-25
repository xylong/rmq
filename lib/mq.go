package lib

import (
	"github.com/streadway/amqp"
	"log"
	"rmq/app"
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
func (mq *MQ) Send(exchange, key, msg string) (err error) {
	return mq.Channel.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}

// Consume 消费消息
func (mq *MQ) Consume(queue, key string, callback func(<-chan amqp.Delivery)) {
	if msgs, err := mq.Channel.Consume(queue, key, false, false, false, false, nil); err != nil {
		log.Fatal(err)
	} else {
		callback(msgs)
	}
}

// DeclareQueueAndBind 申明队列并且绑定
// queues 多个队列逗号分隔
func (mq *MQ) DeclareQueueAndBind(exchange, key string, queues ...string) error {
	for _, queue := range queues {
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
