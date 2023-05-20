package usersDto

type CreateUserDto struct {
	Name  string `json:"name" bson:"name" validate:"required"`
	Email string `json:"email" bson:"email" validate:"required,email"`
	Phone string `json:"phone" bson:"phone" validate:"required,e164"`
	Pwd   string `json:"pwd" bson:"pwd" validate:"required,gte=6"`
}
