package users

import (
	"motivation-bot/users/dto"
)

type UserEntity struct {
	ChatId       int64                 `json:"chatId" bson:"chatId"`
	FirstName    string                `json:"firstName" bson:"firstName"`
	LastName     string                `json:"lastName" bson:"lastName"`
	UserName     string                `json:"userName" bson:"userName"`
	AlertingTime usersDto.AlertingTime `json:"alertingTime" bson:"alertingTime"`
	Lang         string                `json:"lang" bson:"lang"`
}
