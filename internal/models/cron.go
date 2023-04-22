package models

type NotificationsFrom struct {
	IDFrom       int64
	ToIsUser     bool
	IDTo         int64
	NameTo       string
	NameMedicine string
	Email        string
	IDFamily     int64
}
