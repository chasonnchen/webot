package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/chasonnchen/webot/configs"
	"github.com/chasonnchen/webot/dao"
	"github.com/chasonnchen/webot/logic"
	"github.com/chasonnchen/webot/service"
	"github.com/chasonnchen/webot/task"
	"github.com/chasonnchen/webot/wechat"
)

func main() {
	// 0. 初始化
	configs.InitConfig()
	webotConf := configs.GetConf()

	// 1. 启动bot
	wechatBot := wechat.NewWechat(wechat.WechatOptions{
		Account:  webotConf.Wkteam.Account,
		Password: webotConf.Wkteam.Password,
		BaseUrl:  webotConf.Wkteam.BaseUrl,
		AuthKey:  webotConf.Wkteam.AuthKey,
		WcId:     webotConf.Wkteam.WcId, // 要登录的微信ID（每个微信号 唯一且不变)
		Url:      webotConf.Wkteam.Reciver.Url,
		Port:     webotConf.Wkteam.Reciver.Port,
		Uri:      webotConf.Wkteam.Reciver.Uri,
	})
	wechatBot.Start()

	// 2. 初始化服务时候 用到了wcid 所以放到start后面了
	dao.InitDao()
	service.InitService()

	// 3. 注册核心消息处理
	// 事件ID参考 http://www.wkteam.cn/api-wen-dang2/xiao-xi-jie-shou/shou-xiao-xi/callback.html
	wechatBot.On("60001", onMessage) // 私聊消息
	wechatBot.On("80001", onMessage) // 群聊消息
	wechatBot.On("60010", onMessage) // 私聊小程序消息
	wechatBot.On("80010", onMessage) // 群聊小程序消息
	wechatBot.On("60007", onMessage) // 私聊h5消息
	wechatBot.On("80007", onMessage) // 群聊h5消息

	wechatBot.On("85008", onRoomJoin) // 群变更消息
	wechatBot.On("85009", onRoomJoin) // 群变更消息

	// 4. 定时任务初始化
	task.InitTask()

	var quitSig = make(chan os.Signal)
	signal.Notify(quitSig, os.Interrupt, os.Kill)
	select {
	case <-quitSig:
		log.Fatal("exit.by.signal")
	}
}

func onMessage(msg *wechat.Msg) {
	logic.NewMessageLogic().Do(msg)
}
func onRoomJoin(msg *wechat.Msg) {
	logic.NewRoomLogic().DoJoin(msg)
}
