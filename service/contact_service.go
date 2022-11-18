package service

import (
	"log"
	"time"

	"github.com/chasonnchen/webot/dao"
	"github.com/chasonnchen/webot/entity"
	"github.com/chasonnchen/webot/wechat"
)

var (
	contactServiceInstance = &ContactService{}
)

func NewContactService() *ContactService {
	return contactServiceInstance
}

type ContactService struct {
	ContactList map[string]entity.ContactEntity
}

func (c *ContactService) SayTextToContact(contactId string, msgText string, at string) {
	bot := wechat.GetWechatInstance()
	bot.WkteamApi.SendText(bot.WId, contactId, msgText, at)
}

func (c *ContactService) GetById(contactId string) entity.ContactEntity {
	return c.ContactList[contactId]
}

func (c *ContactService) Upsert(contact entity.ContactEntity) entity.ContactEntity {
	//  先检查在不在List里
	contactOri, ok := c.ContactList[contact.Id]
	if ok {
		// TODO 更新逻辑
		_, isUpdate := NewCacheService().Get(contactOri.Id + "_update")
		if isUpdate {
		} else {
			contactOri = c.updateContacts([]entity.ContactEntity{contactOri})[0]
			log.Printf("upsert contact to db, %+v", contactOri)
			dao.Wechat().Model(&contactOri).Update("name", contactOri.Name)
			NewCacheService().Set(contactOri.Id+"_update", 1, 2*time.Minute)
		}
	} else {
		// 插入
		contact.BotId = wechat.GetWechatInstance().WcId
		contact = c.updateContacts([]entity.ContactEntity{contact})[0]
		log.Printf("insert contact to db, %+v", contact)
		dao.Wechat().Create(&contact)
		c.ContactList[contact.Id] = contact
		contactOri = contact
	}

	return contactOri
}

func (c *ContactService) updateContacts(contacts []entity.ContactEntity) []entity.ContactEntity {
	contactMapOri := make(map[string]entity.ContactEntity)
	ids := make([]string, 0)
	for _, contact := range contacts {
		ids = append(ids, contact.Id)
		contactMapOri[contact.Id] = contact
	}
	newContactList := make([]entity.ContactEntity, 0)
	for {
		if ids == nil {
			break
		}

		var idList []string
		if len(ids) > 20 {
			idList = ids[:20]
			ids = ids[20:]
		} else {
			idList = make([]string, len(ids))
			copy(idList, ids)
			ids = nil
		}

		// get
		bot := wechat.GetWechatInstance()
		contactMap, _ := bot.WkteamApi.GetContact(bot.WId, idList)

		// append
		for _, contact := range contactMap {
			newContact := contactMapOri[contact.(map[string]interface{})["userName"].(string)]

			if contact.(map[string]interface{})["nickName"] != nil {
				newContact.Name = contact.(map[string]interface{})["nickName"].(string)
			}
			if contact.(map[string]interface{})["remark"] != nil {
				newContact.Name = contact.(map[string]interface{})["remark"].(string)
			}
			newContactList = append(newContactList, newContact)
		}

		// sleet
		time.Sleep(1 * time.Second)
	}
	return newContactList
}

func (c *ContactService) init() {
	c.load()

	go func() {
		for {
			select {
			case <-time.After(time.Second * 300):
				c.load()
			}
		}
	}()
}

func (c *ContactService) load() {
	// 加载数据库中已保存的联系人
	botId := wechat.GetWechatInstance().WcId
	contactMap := make(map[string]entity.ContactEntity)
	var contactList []entity.ContactEntity
	dao.Wechat().Where("bot_id = ?", botId).Find(&contactList)

	for _, contact := range contactList {
		contactMap[contact.Id] = contact
	}

	log.Printf("contact list  is %#v", contactMap)
	c.ContactList = contactMap
	return
}
