package wkteam

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
)

var (
	wkteamApiInstance = &WkteamApi{}
)

type WkteamApi struct {
	BaseUrl string
	AuthKey string

	client *HttpClient
}

func NewWkteamApi(baseUrl string, authKey string) *WkteamApi {
	wkteamApiInstance.BaseUrl = baseUrl
	wkteamApiInstance.AuthKey = authKey
	wkteamApiInstance.client = NewHttpClient()
	return wkteamApiInstance
}

// 获取登录二维码时候，返回的wId，手机完成登录后，通过此wId获取到登录的微信相关信息，
// 获取到的wcId 就是微信号的唯一ID，后续获取二维码时候带上
func (w *WkteamApi) GetInfoByWId(wId string) (data map[string]interface{}, err error) {
	requestMap := map[string]interface{}{
		"wId": wId,
	}

	resMap, err := w.doRequest("/getIPadLoginInfo", requestMap)
	if err != nil {
		return nil, err
	}

	return resMap.(map[string]interface{}), nil
}

func (w *WkteamApi) SetMsgReciverUrl(url string) error {
	requestMap := map[string]interface{}{
		"httpUrl": url,
		"type":    2, // 这里指定使用优化版协议
	}

	_, err := w.doRequest("/setHttpCallbackUrl", requestMap)
	if err != nil {
		return err
	}

	return nil
}

func (w *WkteamApi) SecondLogin(wId string) (data map[string]interface{}, err error) {
	requestMap := map[string]interface{}{
		"wId": wId,
	}

	resMap, err := w.doRequest("/secondLogin", requestMap)
	if err != nil {
		return nil, err
	}

	return resMap.(map[string]interface{}), nil
}

func (w *WkteamApi) GetAddressList(wId string) (data map[string]interface{}, err error) {
	requestMap := map[string]interface{}{
		"wId": wId,
	}
	resMap, err := w.doRequest("/getAddressList", requestMap)
	if err != nil {
		return nil, err
	}
	return resMap.(map[string]interface{}), err
}

func (w *WkteamApi) GetContact(wId string, wcIds []string) (data map[string]interface{}, err error) {
	requestMap := map[string]interface{}{
		"wId":  wId,
		"wcId": strings.Join(wcIds, ","),
	}
	resMap, err := w.doRequest("/getContact", requestMap)
	if err != nil {
		return nil, err
	}

	contactMap := make(map[string]interface{})
	for _, v := range resMap.([]interface{}) {
		contactMap[v.(map[string]interface{})["userName"].(string)] = v.(map[string]interface{})
	}
	return contactMap, nil
}

func (w *WkteamApi) IsOnline(wId string) (isOnline bool, err error) {
	requestMap := map[string]interface{}{
		"wId": wId,
	}
	resMap, err := w.doRequest("/isOnline", requestMap)

	if err != nil {
		return false, err
	}
	return resMap.(map[string]interface{})["isOnline"].(bool), err
}

func (w *WkteamApi) GetOnlineList() (data map[string]string, err error) {
	resMap, err := w.doRequest("/queryLoginWx", nil)
	if err != nil {
		return nil, err
	}
	onlineMap := make(map[string]string)
	for _, v := range resMap.([]interface{}) {
		onlineMap[v.(map[string]interface{})["wcId"].(string)] = v.(map[string]interface{})["wId"].(string)
	}

	return onlineMap, err
}

func (w *WkteamApi) InitAddressList(wId string) error {
	requestMap := map[string]interface{}{
		"wId": wId,
	}
	_, err := w.doRequest("/initAddressList", requestMap)

	return err
}

// at 注意传字符串
func (w *WkteamApi) SendText(wId string, wcId string, content string, at string) error {
	requestMap := map[string]interface{}{
		"wId":     wId,
		"wcId":    wcId,
		"content": content,
	}
	if len(at) > 0 {
		requestMap["at"] = at
	}
	_, err := w.doRequest("/sendText", requestMap)

	return err
}

func (w *WkteamApi) GetLoginPicUrl(wcId string) (data map[string]interface{}, err error) {
	requestMap := map[string]interface{}{
		"wcId":  wcId,
		"proxy": 1,
	}
	resMap, err := w.doRequest("/iPadLogin", requestMap)

	if err != nil {
		return nil, err
	}
	return resMap.(map[string]interface{}), err
}

func (w *WkteamApi) GetAuthKey(account string, password string) (authKey string, err error) {
	if len(w.AuthKey) > 0 {
		return w.AuthKey, nil
	}

	requestMap := map[string]interface{}{
		"account":  account,
		"password": password,
	}

	resMap, err := w.doRequest("/member/login", requestMap)
	if err != nil {
		return "", err
	}

	return resMap.(map[string]interface{})["Authorization"].(string), nil
}

func (w *WkteamApi) doRequest(uri string, data map[string]interface{}) (response interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("do request recover, err is %+v", err)
		}
	}()

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	if len(w.AuthKey) > 0 {
		header["Authorization"] = w.AuthKey
	}
	postOptions := PostOptions{
		Header: header,
	}

	res, err := w.client.PostJson(w.BaseUrl+uri, data, postOptions)
	if err != nil {
		log.Printf("do request err: %+v", err)
		return nil, err
	}

	resMap := make(map[string]interface{})
	err = json.Unmarshal(res, &resMap)
	if err != nil {
		log.Printf("do request json decode err: %+v", err)
		return nil, err
	}

	if resMap["code"] != "1000" {
		return nil, errors.New("get authkey response code err")
	}

	return resMap["data"], nil
}
