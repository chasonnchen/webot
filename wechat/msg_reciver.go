package wechat

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initServer(port string, uri string) {
	server := gin.Default()
	server.POST(uri, recive)
	if err := server.Run(port); err != nil {
		log.Printf("startup gin server failed, err: %v\n", err)
		panic(err)
	}
}

type Response struct {
	Status int32       `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

type Msg struct {
	Account     string `json:"account"`
	WcId        string `json:"wcId"`
	MessageType string `json:"messageType"`
	Data        Data   `json:"data"`
}
type Data struct {
	FromUser  string `json:"fromUser"`
	FromGroup string `json:"fromGroup"`
	ToUser    string `json:"toUser"`
	MsgId     int64  `json:"msgId"`
	NewMsgId  int64  `json:"newMsgId"`
	Timestamp int64  `json:"timestamp"`
	Content   string `json:"content"`
	Self      bool   `json:"self"`
}

// 接收消息的逻辑实现入口
func recive(ctx *gin.Context) {
	msg := Msg{}
	ctx.BindJSON(&msg)
	log.Printf("recive msg: %+v", &msg)

	// 处理消息
	go GetWechatInstance().Handle(MsgName(msg.MessageType), &msg)
	ctx.JSON(http.StatusOK, Response{
		Status: 20000,
		Msg:    "success",
		Data:   msg,
	})

}
