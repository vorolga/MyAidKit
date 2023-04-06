package models

type ProfileUserDTO struct {
	Name    string `json:"name" form:"name"`
	Surname string `json:"surname" form:"surname"`
	Email   string `json:"email" form:"email"`
	Avatar  string `json:"avatar" form:"avatar"`
	Date    string `json:"date" form:"date"`
	Main    bool   `json:"main" form:"main"`
}

type EditProfileDTO struct {
	Name     string `json:"name" form:"name"`
	Surname  string `json:"surname" form:"surname"`
	Password string `json:"password" form:"password"`
	Date     string `json:"date" form:"date"`
}

type EmailUserDTO struct {
	Email string `json:"email" form:"email"`
}
