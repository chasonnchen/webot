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
		// 更新逻辑
	} else {
		// 插入
		contact.BotId = wechat.GetWechatInstance().WcId
		dao.Wechat().Create(&contact)

		c.updateContacts([]entity.ContactEntity{contact})
		dao.Wechat().Where("id = ?", contact.Id).Take(&contact)
		c.ContactList[contact.Id] = contact
		contactOri = contact
	}

	return contactOri
}

func (c *ContactService) updateContacts(contacts []entity.ContactEntity) {
	ids := make([]string, len(contacts))
	for _, contact := range contacts {
		ids = append(ids, contact.Id)
	}

	newContactList := make([]entity.ContactEntity, len(contacts))

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
			newContact := entity.ContactEntity{
				Id:   contact.(map[string]interface{})["userName"].(string),
				Name: "",
			}
			if contact.(map[string]interface{})["nickName"] != nil {
				newContact.Name = contact.(map[string]interface{})["nickName"].(string)
			}
			if contact.(map[string]interface{})["remark"] != nil {
				newContact.Name = contact.(map[string]interface{})["remark"].(string)
			}
			// 更新下DB
			if len(newContact.Name) > 0 {
				dao.Wechat().Model(&newContact).Update("name", newContact.Name)
			}
			newContactList = append(newContactList, newContact)
		}

		// sleet
		time.Sleep(1 * time.Second)
	}

	return
}

func (c *ContactService) init() {
	// 加载数据库中已保存的联系人
	botId := wechat.GetWechatInstance().WcId
	contactMap := make(map[string]entity.ContactEntity)
	var contactList []entity.ContactEntity
	dao.Wechat().Where("bot_id = ?", botId).Find(&contactList)

	// 更新一下最新信息
	c.updateContacts(contactList)

	// 重新查一次
	dao.Wechat().Where("bot_id = ?", botId).Find(&contactList)

	for _, contact := range contactList {
		contactMap[contact.Id] = contact
	}

	log.Printf("contact list  is %#v", contactMap)
	c.ContactList = contactMap
	return
}
