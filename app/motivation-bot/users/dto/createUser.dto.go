package usersDto

import "time"

type CreateUserDto struct {
	ChatId       int64     `json:"chatId" bson:"chatId" validate:"required"`
	FirstName    string    `json:"firstName" bson:"firstName"`
	LastName     string    `json:"lastName" bson:"lastName"`
	UserName     string    `json:"userName" bson:"userName" validate:"required"`
	AlertingDate time.Time `json:"alertingDate" bson:"alertingDate"`
	Lang         string    `json:"lang" bson:"lang"`
}
