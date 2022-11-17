package dao

import (
	"fmt"
	"sync"

	"github.com/chasonnchen/webot/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbMap map[string]*gorm.DB
	once  sync.Once
)

func InitDao() {
	once.Do(func() {
		dbMap = make(map[string]*gorm.DB)
		for k, v := range configs.GetAppConf().DbList {
			dbClient, err := gorm.Open(mysql.Open(v.Dsn), &gorm.Config{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			dbMap[k] = dbClient
		}
	})
}

func getDb(dbName string) *gorm.DB {
	return dbMap[dbName]
}

func Wechat() *gorm.DB {
	return getDb("wechat")
}
