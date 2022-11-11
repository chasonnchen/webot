package wechat

import (
	"log"
	"time"

	"github.com/chasonnchen/webot/lib/wkteam"
)

var (
	wechatInstance = &Wechat{}
)

type Wechat struct {
	WId       string
	Options   WechatOptions
	WkteamApi *wkteam.WkteamApi
}
type WechatOptions struct {
	Account  string // 平台登录手机号
	Password string // 平台登录密码
	BaseUrl  string // 请求wkteam的域名信息
	WcId     string // 要登录的微信的微信ID
	AuthKey  string // 使用平台用户名密码获取到的authkey，请求后面接口时都需要带到header里面
}

func GetWechatInstance() *Wechat {
	return wechatInstance
}

func NewWechat(options WechatOptions) *Wechat {
	wechatInstance.Options = options
	wechatInstance.WkteamApi = wkteam.NewWkteamApi(options.BaseUrl, options.AuthKey)

	return wechatInstance
}

func (w *Wechat) Start() error {
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

	// 初始化通讯录
	err = w.WkteamApi.InitAddressList(w.WId)
	if err != nil {
		log.Printf("init address list err.")
		return err
	}

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
