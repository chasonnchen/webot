package service

import (
	"log"
)

func InitService() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("service init  recover, err is %+v", err)
		}
	}()

	// 初始化webot基础配置
	NewCacheService().init()
	NewUserService().init()
	NewGroupService().init()
	NewContactService().init()
	NewQaService().init()
	NewCpsService().init()
	// NewForwardService().init()
	// NewForwardMediaService().init()
	// NewRoomService().init()
}
