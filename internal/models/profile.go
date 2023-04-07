package models

type ProfileUserDTO struct {
	ID      int64  `json:"id" form:"id"`
	Name    string `json:"name" form:"name"`
	Surname string `json:"surname" form:"surname"`
	Email   string `json:"email" form:"email"`
	Avatar  string `json:"avatar" form:"avatar"`
	Date    string `json:"date" form:"date"`
	Main    bool   `json:"main" form:"main"`
	Adult   bool   `json:"adult" form:"adult"`
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

type UserIDDTO struct {
	ID int64 `json:"id" form:"id"`
}

type MemberDTO struct {
	Name   string `json:"name" form:"name"`
	Avatar string `json:"avatar" form:"avatar"`
}

type Member struct {
	ID     int64  `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	Avatar string `json:"avatar" form:"avatar"`
	Adult  bool   `json:"adult" form:"adult"`
	User   bool   `json:"user" form:"user"`
}

type InviteUserDTO struct {
	Email string `json:"email" form:"email"`
	Adult bool   `json:"adult" form:"adult"`
}
