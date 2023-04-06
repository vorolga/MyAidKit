package models

type CreateUserDTO struct {
	Name     string `json:"name" form:"name"`
	Surname  string `json:"surname" form:"surname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Date     string `json:"date" form:"date"`
}

type LogInUserDTO struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
