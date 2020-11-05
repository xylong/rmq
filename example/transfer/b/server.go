package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
)

func init() {
	transfer.InitDB("b")
	transfer.GetDB().AutoMigrate(&model.Trans{}, &model.AMoney{})
}

func main() {
	router := gin.Default()
	router.Handle(http.MethodPost, "/", func(context *gin.Context) {

	})

	c := make(chan error)

	go func() {
		if err := router.Run(":8081"); err != nil {
			c <- err
		}
	}()

	<-c
}
