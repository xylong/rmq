package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rmq/example/transfer"
	"rmq/example/transfer/model"
)

func init() {
	fmt.Println("b")
	transfer.InitDB("b")
	transfer.GetDB().AutoMigrate(&model.BMoney{}, &model.BTransLog{})
}

func main() {
	router := gin.Default()

	router.Handle(http.MethodPost, "/", func(context *gin.Context) {
		ts := model.NewTrans()
		err := context.BindJSON(&ts)
		if err != nil {
			context.JSON(200, gin.H{"result": err.Error()})
		} else {
			context.JSON(200, gin.H{"result": ts.String()})
		}
	})

	c := make(chan error)

	go func() {
		err := router.Run(":8081")
		if err != nil {
			c <- err
		}
	}()

	log.Fatal(<-c)
}
