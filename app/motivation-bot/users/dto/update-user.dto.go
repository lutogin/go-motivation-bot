package usersDto

type UpdateUserDto struct {
	Id        string `json:"id" bson:"-" validate:"required,mongodb"`
	ChatId    int64  `json:"chatId" bson:"chatId" validate:"omitempty"`
	FirstName string `json:"firstName" bson:"firstName,omitempty" validate:"omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty" validate:"omitempty"`
	UserName  string `json:"userName" bson:"userName" validate:"omitempty"`
	Cron      string `json:"cron" bson:"cron" validate:"omitempty"`
}
