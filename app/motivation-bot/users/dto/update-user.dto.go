package usersDto

type UpdateUserDto struct {
	Id      string `json:"id" bson:"-" validate:"required,mongodb"`
	Name    string `json:"name" bson:"name,omitempty"`
	Email   string `json:"email" bson:"email,omitempty" validate:"omitempty,email"`
	Phone   string `json:"phone" bson:"phone,omitempty" validate:"omitempty,e164"`
	PwdHash string `json:"-" bson:"pwdHash,omitempty"`
}
