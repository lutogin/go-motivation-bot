package usersDto

type DeleteUserByUserNameDto struct {
	UserName string `json:"userName" bson:"userName" validate:"required"`
}
