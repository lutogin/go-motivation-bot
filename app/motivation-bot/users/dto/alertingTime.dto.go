package usersDto

type AlertingTime struct {
	Hours   int `json:"hours" bson:"hours"`
	Minutes int `json:"minutes" bson:"minutes"`
}
