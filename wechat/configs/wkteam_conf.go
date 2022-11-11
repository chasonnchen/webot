package configs

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type WkteamConf struct {
	Name     string   `yaml:"name"`
	Baseinfo Baseinfo `yaml:"baseinfo"`
}

type Baseinfo struct {
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
	BaseUrl  string `yaml:"baseurl"`
	AuthKey  string `yaml:"authkey"`
	WcId     string `yaml:"wcid"`
}

var (
	wkteamConfig *WkteamConf
	once         sync.Once
)

func GetConf() *WkteamConf {
	initWkteamConf()
	return wkteamConfig
}

func initWkteamConf() {
	once.Do(func() {
		yamlFile, err := ioutil.ReadFile("./conf/wkteam.yml")
		if err != nil {
			fmt.Println(err.Error())
		}
		var confTmp WkteamConf
		err = yaml.Unmarshal(yamlFile, &confTmp)
		if err != nil {
			fmt.Println(err.Error())
		}
		wkteamConfig = &confTmp
	})
}
