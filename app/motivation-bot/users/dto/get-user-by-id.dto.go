package usersDto

type GetUserByIdDto struct {
	Id string `json:"id" bson:"_id,omitempty" validate:"required,mongodb"` // omitempty means can be null
}
