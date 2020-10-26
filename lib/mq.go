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
	Channel       *amqp.Channel
	notifyConfirm chan amqp.Confirmation
	notifyReturn  chan amqp.Return
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
	// mandatory 为true时，发送失败将消息返还生产者
	return mq.Channel.Publish(exchange, key, true, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}

// Consume 消费消息
// queue 队列
// consumer 消费者
func (mq *MQ) Consume(queue, consumer string, callback func(<-chan amqp.Delivery, string)) {
	if msgs, err := mq.Channel.Consume(queue, consumer, false, false, false, false, nil); err != nil {
		log.Fatal(err)
	} else {
		callback(msgs, consumer)
	}
}

// SetConfirm 发送设置confirm模式
func (mq *MQ) SetConfirm() {
	if err := mq.Channel.Confirm(false); err != nil {
		log.Fatal(err)
	}
	mq.notifyConfirm = mq.Channel.NotifyPublish(make(chan amqp.Confirmation))
	go mq.ListenConfirm()
}

func (mq *MQ) ListenConfirm() {
	defer mq.Channel.Close()
	res := <-mq.notifyConfirm
	if res.Ack {
		log.Println("send success")
	} else {
		log.Println("send fail")
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

// NotifyReturn 入队失败回执
func (mq *MQ) NotifyReturn() {
	mq.notifyReturn = mq.Channel.NotifyReturn(make(chan amqp.Return))
	go mq.ListenReturn() // 协程执行
}

// ListenReturn 监听回执
func (mq *MQ) ListenReturn() {
	res := <-mq.notifyReturn
	msg := string(res.Body)
	if msg != "" {
		log.Println("消息没有正确入列:", msg)
	}
}
