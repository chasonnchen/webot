package service

import (
	"log"
	"strings"
	//"time"
	"encoding/xml"
	"errors"

	//"github.com/chasonnchen/webot/dao"
	"github.com/chasonnchen/webot/entity"
	"github.com/chasonnchen/webot/wechat"

	"github.com/irebit/jingdong_union_go"
)

var (
	cpsServiceInstance = &CpsService{}
)

func NewCpsService() *CpsService {
	return cpsServiceInstance
}

func (c *CpsService) init() {
	c.app = &jingdong_union_go.App{
		ID:     "1000001619",
		Key:    "44b7c53b0e44177f66c32adb13f7585f",
		Secret: "336e29f137b54648a7309f4b6b955988",
	}
}

type CpsService struct {
	app *jingdong_union_go.App
}

type Msg struct {
	XMLName xml.Name `xml:"msg"`
	Appmsg  Appmsg   `xml:"appmsg"`
}
type Appmsg struct {
	XMLName xml.Name `xml:"appmsg"`
	Title   string   `xml:"title"`
	Desc    string   `xml:"desc"`
	Url     string   `xml:"url"`
}

func (c *CpsService) neesIgnore(msg *wechat.Msg) bool {
	// 过滤掉自己发出去的消息
	if msg.Data.Self {
		return true
	}

	return false
}

func (c *CpsService) checkUrl(url string) (string, error) {
	// 获得URL
	urlOri := strings.Split(url, "?")[0]
	log.Printf("Cps service recive url: %s", urlOri)

	// https://github.com/irebit/jingdong_union_go
	// 判断如果是京东的，获取返利链接关返回
	if strings.HasPrefix(urlOri, "http") && strings.Contains(urlOri, ".jd.") {
		return urlOri, nil
	}
	return url, errors.New("这个不是京东链接")
}

func (c *CpsService) GetCpsLinkByUrl(url string) (string, error) {

	// 获得URL
	urlOri := strings.Split(url, "?")[0]
	log.Printf("Cps service recive url: %s", urlOri)
	resContent := ""
	// https://github.com/irebit/jingdong_union_go
	// 判断如果是京东的，获取返利链接关返回
	if strings.HasPrefix(urlOri, "http") && strings.Contains(urlOri, ".jd.") {
		res, err := c.app.JdUnionOpenPromotionCommonGet(map[string]interface{}{
			//res, err := c.app.JdUnionOpenPromotionBysubunionidGet(map[string]interface{}{
			"materialId": urlOri,
			"siteId":     "4101166550",
		})
		if err != nil {
			return "", errors.New("请求京东联盟出错")
		}
		if res.Code == 200 {
			resContent = "已生成专属链接，点击后跳转到京东购买\n" + res.Data.ClickURL
		} else {
			resContent = "哎呀，转链失败，管理员在路上了~\n错误提示：" + res.Message
		}
		return resContent, nil
	}

	return "", errors.New("不是京东链接")
}

func (c *CpsService) DoCps(contact entity.ContactEntity, message *wechat.Msg) error {
	// 1. 检查是否需要忽略
	if c.neesIgnore(message) {
		return errors.New("ignor")
	}
	inputContent := ""

	// 判断是否指定群组
	if contact.Type == 2 && !NewGroupService().HasContact(3, contact.Id) {
		log.Printf("Cps service 不是指定群")
		return errors.New("不是指定群")
	}

	// 不同消息类型，特殊处理
	if message.MessageType == "60010" || message.MessageType == "80010" || message.MessageType == "60007" || message.MessageType == "80007" {
		// 解析xml
		msg := Msg{}
		xml.Unmarshal([]byte(message.Data.Content), &msg)
		url, err := c.checkUrl(msg.Appmsg.Url)

		// 特殊逻辑，对fenglinyexing展示原始链接
		if message.Data.FromUser == "fenglinyexing" {
			NewContactService().SayTextToContact(contact.Id, "原始链接："+url, "")
		}
		if err != nil {
			NewContactService().SayTextToContact(contact.Id, "亲分享链接不太对哦，在京东App直接分享到群里即可。（京东微信小程序里分享可能无法识别）", "")
			return nil
		}
		inputContent = url
	}
	if message.MessageType == "60001" || message.MessageType == "80001" {
		url, err := c.checkUrl(message.Data.Content)
		if err != nil {
			return errors.New("不是京东链接")
		}
		inputContent = url
	}

	resContent, err := c.GetCpsLinkByUrl(inputContent)
	if err == nil {
		NewContactService().SayTextToContact(contact.Id, resContent, "")
	}

	return nil
}
