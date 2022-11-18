package logic

import (
	"log"

	"github.com/chasonnchen/webot/entity"
	"github.com/chasonnchen/webot/lib/baidu/unit"
	"github.com/chasonnchen/webot/service"
	"github.com/chasonnchen/webot/wechat"
)

const (
	contactTypeRoom = 2
	contactTypeUser = 1
)

var (
	messageLogicInstance = &MessageLogic{}
)

func NewMessageLogic() *MessageLogic {
	return messageLogicInstance
}

type MessageLogic struct {
}

func (m *MessageLogic) buildContact(message *wechat.Msg) entity.ContactEntity {
	contact := entity.ContactEntity{}

	if message.MessageType == "80001" {
		contact.Id = message.Data.FromGroup
		contact.Type = 2
		contact.Status = 1
	} else {
		contact.Id = message.Data.FromUser
		contact.Type = 1
		contact.Status = 1
	}

	return contact
}

func (m *MessageLogic) buildMsgText(message *wechat.Msg) string {
	msgText := "[" + message.Data.FromUser + "]"

	msgText = msgText + ": " + message.Data.Content

	return msgText
}

func (m *MessageLogic) Do(message *wechat.Msg) {
	// 0. log
	log.Printf("MessageLogic recive message: %s", m.buildMsgText(message))
	contact := m.buildContact(message)
	// 1. 更新联系人
	contact = service.NewContactService().Upsert(contact)

	// 2. 问答
	service.NewQaService().DoQa(contact, message)

	// 3. 转发
	// service.NewForwardService().DoForward(contact, message)
	// service.NewForwardMediaService().DoForward(contact, message)

	// 4. 暗号加群
	// service.NewRoomService().AutoInvite(message.From(), message, "")

	// 5. 智能聊天
	// 5.1 好友聊天 && 打开智能聊天配置 && 文本消息
	if contact.Type == 1 && contact.OpenAi == 1 {
		log.Print("start ai\n")
		aiRes, _ := unit.NewClient().Chat(contact.Id, message.Data.Content)
		if len(aiRes) > 1 {
			service.NewContactService().SayTextToContact(contact.Id, aiRes, "")
		}
	}
	// 5.2 群里聊天 && @机器人 && 文本消息 && 非群公告
	// @发言人 并回复智能聊天的结果
	/*selfAliasName, _ := message.Room().Alias(service.NewGlobleService().GetBot().UserSelf())
	if len(selfAliasName) < 1 {
		selfAliasName = service.NewGlobleService().GetBot().UserSelf().Name()
	}
	log.Printf("self room alias name is %s", selfAliasName)
	if contact.Type == 2 && message.MentionSelf() && message.Type() == schemas.MessageTypeText && strings.Contains(message.Text(), selfAliasName) {
		aiRes, _ := unit.NewClient().Chat(message.From().ID(), message.Text()[strings.Index(message.Text(), string(rune(8197)))+3:])
		if len(aiRes) > 1 {
			message.Room().Say(aiRes, message.From())
		}
	}*/

	// 6. 内容上传
	// service.NewUploadService().DoUpload(contact, message)
}
