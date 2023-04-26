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
	DeleteFile(string, string) error

	HasFamily(userID int64) (bool, int64, int64, bool, error)
	CreateFamily(userID int64) error
	DeleteFamily(userID int64) error

	AcceptInvitationToFamily(data *proto.AddToFamily) error
	DeleteFromFamily(userID int64) error
	DeleteMember(userID int64) (string, error)
	AddMember(data *proto.MemberData) error
	GetFamily(userID int64) ([]*proto.ResponseMemberData, error)

	IsUserExists(data *proto.EmailData) (bool, error)

	AddMedicine(data *proto.AddMed) error
	DeleteMedicine(data *proto.DeleteMed) (string, error)
	GetMedicine(userID int64) ([]*proto.GetMedicineData, error)
	GetMedicineFamily(familyID int64) ([]*proto.GetMedicineData, error)
	EditMedicine(data *proto.GetMedicineData) (string, error)

	AddNotification(data *proto.NotificationData) error
	DeleteNotification(data *proto.DeleteNotificationData) error
	GetNotifications(userID int64) ([]*proto.GetNotificationData, error)
	GetNotificationsFamily(familyID int64) ([]*proto.GetNotificationData, error)
	AcceptNotification(data *proto.Accept) (int64, error)
	Substruct(idMedicine, count int64) error
}
