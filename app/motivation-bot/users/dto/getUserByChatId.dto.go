package usersDto

type GetUserByChatIdDto struct {
	ChatId int64 `json:"chatId" bson:"chatId,omitempty" validate:"required"` // omitempty means can be null
}
