package configs

import (
	"fmt"
	"io/ioutil"
	"sync"

	"gopkg.in/yaml.v2"
)

type AppConf struct {
	Name   string             `yaml:"name"`
	DbList map[string]*DbConf `yaml:"db"`
	Baidu  Baidu              `yaml:"baidu"`
	Upload Upload             `yaml:"upload"`
}

type Upload struct {
	Path string `yaml:"path"`
}

type Baidu struct {
	Ak string `yaml:"ak"`
	Sk string `yaml:"sk"`
}

type DbConf struct {
	Dsn string `yaml:"dsn"`
}

var (
	appConfig *AppConf
	onceApp   sync.Once
)

func GetAppConf() *AppConf {
	initAppConf()
	return appConfig
}

// initAppConf 初始化配置
func initAppConf() {
	onceApp.Do(func() {
		yamlFile, err := ioutil.ReadFile("./conf/app.yml")
		if err != nil {
			fmt.Println(err.Error())
		}
		var confTmp AppConf
		err = yaml.Unmarshal(yamlFile, &confTmp)
		if err != nil {
			fmt.Println(err.Error())
		}
		appConfig = &confTmp
	})
}
