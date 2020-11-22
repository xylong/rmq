package main

import (
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"rmq/app"
	"rmq/example/asyncorder/model"
	"rmq/example/asyncorder/service"
	"rmq/lib"
	"time"
)

type orderReq struct {
	Uid   int    // 用户🆔
	No    string // 订单号
	Money int    // 金额
}

func init() {
	app.GetDB().AutoMigrate(&model.Order{})
}

// getID 获取🆔
func getID() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}
	return time.Now().Format("20060102") + node.Generate().String(), nil
}

// generateUid 生成uid
func generateUid() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(cors())
	// 用户下单
	router.Handle(http.MethodPost, "/", func(context *gin.Context) {
		id, err := getID()

		mq := lib.NewMQ()
		if mq == nil {
			log.Println(err)
			return
		}
		defer mq.Channel.Close()

		req := &orderReq{
			Uid:   generateUid(),
			No:    id,
			Money: generateUid(),
		}
		str, _ := json.Marshal(req)
		if err := mq.Send(service.OrderExchange, service.OrderRouter, string(str)); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "create order fail",
			})
		}

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": "order fail",
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"no": id,
			})
		}
	})

	router.Handle(http.MethodGet, "/result", func(context *gin.Context) {
		orderNo := context.Query("no")
		order := model.NewOrder()
		if err := app.GetDB().Where("no=?", orderNo).Select("id").First(order).Error; err != nil {
			context.JSON(http.StatusOK, gin.H{"result": 0})
		} else {
			context.JSON(http.StatusOK, gin.H{"result": order.ID})
		}
	})

	c := make(chan error)

	go func() {
		if err := router.Run(); err != nil {
			c <- err
		}
	}()

	go func() {
		if err := service.OrderInit(); err != nil {
			c <- err
		}
	}()

	// 消费者
	go func() {
		client := lib.NewMQ()
		if client == nil {
			return
		}
		defer client.Channel.Close()

		if err := client.Channel.Qos(5, 0, false); err != nil {
			log.Fatal(err)
		}
		client.Receive(service.OrderQueue, "c1", service.CreateOrder)
	}()

	log.Fatal(<-c)
}
