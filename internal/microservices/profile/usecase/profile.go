package usecase

import (
	"context"
	"fmt"
	"main/internal/constants"
	"main/internal/microservices/profile"
	proto "main/internal/microservices/profile/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	storage profile.Storage
}

func NewService(storage profile.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetUserProfile(ctx context.Context, userID *proto.UserID) (*proto.ProfileData, error) {
	userData, err := s.storage.GetUserProfile(userID.ID)
	if err != nil {
		return &proto.ProfileData{}, status.Error(codes.Internal, err.Error())
	}
	fmt.Println(userData)

	return userData, nil
}

func (s *Service) EditProfile(ctx context.Context, data *proto.EditProfileData) (*proto.Empty, error) {
	err := s.storage.EditProfile(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) EditAvatar(ctx context.Context, data *proto.EditAvatarData) (*proto.Empty, error) {
	oldAvatar, err := s.storage.EditAvatar(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if oldAvatar != constants.DefaultImage {
		err = s.storage.DeleteFile(oldAvatar)
		if err != nil {
			return &proto.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}

	return &proto.Empty{}, nil
}

func (s *Service) UploadAvatar(ctx context.Context, data *proto.UploadInputFile) (*proto.FileName, error) {
	name, err := s.storage.UploadAvatar(data)
	if err != nil {
		return &proto.FileName{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.FileName{Name: name}, nil
}

func (s *Service) GetAvatar(ctx context.Context, userID *proto.UserID) (*proto.FileName, error) {
	name, err := s.storage.GetAvatar(userID.ID)
	if err != nil {
		return &proto.FileName{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.FileName{Name: name}, nil
}

func (s *Service) AcceptInvitationToFamily(ctx context.Context, data *proto.AddToFamily) error {
	err := s.storage.AcceptInvitationToFamily(data)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *Service) CreateFamily(ctx context.Context, userID *proto.UserID) error {
	hasFamily, _, _, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if hasFamily == true {
		return status.Error(codes.AlreadyExists, constants.ErrFamilyAlreadyExists.Error())
	}

	err = s.storage.CreateFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *Service) DeleteFamily(ctx context.Context, userID *proto.UserID) error {
	hasFamily, idMainUser, _, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != userID.ID {
		return status.Error(codes.Internal, constants.ErrNotMainUser.Error())
	}

	err = s.storage.DeleteFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *Service) ExitFromFamily(ctx context.Context, userID *proto.UserID) error {
	hasFamily, idMainUser, _, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser == userID.ID {
		return status.Error(codes.Internal, constants.ErrMainUser.Error())
	}

	err = s.storage.ExitFromFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *Service) DeleteFromFamily(ctx context.Context, userID *proto.UserID, userIDToDelete *proto.UserID) error {
	hasFamily, idMainUser, _, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != userID.ID || idMainUser == userIDToDelete.ID {
		return status.Error(codes.Internal, constants.ErrNotAvailableForDelete.Error())
	}

	err = s.storage.DeleteFromFamily(userIDToDelete.ID)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *Service) AddMember(ctx context.Context, data *proto.MemberData) error {
	hasFamily, idMainUser, idFamily, err := s.storage.HasFamily(data.IDMainUser)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != data.IDMainUser {
		return status.Error(codes.Internal, constants.ErrNotAvailableForAdd.Error())
	}

	data.IDFamily = idFamily

	err = s.storage.AddMember(data)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s *Service) GetFamily(ctx context.Context, userID *proto.UserID) (*proto.ResponseMemberDataArr, error) {
	members, err := s.storage.GetFamily(userID.ID)
	if err != nil {
		return &proto.ResponseMemberDataArr{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.ResponseMemberDataArr{ResponseMemberData: members}, nil
}

func (s *Service) HasFamily(ctx context.Context, userID *proto.UserID) (*proto.HasFamilyResp, error) {
	has, user, family, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return &proto.HasFamilyResp{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.HasFamilyResp{
		Has:        has,
		IDMainUser: user,
		IDFamily:   family,
	}, nil
}
