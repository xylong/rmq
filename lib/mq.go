package lib

import (
	"github.com/streadway/amqp"
	"log"
	"rmq/app"
)

// MQ RabbitMQ
type MQ struct {
	// Channel rabbitMQ的channel
	Channel *amqp.Channel
	// notifyConfirm 发送时的确认通知消息
	notifyConfirm chan amqp.Confirmation
	// notifyReturn 入队失败，返还通知
	notifyReturn chan amqp.Return
}

// NewMQ 创建mq
func NewMQ() *MQ {
	if c, err := app.GetConn().Channel(); err != nil {
		log.Println(err)
		return nil
	} else {
		return &MQ{
			Channel: c,
		}
	}
}

// DeclareAndBind
// key 路由
// exchange 交换机
// queues 队列
func (mq *MQ) DeclareAndBind(key, exchange string, queues ...string) error {
	for _, queue := range queues {
		if q, err := mq.Channel.QueueDeclare(queue, false, false, false, false, nil); err == nil {
			if err = mq.Channel.QueueBind(q.Name, key, exchange, false, nil); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// SetConfirm 开启发送确认
// 开启此功能会消耗性能
func (mq *MQ) SetConfirm() {
	if err := mq.Channel.Confirm(false); err != nil {
		log.Println(err)
	} else {
		mq.notifyConfirm = mq.Channel.NotifyPublish(make(chan amqp.Confirmation))
	}
}

// ListenConfirm 监听确认消息
func (mq *MQ) ListenConfirm() {
	defer mq.Channel.Close()
	res := <-mq.notifyConfirm
	if res.Ack {
		log.Println("send succeed")
	} else {
		log.Println("send failed")
	}
}

// NotifyReturn 入队返还通知
// 阻塞
func (mq *MQ) NotifyReturn() {
	mq.notifyReturn = mq.Channel.NotifyReturn(make(chan amqp.Return))
	go mq.listenReturn() // 协程执行
}

// listenReturn 监听入队返回通知
func (mq *MQ) listenReturn() {
	notify := <-mq.notifyReturn
	if msg := string(notify.Body); msg != "" {
		log.Printf("message:\"%s\" not listed correctly", msg)
	}
}

// Send 发送消息
// exchange 交换机
// key 路由
// message 消息
func (mq *MQ) Send(exchange, key, message string) error {
	// mandatory 若为true，在exchange正常且可以到达的情况下，如果exchange+routerKey无法将消息投递给queue，那么MQ会将消息返还给生产者
	//           若为false，则直接丢弃
	return mq.Channel.Publish(exchange, key, true, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

// SendDelay 发送延迟消息
// delay 延迟时间 (毫秒)
func (mq *MQ) SendDelay(exchange, key, message string, delay int) error {
	return mq.Channel.Publish(exchange, key, true, false, amqp.Publishing{
		Headers: map[string]interface{}{
			"x-delay": delay,
		},
		ContentType: "text/plain",
		Body:        []byte(message),
	})
}

// Receive 接收消息
// queue 队列名
// consumer 消费者
// callback 收到消息后执行的回调函数
func (mq *MQ) Receive(queue, consumer string, callback func(<-chan amqp.Delivery, string)) {
	if msg, err := mq.Channel.Consume(queue, consumer, false, false, false, false, nil); err != nil {
		log.Fatal(err)
	} else {
		callback(msg, consumer)
	}
}
