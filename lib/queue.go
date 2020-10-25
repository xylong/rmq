package lib

import "fmt"

// UserInit 初始化用户相关队列
func UserInit() error {
	mq := NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	// 申明交换机
	if err := mq.Channel.ExchangeDeclare(ExchangeUser, "direct", false, false, false, false, nil); err != nil {
		return err
	}
	qs := fmt.Sprintf("%s,%s", QueueRegister, QueueRegisterUnion)
	if err := mq.DeclareQueueAndBind(qs, RouterKeyUser, ExchangeUser); err != nil {
		return err
	}
	return nil
}
