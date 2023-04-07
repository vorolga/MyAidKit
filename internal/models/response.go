package models

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseUserProfile struct {
	Status   int             `json:"status"`
	UserData *ProfileUserDTO `json:"user"`
}

type ResponseMembers struct {
	Status  int      `json:"status"`
	Members []Member `json:"members"`
}
