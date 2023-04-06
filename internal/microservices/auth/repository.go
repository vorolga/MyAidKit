package auth

import proto "main/internal/microservices/auth/proto"

type Storage interface {
	IsUserExists(data *proto.LogInData) (int64, error)
	IsUserUnique(email string) (bool, error)
	CreateUser(data *proto.SignUpData) (*proto.Hash, error)
	ConfirmEmail(data *proto.Hash) error
	GetEmailLink(domen string) (*proto.EmailLink, error)

	StoreSession(userID int64) (string, error)
	GetUserID(session string) (int64, error)
	DeleteSession(session string) error
}
