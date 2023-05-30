package usersDto

type DeleteUserByChatIdDto struct {
	ChatId int64 `json:"chatId" bson:"chatId" validate:"required"`
}
