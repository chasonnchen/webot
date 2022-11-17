package entity

type BotEntity struct {
	Id      string
	WcId    string
	AuthKey string
	Name    string
	Alias   string
	Head    string
	Uid     int32
	Status  int32
}

func (BotEntity) TableName() string {
	return "t_bot"
}
