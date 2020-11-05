package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
	"rmq/example/transfer/service"
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
			_ = service.Transfer(ts)
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

	//go func() {
	//	if err := transfer.InitDB("a"); err != nil {
	//		c <- err
	//	}
	//}()

	<-c
}
