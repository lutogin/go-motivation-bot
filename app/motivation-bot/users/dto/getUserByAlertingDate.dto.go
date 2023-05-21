package usersDto

import "time"

type GetUserByAlertingDateDto struct {
	From time.Time `json:"from" bson:"from" validate:"required"`
	To   time.Time `json:"to" bson:"to" validate:"required"`
}
