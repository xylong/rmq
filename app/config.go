package app

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var (
	mq     mqConfig
	dbConf *dbConfig
)

type mqConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
}

type dbConfig struct {
	Driver          string
	Name            string
	IP              string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifeTime time.Duration
	LogEnable       bool
	Charset         string
}

func readConfig() {
	dir, _ := os.Getwd()
	viper.SetConfigName("config")                       // 配置文件名
	viper.SetConfigType("toml")                         // 配置扩展名
	viper.AddConfigPath(fmt.Sprintf("%s/%s", dir, ".")) // 查找配置文件所在路径

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
			panic("未找到配置文件")
		} else {
			// 配置文件被找到，但产生了另外的错误
			panic("配置文件内容错误")
		}
	}

	loadMq()
	loadDb()
}

func loadMq() {
	if err := viper.UnmarshalKey("mq", &mq); err != nil {
		log.Fatal(err)
	}
}

func loadDb() {
	if err := viper.UnmarshalKey("database", &dbConf); err != nil {
		log.Fatal(err)
	}
}
