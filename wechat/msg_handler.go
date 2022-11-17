package wechat

import (
	"log"
)

var (
	msgHandler = &MsgHandler{}
)

func GetMsgHandler() *MsgHandler {
	if msgHandler.handlerMap == nil {
		msgHandler.handlerMap = make(map[MsgName]Handler)
	}
	return msgHandler
}

type MsgName string
type Handler func(...interface{})

type MsgHandler struct {
	handlerMap map[MsgName]Handler // 只实现1对1
}

func (m *MsgHandler) On(msgName MsgName, handler Handler) {
	m.handlerMap[msgName] = handler
}

func (m *MsgHandler) Handle(msgName MsgName, data interface{}) {
	if m.handlerMap == nil {
		log.Print("msghandler map is empty.")
		return
	}

	handler, ok := m.handlerMap[msgName]
	if ok {
		handler(data) // 这里实际执行了处理过程
	} else {
		log.Print("msg has no handler,name is %s.", msgName)
		return
	}
}
