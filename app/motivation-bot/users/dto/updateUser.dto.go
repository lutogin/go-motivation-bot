package usersDto

import "time"

type UpdateUserDto struct {
	ChatId       int64     `json:"chatId,omitempty" bson:"chatId,omitempty"`
	FirstName    string    `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName     string    `json:"lastName,omitempty" bson:"lastName,omitempty"`
	UserName     string    `json:"userName,omitempty" bson:"userName,omitempty"`
	AlertingDate time.Time `json:"alertingDate,omitempty" bson:"alertingDate,omitempty"`
	Lang         string    `json:"lang,omitempty" bson:"lang,omitempty"`
}
