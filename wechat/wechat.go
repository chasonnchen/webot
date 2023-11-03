package wechat

import (
	"log"
	"reflect"
	"time"

	"github.com/chasonnchen/webot/entity"
	"github.com/chasonnchen/webot/lib/wkteam"
)

var (
	wechatInstance = &Wechat{}
)

type Wechat struct {
	WId         string
	WcId        string
	Options     WechatOptions
	WkteamApi   *wkteam.WkteamApi
	MsgHandler  *MsgHandler // 消息处理器
	ContactList map[string]entity.ContactEntity
}
type WechatOptions struct {
	Account  string // 平台登录手机号
	Password string // 平台登录密码
	BaseUrl  string // 请求wkteam的域名信息
	WcId     string // 要登录的微信的微信ID
	AuthKey  string // 使用平台用户名密码获取到的authkey，请求后面接口时都需要带到header里面
	Url      string // 消息回调URL，需要上报到wkteam，
	Port     string // 本地起的gin server服务监听的端口
	Uri      string // 接口URI
}

type MsgLogicFunc func(msg *Msg)

func GetWechatInstance() *Wechat {
	return wechatInstance
}

func NewWechat(options WechatOptions) *Wechat {
	wechatInstance.Options = options
	wechatInstance.WkteamApi = wkteam.NewWkteamApi(options.BaseUrl, options.AuthKey)
	wechatInstance.MsgHandler = GetMsgHandler()
	wechatInstance.WcId = options.WcId

	return wechatInstance
}

func (w *Wechat) On(msgName MsgName, handler MsgLogicFunc) {
	w.MsgHandler.On(msgName, func(data ...interface{}) {
		values := make([]reflect.Value, 0, len(data))
		for _, v := range data {
			values = append(values, reflect.ValueOf(v))
		}
		_ = reflect.ValueOf(handler).Call(values)
	})
	return
}

func (w *Wechat) Handle(msgName MsgName, data interface{}) {
	w.MsgHandler.Handle(msgName, data)
	return
}

func (w *Wechat) Start() error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("wechat start recover, err is %+v", err)
		}
	}()
	// init gin server, gin run会阻塞，所以这里单独一个协程
	go initServer(w.Options.Port, w.Options.Uri)

	// check options参数，如果缺少authkey wcid 等，需要初始化
	if w.Options.AuthKey == "" {
		authKey, err := w.WkteamApi.GetAuthKey(w.Options.Account, w.Options.Password)
		if err != nil {
			return err
		}
		w.Options.AuthKey = authKey
	}

	// login
	wxInfo, err := w.login()
	if err != nil {
		log.Printf("login err.")
		return err
	}
	log.Printf("[%s]登录成功! wcId[%s], wId[%s], 显示账号[%s]", wxInfo["nickName"].(string), wxInfo["wcId"].(string), wxInfo["wId"].(string), wxInfo["wAccount"].(string))

	// 设置接收消息地址
	err = w.WkteamApi.SetMsgReciverUrl(w.Options.Url)

	// 初始化通讯录
	err = w.WkteamApi.InitAddressList(w.WId)
	if err != nil {
		log.Printf("init address list err.")
		return err
	}

	// 持续更新通信录 暂定30min

	return nil
}

func (w *Wechat) reLogin() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("re login recover, err is %+v", err)
		}
	}()
	isOnline, _ := w.WkteamApi.IsOnline(w.WId)
	if isOnline == false {
		log.Printf("Online false,start relogin wid is %s", w.WId)
		isLogin := false
		for isLogin == false {
			select {
			case <-time.After(time.Second * 10):
				wxInfo, err := w.WkteamApi.SecondLogin(w.WId)
				if err == nil {
					isLogin = true
					w.WId = wxInfo["wcId"].(string)
					break
				}
				log.Printf("Online false,start relogin wid is %s", w.WId)
			}
		}
	}
}

func (w *Wechat) login() (wxInfo map[string]interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("login recover, err is %+v", err)
		}
	}()

	// 查在线的微信列表，判断是否在线
	onlineList, err := w.WkteamApi.GetOnlineList()
	if err != nil {
		return nil, err
	}
	wId, ok := onlineList[w.Options.WcId]
	w.WId = wId
	if ok {
		return w.WkteamApi.GetInfoByWId(wId)
	}

	// 获取登录二维码URL
	loginUrl, err := w.WkteamApi.GetLoginPicUrl(w.Options.WcId)

	// 循环获取用户信息，直到登录成功
	log.Printf("login QrCode url is %s", loginUrl["qrCodeUrl"].(string))
	isLogin := false
	for isLogin == false {
		select {
		case <-time.After(time.Second * 20):
			wxInfo, err = w.WkteamApi.GetInfoByWId(loginUrl["wId"].(string))
			if err == nil {
				isLogin = true
				w.WId = loginUrl["wId"].(string)
				break
			}
			log.Printf("login QrCode url is %s", loginUrl["qrCodeUrl"])
		}
	}

	// 后台检查是否登录，以及二次登录逻辑
	go func() {
		for {
			select {
			case <-time.After(time.Second * 300):
				w.reLogin()
			}
		}
	}()

	return wxInfo, nil
}
