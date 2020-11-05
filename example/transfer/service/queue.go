package service

import (
	"fmt"
	"rmq/lib"
)

// TransInit 转账队列初始化
func TransInit() error {
	mq := lib.NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	// 申明交换机
	if err := mq.Channel.ExchangeDeclare(TransExchange, "direct", false, false, false, false, nil); err != nil {
		return fmt.Errorf("exchange error:%s", err)
	}

	// 绑定队列
	if err := mq.DeclareAndBind(TransrRouter, TransExchange, TransQueue); err != nil {
		return fmt.Errorf("queue bind error:%s", err)
	}

	return nil
}
