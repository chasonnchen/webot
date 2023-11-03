package task

import (
	"log"
	"strconv"
	//"strings"
	"time"

	"github.com/chasonnchen/webot/lib/sina7x24"
	"github.com/chasonnchen/webot/wechat"
)

var (
	sina7x24Task = &Sina7x24Task{LastId: 0}
)

type Sina7x24Task struct {
	LastId int32
}

func NewSina7x24Task() *Sina7x24Task {
	return sina7x24Task
}

func (s *Sina7x24Task) Start() {
	s.work()
	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				s.work()
			}
		}
	}()
}

func (s *Sina7x24Task) work() {
	msg, id := sina7x24.NewClient().GetMsgs(0, s.LastId)
	if id > 0 {
		s.LastId = id
	}

	// 晚上10点半到早上8点半 不推送
	layout := "1504"
	timeStr, _ := strconv.Atoi(time.Now().Format(layout))
	if timeStr > 2230 || timeStr < 830 {
		//if timeStr > 10000 {
		//if !strings.Contains(msg, "俄") && !strings.Contains(msg, "乌") {
		log.Println("It is not good time")
		return
		//}
	}

	if len(msg) > 0 {
		wk := wechat.GetWechatInstance()
		wk.WkteamApi.SendText(wk.WId, "fenglinyexing", msg, "")
		time.Sleep(3 * time.Second)
		wk.WkteamApi.SendText(wk.WId, "liuzhaoliang-1", msg, "")
		time.Sleep(3 * time.Second)
		wk.WkteamApi.SendText(wk.WId, "pww932589183", msg, "")
		time.Sleep(3 * time.Second)
		wk.WkteamApi.SendText(wk.WId, "jaytudediaozha", msg, "")
	}
}
