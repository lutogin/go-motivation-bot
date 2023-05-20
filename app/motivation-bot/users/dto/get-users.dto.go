package usersDto

type GetUsersDto struct {
	Name  string `json:"name" bson:"name,omitempty" validate:"omitempty"`
	Email string `json:"email" bson:"email,omitempty" validate:"omitempty,email"`
	Phone string `json:"phone" bson:"phone,omitempty"`
}
