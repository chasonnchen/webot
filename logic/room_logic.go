package logic

import (
	"encoding/xml"
	"log"
	"strings"

	"github.com/chasonnchen/webot/service"
	"github.com/chasonnchen/webot/wechat"
)

var (
	roomLogicInstance = &RoomLogic{}
)

func NewRoomLogic() *RoomLogic {
	return roomLogicInstance
}

type RoomLogic struct {
}

type RoomJoinXml struct {
	XMLName        xml.Name       `xml:"sysmsg"`
	Sysmsgtemplate Sysmsgtemplate `xml:"sysmsgtemplate"`
}
type Sysmsgtemplate struct {
	XMLName         xml.Name        `xml:"sysmsgtemplate"`
	Contenttemplate Contenttemplate `xml:"content_template"`
}
type Contenttemplate struct {
	XMLName  xml.Name `xml:"content_template"`
	Linklist Linklist `xml:"link_list"`
}
type Linklist struct {
	XMLName xml.Name `xml:"link_list"`
	Link    []Link   `xml:"link"`
}
type Link struct {
	XMLName    xml.Name   `xml:"link"`
	Name       string     `xml:"name,attr"`
	Memberlist Memberlist `xml:"memberlist"`
}
type Memberlist struct {
	XMLName xml.Name `xml:"memberlist"`
	Member  []Member `xml:"member"`
}
type Member struct {
	XMLName  xml.Name `xml:"member"`
	Username Cdata    `xml:"username"`
	Nickname Cdata    `xml:"nickname"`
}
type Cdata struct {
	Value string `xml:",cdata"`
}

func (r *RoomLogic) DoJoin(msg *wechat.Msg) {
	contact := service.NewContactService().GetById(msg.Data.FromGroup)
	// 尝试解析一下XML
	xmlMap := RoomJoinXml{}
	err := xml.Unmarshal([]byte(msg.Data.Content), &xmlMap)
	if err != nil {
		log.Printf("xml to map fail. err is %+v", err)
	}
	log.Printf("xml map is %+v", xmlMap)

	var mlist []Member
	for _, link := range xmlMap.Sysmsgtemplate.Contenttemplate.Linklist.Link {
		if link.Name == "names" {
			mlist = link.Memberlist.Member
		}
	}
	atStr := ""
	idlist := make([]string, 0)
	for _, m := range mlist {
		atStr += "@" + m.Nickname.Value + string(rune(8197))
		idlist = append(idlist, m.Username.Value)
	}

	if len(contact.Hello) > 0 {
		service.NewContactService().SayTextToContact(contact.Id, atStr+contact.Hello, strings.Join(idlist, ","))
	}
}
