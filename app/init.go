package app

import (
	"fmt"
	"github.com/streadway/amqp"
)

var (
	// Conn mq连接
	Conn *amqp.Connection
)

func init() {
	var err error

	readConfig()
	// 连接mq
	if Conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		mq.User, mq.Password, mq.IP, mq.Port)); err != nil {
		panic(err)
	}

	// 连接数据库
	initDB()
}

// GetConn 获取mq连接
func GetConn() *amqp.Connection {
	return Conn
}
