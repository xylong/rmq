package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rmq/app"
	"rmq/example/register/model"
	"rmq/example/register/service"
	"rmq/lib"
	"strconv"
)

func init() {
	app.GetDB().AutoMigrate(
		&model.User{},
		&model.UserRegisterNotify{},
	)
}

func main() {
	router := gin.Default()

	router.Handle(http.MethodPost, "/users", func(context *gin.Context) {
		user := model.NewUser()
		if err := context.ShouldBind(user); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"msg": "参数错误",
			})
		}

		if id, err := user.Create(user); err == nil && id > 0 {
			// api返回
			context.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"data": user,
			})
			// 队列
			mq := lib.NewMQ()
			defer mq.Channel.Close()
			mq.SetConfirm() // 开启确认
			mq.NotifyReturn()
			//_ = mq.Send(service.UserExchange, service.RegisterRouter, strconv.Itoa(int(id)))
			_ = mq.SendDelay(service.UserDelayExchange, service.RegisterRouter, strconv.Itoa(int(id)), 1500)
			mq.ListenConfirm() // 监听是否发送成功
		} else {
			context.JSON(http.StatusOK, gin.H{
				"msg":  "fail",
				"data": user.ID,
			})
		}
	})

	c := make(chan error)

	go func() {
		if err := router.Run(); err != nil {
			c <- err
		}
	}()

	go func() {
		if err := service.UserInit(); err != nil {
			c <- err
		}
		if err := service.UserDelayInit(); err != nil {
			c <- err
		}
	}()

	log.Fatal(<-c)
}
