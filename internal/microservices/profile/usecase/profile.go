package usecase

import (
	"context"
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

func (s *Service) AcceptInvitationToFamily(ctx context.Context, data *proto.AddToFamily) (*proto.Empty, error) {
	err := s.storage.AcceptInvitationToFamily(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) CreateFamily(ctx context.Context, userID *proto.UserID) (*proto.Empty, error) {
	hasFamily, _, _, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if hasFamily == true {
		return &proto.Empty{}, status.Error(codes.AlreadyExists, constants.ErrFamilyAlreadyExists.Error())
	}

	err = s.storage.CreateFamily(userID.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) DeleteFamily(ctx context.Context, userID *proto.UserID) (*proto.Empty, error) {
	hasFamily, idMainUser, _, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != userID.ID {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNotMainUser.Error())
	}

	err = s.storage.DeleteFamily(userID.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) DeleteFromFamily(ctx context.Context, Delete *proto.Delete) (*proto.Empty, error) {
	hasFamily, idMainUser, _, err := s.storage.HasFamily(Delete.UserID.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != Delete.UserID.ID || idMainUser == Delete.UserToDelete.ID {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNotAvailableForDelete.Error())
	}

	err = s.storage.DeleteFromFamily(Delete.UserToDelete.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) AddMember(ctx context.Context, data *proto.MemberData) (*proto.Empty, error) {
	hasFamily, idMainUser, idFamily, err := s.storage.HasFamily(data.IDMainUser)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != data.IDMainUser {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNotAvailableForAdd.Error())
	}

	data.IDFamily = idFamily

	err = s.storage.AddMember(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Empty{}, nil
}

func (s *Service) GetFamily(ctx context.Context, userID *proto.UserID) (*proto.ResponseMemberDataArr, error) {
	members, err := s.storage.GetFamily(userID.ID)
	if err != nil {
		return &proto.ResponseMemberDataArr{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.ResponseMemberDataArr{ResponseMemberData: members}, nil
}

func (s *Service) DeleteMember(ctx context.Context, Delete *proto.Delete) (*proto.Empty, error) {
	hasFamily, idMainUser, _, err := s.storage.HasFamily(Delete.UserID.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if hasFamily == false {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNoFamily.Error())
	}

	if idMainUser != Delete.UserID.ID {
		return &proto.Empty{}, status.Error(codes.Internal, constants.ErrNotAvailableForDelete.Error())
	}

	avatar, err := s.storage.DeleteMember(Delete.UserToDelete.ID)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}

	if avatar != constants.DefaultImage {
		err = s.storage.DeleteFile(avatar)
		if err != nil {
			return &proto.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}

	return &proto.Empty{}, nil
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

func (s *Service) UserExists(ctx context.Context, data *proto.EmailData) (*proto.Exists, error) {
	exists, err := s.storage.IsUserExists(data)
	if err != nil {
		return &proto.Exists{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Exists{
		Exists: exists,
	}, nil
}

func (s *Service) AddMedicine(ctx context.Context, data *proto.AddMed) (*proto.Empty, error) {
	err := s.storage.AddMedicine(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (s *Service) DeleteMedicine(ctx context.Context, data *proto.DeleteMed) (*proto.Empty, error) {
	err := s.storage.DeleteMedicine(data)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (s *Service) GetMedicine(ctx context.Context, userID *proto.UserID) (*proto.MedicineArr, error) {
	has, _, family, err := s.storage.HasFamily(userID.ID)
	if err != nil {
		return &proto.MedicineArr{}, status.Error(codes.Internal, err.Error())
	}

	if !has {
		medicines, err := s.storage.GetMedicine(userID.ID)
		if err != nil {
			return &proto.MedicineArr{}, status.Error(codes.Internal, err.Error())
		}
		return &proto.MedicineArr{MedicineArr: medicines}, nil
	}

	medicines, err := s.storage.GetMedicineFamily(family)
	if err != nil {
		return &proto.MedicineArr{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.MedicineArr{MedicineArr: medicines}, nil
}
