package service

import (
	"log"
	"strings"
	"time"

	"github.com/chasonnchen/webot/dao"
	"github.com/chasonnchen/webot/entity"
	"github.com/chasonnchen/webot/wechat"
)

var (
	qaServiceInstance = &QaService{}
)

func NewQaService() *QaService {
	return qaServiceInstance
}

type QaService struct {
	QaConf map[string][]entity.SkillQaEntity
}

func (q *QaService) neesIgnore(msg *wechat.Msg) bool {
	// 过滤掉自己发出去的消息
	if msg.Data.Self {
		return true
	}

	return false
}

func (q *QaService) DoQa(contact entity.ContactEntity, message *wechat.Msg) {
	// 1. 检查是否需要忽略
	if q.neesIgnore(message) {
		return
	}

	// 2. 开始匹配问答
	for _, qaItem := range q.QaConf[contact.Id] {
		for _, keyword := range strings.Split(qaItem.QaKey, ",") {
			if strings.Contains(message.Data.Content, keyword) {
				if contact.Type == 2 {
					//TODO 实现群里回答时，at的逻辑
					NewContactService().SayTextToContact(contact.Id, strings.Trim(qaItem.QaValue, "\n"))
				} else {
					NewContactService().SayTextToContact(contact.Id, strings.Trim(qaItem.QaValue, "\n"))
				}
				log.Printf("Message response is %s", qaItem.QaValue)
			}
		}
	}

	return
}

func (q *QaService) init() {
	q.load()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 60):
				q.load()
			}
		}
	}()
}

func (q *QaService) load() {
	bot := wechat.GetWechatInstance()
	qaConf := make(map[string][]entity.SkillQaEntity)
	var qaListOri []entity.SkillQaEntity
	dao.Wechat().Where("status = ? and bot_id = ?", "1", bot.WcId).Find(&qaListOri)

	// 把组级别的配置，扩散成群配置
	var qaList []entity.SkillQaEntity
	for _, qaItem := range qaListOri {
		if len(qaItem.ContactId) > 1 {
			qaList = append(qaList, qaItem)
		}
		if qaItem.GroupId > 0 {
			contactIdList := NewGroupService().GetContactIdListByGroupId(qaItem.GroupId)
			for _, cid := range contactIdList {
				newQaItem := entity.SkillQaEntity{
					ContactId: cid,
					Name:      qaItem.Name,
					QaKey:     qaItem.QaKey,
					QaValue:   qaItem.QaValue,
					CallOwner: qaItem.CallOwner,
					Status:    qaItem.Status,
				}
				qaList = append(qaList, newQaItem)
			}
		}
	}

	for _, qaItem := range qaList {
		if len(qaItem.ContactId) < 1 {
			continue
		}
		if len(qaConf[qaItem.ContactId]) < 1 {
			confItem := make([]entity.SkillQaEntity, 0)
			qaConf[qaItem.ContactId] = confItem
		}
		qaConf[qaItem.ContactId] = append(qaConf[qaItem.ContactId], qaItem)
	}

	log.Printf("qa conf is %#v", qaConf)
	q.QaConf = qaConf
}
