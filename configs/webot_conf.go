package configs

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type WebotConf struct {
	Name   string `yaml:"name"`
	Wkteam Wkteam `yaml:"wkteam"`
}

type Wkteam struct {
	Account  string  `yaml:"account"`
	Password string  `yaml:"password"`
	BaseUrl  string  `yaml:"baseurl"`
	AuthKey  string  `yaml:"authkey"`
	WcId     string  `yaml:"wcid"`
	Reciver  Reciver `yaml:"reciver"`
}
type Reciver struct {
	Port string `yaml:"port"` // 本地gin监听的端口，可以直接对外暴漏，也可以走ngixn转发过来
	Uri  string `yaml:"uri"`  // 本地gin服务监听的URI
	Url  string `yaml:"url"`  // 上调到wkteam的回调URL地址
}

var (
	webotConfig *WebotConf
	once        sync.Once
)

func GetConf() *WebotConf {
	initWebotConf()
	return webotConfig
}

func initWebotConf() {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile("./conf/webot.yml")
		if err != nil {
			fmt.Println(err.Error())
		}
		var confTmp WebotConf
		err = yaml.Unmarshal(yamlFile, &confTmp)
		if err != nil {
			fmt.Println(err.Error())
		}
		webotConfig = &confTmp
	})
}
