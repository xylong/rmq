package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rmq/lib"
	"rmq/model"
	"strconv"
)

const QueueRegister = "register"

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
			mq := lib.NewMQ()
			if err = mq.Send(lib.RouterKeyUser, lib.ExchangeUser, strconv.Itoa(int(id))); err != nil {
				log.Println(err)
			}
			defer mq.Channel.Close()

			context.JSON(http.StatusOK, gin.H{
				"msg":  "ok",
				"data": user,
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"msg":  "fail",
				"data": id,
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
		if err := lib.UserInit(); err != nil {
			c <- err
		}
	}()
	log.Fatal(<-c)
}
