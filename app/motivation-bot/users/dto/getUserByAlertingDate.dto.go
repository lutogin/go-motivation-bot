package usersDto

import "time"

type GetUserByAlertingTimeDto struct {
	//From time.Time `json:"from" bson:"from" validate:"required"`
	//To   time.Time `json:"to" bson:"to" validate:"required"`
	Date time.Time `json:"date" bson:"date" validate:"required"`
}
