package entity

type UserEntity struct {
	Id            int32
	LoginName     string
	LoginPassword string
	Name          string
	ContactId     string
	Tel           string
	AppId         int32
	AppKey        string
	Status        int32
}

func (UserEntity) TableName() string {
	return "t_user"
}
