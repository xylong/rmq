package service

import (
	"fmt"
	"rmq/lib"
)

const (
	RegisterQueue       = "user_register"       // 用户注册
	RegisterUnionrQueue = "union_user_register" // 合作单位用户注册
	UserExchange        = "user_exchange"       // 用户模块交换机
	UserDelayExchange   = "user_delay_exchange" // 用户延迟模块交换机
	RegisterRouter      = "user_register"       // 用户注册路由
)

// UserInit 初始化用户相关队列
func UserInit() error {
	mq := lib.NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	// 申明交换机
	if err := mq.Channel.ExchangeDeclare(UserExchange, "direct", false, false, false, false, nil); err != nil {
		return fmt.Errorf("exchange error:%s", err)
	}
	// 绑定队列
	if err := mq.DeclareAndBind(RegisterRouter, UserExchange, RegisterQueue, RegisterUnionrQueue); err != nil {
		return fmt.Errorf("queue bind error:%s", err)
	}

	return nil
}

// UserDelayInit 用户延迟队列初始化
func UserDelayInit() error {
	mq := lib.NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	// 申明交换机
	if err := mq.Channel.ExchangeDeclare(UserDelayExchange, "x-delayed-message", false, false, false, false, map[string]interface{}{
		"x-delayed-type": "direct",
	}); err != nil {
		return fmt.Errorf("delay exchange error:%s", err)
	}
	// 绑定队列
	if err := mq.DeclareAndBind(RegisterRouter, UserDelayExchange, RegisterQueue); err != nil {
		return fmt.Errorf("queue bind error:%s", err)
	}

	return nil
}
