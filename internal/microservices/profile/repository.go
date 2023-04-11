package profile

import (
	proto "main/internal/microservices/profile/proto"
)

type Storage interface {
	GetUserProfile(userID int64) (*proto.ProfileData, error)
	EditProfile(data *proto.EditProfileData) error
	EditAvatar(data *proto.EditAvatarData) (string, error)
	GetAvatar(userID int64) (string, error)
	UploadAvatar(data *proto.UploadInputFile) (string, error)
	DeleteFile(string) error

	HasFamily(userID int64) (bool, int64, int64, error)
	CreateFamily(userID int64) error
	DeleteFamily(userID int64) error

	AcceptInvitationToFamily(data *proto.AddToFamily) error
	DeleteFromFamily(userID int64) error
	DeleteMember(userID int64) (string, error)
	AddMember(data *proto.MemberData) error
	GetFamily(userID int64) ([]*proto.ResponseMemberData, error)

	IsUserExists(data *proto.EmailData) (bool, error)

	AddMedicine(data *proto.AddMed) error
	DeleteMedicine(data *proto.DeleteMed) error
	GetMedicine(userID int64) ([]*proto.GetMedicine, error)
	GetMedicineFamily(familyID int64) ([]*proto.GetMedicine, error)
}
