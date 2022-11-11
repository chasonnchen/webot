package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/chasonnchen/webot/configs"
	"github.com/chasonnchen/webot/task"
	"github.com/chasonnchen/webot/wechat"
)

func main() {
	// 0. 读配置
	configs.InitConfig()
	baseInfo := configs.GetConf().Baseinfo
	log.Printf("conf is %+v", baseInfo)

	// 1. 启动bot
	wechatBot := wechat.NewWechat(wechat.WechatOptions{
		Account:  baseInfo.Account,
		Password: baseInfo.Password,
		BaseUrl:  baseInfo.BaseUrl,
		AuthKey:  baseInfo.AuthKey,
		WcId:     baseInfo.WcId, // 要登录的微信ID（每个微信号 唯一且不变)
	})
	wechatBot.Start()

	// 2. 绑定核心事件处理

	// 3. task
	task.InitTask()

	var quitSig = make(chan os.Signal)
	signal.Notify(quitSig, os.Interrupt, os.Kill)
	select {
	case <-quitSig:
		log.Fatal("exit.by.signal")
	}
}
