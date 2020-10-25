package app

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"os"
)

var (
	// Conn mq连接
	Conn *amqp.Connection
)

// Config 配置
type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
}

func init() {
	var (
		err    error
		config Config
	)

	dir, _ := os.Getwd()
	viper.SetConfigName("config")                       // 配置文件名
	viper.SetConfigType("toml")                         // 配置扩展名
	viper.AddConfigPath(fmt.Sprintf("%s/%s", dir, ".")) // 查找配置文件所在路径
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic("未找到配置文件")
		} else {
			// 配置文件被找到，但产生了另外的错误
			panic("配置文件内容错误")
		}
	}
	// 将配置读取到struct
	if err = viper.Unmarshal(&config); err != nil {
		fmt.Println(err.Error())
	}
	// 连接mq
	if Conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.User, config.Password, config.IP, config.Port)); err != nil {
		panic(err)
	}
}

// GetConn 获取mq连接
func GetConn() *amqp.Connection {
	return Conn
}
