package users

import "time"

type UserEntity struct {
	ChatId       int64     `json:"chatId" bson:"chatId"`
	FirstName    string    `json:"firstName" bson:"firstName"`
	LastName     string    `json:"lastName" bson:"lastName"`
	UserName     string    `json:"userName" bson:"userName"`
	AlertingDate time.Time `json:"alertingDate" bson:"alertingDate"`
	Lang         string    `json:"lang" bson:"lang"`
}
