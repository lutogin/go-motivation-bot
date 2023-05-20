package users

type UserEntity struct {
	ChatId       string `json:"chatId" bson:"chatID,omitempty"` // omitempty means can be null
	FirstName    string `json:"firstName" bson:"firstName"`
	LastName     string `json:"lastName" bson:"lastName"`
	UserName     string `json:"userName" bson:"userName"`
	AlertingDate string `json:"alertingDate" bson:"alertingDate"` // json:"-" means don't be get in JSON
}
