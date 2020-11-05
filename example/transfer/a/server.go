package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
	"rmq/example/transfer/service"
	"rmq/lib"
)

func init() {
	transfer.InitDB("a")
	transfer.GetDB().AutoMigrate(&model.AMoney{}, &model.TransLog{})
}

func main() {
	router := gin.Default()
	router.Use(transfer.ErrorMiddleware())
	router.Handle(http.MethodPost, "/", func(context *gin.Context) {
		ts := model.NewTrans()
		if err := context.ShouldBindJSON(ts); err == nil {
			err = service.Transfer(ts)
			transfer.CheckError(err, "A转账失败：")

			mq := lib.NewMQ()
			jsonStr, _ := json.Marshal(ts)
			if err := mq.Send(service.TransExchange, service.TransrRouter, string(jsonStr)); err != nil {
				log.Println(err)
			}

			context.JSON(http.StatusOK, gin.H{
				"data": ts.String(),
			})
		} else {
			context.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	})

	c := make(chan error)

	go func() {
		if err := router.Run(":8080"); err != nil {
			c <- err
		}
	}()

	go func() {
		if err := service.TransInit(); err != nil {
			c <- err
		}
	}()

	<-c
}
