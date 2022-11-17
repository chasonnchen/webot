package entity

type SkillInviteEntity struct {
	Id        int32
	Name      string
	Keyword   string
	ContactId string
	Hello     string
	BotId     string
	Status    int32
}

func (SkillInviteEntity) TableName() string {
	return "t_skill_invite"
}
