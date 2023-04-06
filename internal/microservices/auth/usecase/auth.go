package usecase

import (
	"context"
	"errors"
	"main/internal/constants"
	"main/internal/microservices/auth"
	proto "main/internal/microservices/auth/proto"
	"main/internal/microservices/auth/utils/validation"
	"net/smtp"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	storage auth.Storage
}

func NewService(storage auth.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) SignUp(ctx context.Context, data *proto.SignUpData) (*proto.EmailLink, error) {
	if err := validation.ValidateUser(data); err != nil {
		return &proto.EmailLink{}, status.Error(codes.Internal, err.Error())
	}

	isUnique, err := s.storage.IsUserUnique(data.Email)
	if err != nil {
		return &proto.EmailLink{}, status.Error(codes.Internal, err.Error())
	}

	if !isUnique {
		return &proto.EmailLink{}, status.Error(codes.InvalidArgument, constants.ErrEmailIsNotUnique.Error())
	}

	hash, err := s.storage.CreateUser(data)
	if err != nil {
		return &proto.EmailLink{}, status.Error(codes.Internal, err.Error())
	}

	from := "vorrovvorrov@gmail.com"
	password := os.Getenv("EMAILPASSWORD")

	toList := []string{data.Email}

	host := "smtp.gmail.com"
	port := "587"

	msg := "Подтвердите Email\r\n" +
		"Что бы подтвердить Email, перейдите по ссылке: " +
		"http://" + os.Getenv("HOST") + "/confirm?hash=" + hash.Hash

	body := []byte(msg)

	authSMTP := smtp.PlainAuth("", from, password, host)
	err = smtp.SendMail(host+":"+port, authSMTP, from, toList, body)
	if err != nil {
		return &proto.EmailLink{}, status.Error(codes.Internal, err.Error())
	}

	result := strings.Split(data.Email, "@")
	link, err := s.storage.GetEmailLink(result[1])

	return link, nil
}

func (s *Service) LogIn(ctx context.Context, data *proto.LogInData) (*proto.Cookie, error) {
	userID, err := s.storage.IsUserExists(data)
	if err != nil {
		if errors.Is(err, constants.ErrWrongData) {
			return &proto.Cookie{}, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, constants.ErrEmailIsNotConfirmed) {
			return &proto.Cookie{}, status.Error(codes.Unauthenticated, err.Error())
		}
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}

	session, err := s.storage.StoreSession(userID)
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Cookie{
		Cookie: session,
	}, nil
}

func (s *Service) ConfirmEmail(ctx context.Context, data *proto.Hash) (*proto.Empty, error) {
	err := s.storage.ConfirmEmail(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (s *Service) LogOut(ctx context.Context, cookie *proto.Cookie) (*proto.Empty, error) {
	err := s.storage.DeleteSession(cookie.Cookie)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (s *Service) CheckAuthorization(ctx context.Context, cookie *proto.Cookie) (*proto.UserID, error) {
	userID, err := s.storage.GetUserID(cookie.Cookie)
	if err != nil {
		return &proto.UserID{ID: -1}, status.Error(codes.Internal, err.Error())
	}

	return &proto.UserID{ID: userID}, nil
}
