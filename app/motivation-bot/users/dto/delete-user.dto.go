package usersDto

type DeleteUserDto struct {
	Id string `json:"id" bson:"_id" validate:"required,mongodb"` // omitempty means can be null
}
