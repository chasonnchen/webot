package entity

type SkillForwardEntity struct {
	Id   int32
	Name string
	// Spekers 为空时，转发所有人消息，不为空时判断发言人
	// Spekers 支持按英文逗号分割，配置多个，任何一人发的都会触发
	Spekers string
	// Keywords 为空时，转发所有消息，不为空时判断包含关键字时触发
	// Keywords 支持按英文逗号分割，配置多个，命中任何一个就会触发
	Keywords string

	// 屏蔽词做为约束条件，只要包含之一就不转发
	BadKeywords string
	FromGroupId int32
	ToGroupId   int32
	BotId       string
	Status      int32
}

func (SkillForwardEntity) TableName() string {
	return "t_skill_forward"
}
