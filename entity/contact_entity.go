package entity

type ContactEntity struct {
	Id     string
	Name   string
	Type   int32
	Hello  string
	OpenAi int32
	BotId  string
	Status int32
}

func (ContactEntity) TableName() string {
	return "t_contact"
}
