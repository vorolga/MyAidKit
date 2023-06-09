package repository

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"main/internal/constants"
	"main/internal/microservices/auth/utils/hash"
	"main/internal/microservices/profile"
	proto "main/internal/microservices/profile/proto"
	"main/internal/microservices/profile/utils/images"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gomodule/redigo/redis"
	"github.com/minio/minio-go/v7"
)

type Storage struct {
	db    *sql.DB
	minio *minio.Client
	redis *redis.Pool
}

func NewStorage(db *sql.DB, minio *minio.Client, redis *redis.Pool) profile.Storage {
	return &Storage{db: db, minio: minio, redis: redis}
}

func (s Storage) GetUserProfile(userID int64) (*proto.ProfileData, error) {
	sqlScript := "SELECT name, surname, email, avatar, birthday, is_adult FROM users WHERE id=$1"

	var name, surname, email, avatar, birthday string
	var isAdult bool
	err := s.db.QueryRow(sqlScript, userID).Scan(&name, &surname, &email, &avatar, &birthday, &isAdult)

	if err != nil {
		return nil, err
	}

	avatarUrl, err := images.GenerateFileURL(avatar, constants.UserObjectsBucketName)
	if err != nil {
		return nil, err
	}

	has, user, _, _, err := s.HasFamily(userID)
	if err != nil {
		return nil, err
	}

	main := false

	if has == true {
		if user == userID {
			main = true
		}
	}

	return &proto.ProfileData{
		Name:    name,
		Surname: surname,
		Email:   email,
		Avatar:  avatarUrl,
		Date:    birthday[:10],
		Main:    main,
		Adult:   isAdult,
	}, nil
}

func (s Storage) EditProfile(data *proto.EditProfileData) error {
	sqlScript := "SELECT name, surname, password, salt, birthday FROM users WHERE id=$1"

	var oldName, oldSurname, oldPassword, oldSalt, oldBirthday string
	err := s.db.QueryRow(sqlScript, data.ID).Scan(&oldName, &oldSurname, &oldPassword, &oldSalt, &oldBirthday)
	if err != nil {
		return err
	}

	notChangedPassword, _ := hash.ComparePasswords(oldPassword, oldSalt, data.Password)

	if !notChangedPassword && len(data.Password) != 0 {
		salt, err := uuid.NewV4()
		if err != nil {
			return err
		}

		hashPassword, err := hash.HashAndSalt(data.Password, salt.String())
		if err != nil {
			return err
		}

		oldPassword = hashPassword
		oldSalt = salt.String()
	}

	if data.Name != oldName && len(data.Name) != 0 {
		oldName = data.Name
	}

	if data.Surname != oldSurname && len(data.Surname) != 0 {
		oldSurname = data.Surname
	}

	if data.Date != oldBirthday && data.Date != "" {
		oldBirthday = data.Date
	}

	sqlScript = "UPDATE users SET name = $2, surname = $3, password = $4, salt = $5, birthday = TO_TIMESTAMP($6, 'YYYY-MM-DD') WHERE id = $1"

	_, err = s.db.Exec(sqlScript, data.ID, oldName, oldSurname, oldPassword, oldSalt, oldBirthday)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) EditAvatar(data *proto.EditAvatarData) (string, error) {
	sqlScript := "SELECT avatar FROM users WHERE id=$1"

	var oldAvatar string
	err := s.db.QueryRow(sqlScript, data.ID).Scan(&oldAvatar)
	if err != nil {
		return "", err
	}

	if len(data.Avatar) != 0 {
		sqlScript := "UPDATE users SET avatar = $2 WHERE id = $1"

		_, err = s.db.Exec(sqlScript, data.ID, data.Avatar)
		if err != nil {
			return "", err
		}

		return oldAvatar, nil
	}

	return "", nil
}

func (s Storage) GetAvatar(userID int64) (string, error) {
	sqlScript := "SELECT avatar FROM users WHERE id=$1"

	var avatar string
	err := s.db.QueryRow(sqlScript, userID).Scan(&avatar)

	if err != nil {
		return "", err
	}

	return avatar, nil
}

func (s Storage) UploadAvatar(data *proto.UploadInputFile) (string, error) {
	imageName := images.GenerateObjectName(data.ID)

	opts := minio.PutObjectOptions{
		ContentType:  data.ContentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	}

	_, err := s.minio.PutObject(
		context.Background(),
		data.BucketName, // Константа с именем бакета
		imageName,
		bytes.NewReader(data.File),
		data.Size,
		opts,
	)
	if err != nil {
		return "", err
	}

	return imageName, nil
}

func (s Storage) DeleteFile(name string, bucket string) error {
	opts := minio.RemoveObjectOptions{}

	err := s.minio.RemoveObject(
		context.Background(),
		bucket,
		name,
		opts,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) AcceptInvitationToFamily(data *proto.AddToFamily) error {
	sqlScript := "UPDATE users SET id_family = $2, is_adult = $3 WHERE email = $1"

	_, err := s.db.Exec(sqlScript, data.Email, data.ID, data.IsAdult)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) CreateFamily(userID int64) error {
	var familyID int64

	sqlScript := "INSERT INTO family(id_main_user) VALUES($1) RETURNING id"

	if err := s.db.QueryRow(sqlScript, userID).Scan(&familyID); err != nil {
		return err
	}

	sqlScript = "UPDATE users SET id_family = $2 WHERE id = $1"
	_, err := s.db.Exec(sqlScript, userID, familyID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) HasFamily(userID int64) (bool, int64, int64, bool, error) {
	sqlScript := "SELECT id_family, is_adult FROM users WHERE id=$1"

	var idFamily int64
	var isAdult bool
	err := s.db.QueryRow(sqlScript, userID).Scan(&idFamily, &isAdult)
	if err != nil {
		return false, 0, 0, false, err
	}

	if idFamily == 0 {
		return false, 0, 0, isAdult, nil
	}

	sqlScript = "SELECT id_main_user FROM family WHERE id=$1"

	var idMainUser int64
	err = s.db.QueryRow(sqlScript, idFamily).Scan(&idMainUser)
	if err != nil {
		return false, 0, 0, false, err
	}

	fmt.Println(isAdult)
	return true, idMainUser, idFamily, isAdult, nil
}

func (s Storage) DeleteFamily(userID int64) error {
	sqlScript := "SELECT id_family FROM users WHERE id=$1"

	var idFamily int64
	err := s.db.QueryRow(sqlScript, userID).Scan(&idFamily)

	if err != nil {
		return err
	}

	sqlScript = "UPDATE users SET id_family = 0 WHERE id_family = $1"
	_, err = s.db.Exec(sqlScript, idFamily)
	if err != nil {
		return err
	}

	sqlScript = "DELETE FROM family WHERE id_main_user = $1"
	_, err = s.db.Exec(sqlScript, userID)
	if err != nil {
		return err
	}

	sqlScript = "DELETE FROM notification_user WHERE to_is_user = false and id_to_user IN (SELECT id FROM members WHERE id_family = $1)"
	_, err = s.db.Exec(sqlScript, idFamily)
	if err != nil {
		return err
	}

	sqlScript = "DELETE FROM members WHERE id_family = $1"
	_, err = s.db.Exec(sqlScript, idFamily)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) DeleteFromFamily(userID int64) error {
	sqlScript := "UPDATE users SET id_family = 0 WHERE id = $1"
	_, err := s.db.Exec(sqlScript, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) AddMember(data *proto.MemberData) error {
	sqlScript := "INSERT INTO members(id_main_user, id_family, name, avatar) VALUES($1, $2, $3, $4)"

	if _, err := s.db.Exec(sqlScript, data.IDMainUser, data.IDFamily, data.Name, data.Avatar); err != nil {
		return err
	}

	return nil
}

func (s Storage) GetFamily(userID int64) ([]*proto.ResponseMemberData, error) {
	sqlScript := "SELECT id_family FROM users WHERE id=$1"

	var idFamily int64
	err := s.db.QueryRow(sqlScript, userID).Scan(&idFamily)
	if err != nil {
		return nil, err
	}

	if idFamily == 0 {
		return nil, nil
	}

	members := make([]*proto.ResponseMemberData, 0)
	sqlScript = "SELECT id, name, avatar, is_adult FROM users WHERE id_family = $1"

	rows, err := s.db.Query(sqlScript, idFamily)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member proto.ResponseMemberData
		member.IsUser = true
		if err = rows.Scan(&member.ID, &member.Name, &member.Avatar, &member.IsAdult); err != nil {
			return nil, err
		}
		member.Avatar, err = images.GenerateFileURL(member.Avatar, constants.UserObjectsBucketName)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	sqlScript = "SELECT id, name, avatar FROM members WHERE id_family = $1"

	rows, err = s.db.Query(sqlScript, idFamily)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member proto.ResponseMemberData
		member.IsAdult = false
		member.IsUser = false
		if err = rows.Scan(&member.ID, &member.Name, &member.Avatar); err != nil {
			return nil, err
		}
		member.Avatar, err = images.GenerateFileURL(member.Avatar, constants.UserObjectsBucketName)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	return members, nil
}

func (s Storage) DeleteMember(userID int64) (string, error) {
	sqlScript := "SELECT avatar FROM members WHERE id=$1"

	var avatar string
	err := s.db.QueryRow(sqlScript, userID).Scan(&avatar)
	if err != nil {
		return "", err
	}

	sqlScript = "DELETE FROM notification_user WHERE to_is_user = false AND id_to_user=$1"
	_, err = s.db.Exec(sqlScript, userID)
	if err != nil {
		return "", err
	}

	sqlScript = "DELETE FROM members WHERE id = $1"
	_, err = s.db.Exec(sqlScript, userID)
	if err != nil {
		return "", err
	}

	return avatar, nil
}

func (s Storage) IsUserExists(data *proto.EmailData) (bool, error) {
	sqlScript := "SELECT id FROM users WHERE email=$1"
	rows, err := s.db.Query(sqlScript, data.Email)
	if err != nil {
		return false, err
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	// убедимся, что всё закроется при выходе из программы
	defer func() {
		rows.Close()
	}()

	// Из базы пришел пустой запрос, значит пользователя в базе данных нет
	if !rows.Next() {
		return false, constants.ErrWrongData
	}

	return true, nil
}

func (s Storage) AddMedicine(data *proto.AddMed) error {
	sqlScript := "INSERT INTO medicine(id_user, name, count, image, is_tablets) VALUES($1, $2, $3, $4, $5)"

	if _, err := s.db.Exec(sqlScript, data.UserID, data.Medicine.Name, data.Medicine.Count, data.Medicine.Image, data.Medicine.IsTablets); err != nil {
		return err
	}

	return nil
}

func (s Storage) DeleteMedicine(data *proto.DeleteMed) (string, error) {
	sqlScript := "SELECT image FROM medicine WHERE id=$1"

	var image string
	err := s.db.QueryRow(sqlScript, data.MedicineID).Scan(&image)
	if err != nil {
		return "", err
	}
	sqlScript = "DELETE FROM medicine WHERE id=$1"

	_, err = s.db.Exec(sqlScript, data.MedicineID)
	if err != nil {
		return "", err
	}
	return image, nil
}

func (s Storage) GetMedicine(userID int64) ([]*proto.GetMedicineData, error) {
	medicines := make([]*proto.GetMedicineData, 0)
	sqlScript := "SELECT id, name, count, image, is_tablets FROM medicine WHERE id_user = $1"

	rows, err := s.db.Query(sqlScript, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var medicine proto.GetMedicineData
		medicine.Medicine = &proto.Medicine{
			Image:     "",
			Name:      "",
			IsTablets: false,
			Count:     0,
		}
		if err = rows.Scan(&medicine.ID, &medicine.Medicine.Name, &medicine.Medicine.Count, &medicine.Medicine.Image, &medicine.Medicine.IsTablets); err != nil {
			return nil, err
		}
		medicine.Medicine.Image, err = images.GenerateFileURL(medicine.Medicine.Image, constants.MedicinesObjectsBucketName)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, &medicine)
	}

	return medicines, nil
}

func (s Storage) GetMedicineFamily(familyID int64) ([]*proto.GetMedicineData, error) {
	sqlScript := "SELECT medicine.id, medicine.name, medicine.count, medicine.image, medicine.is_tablets " +
		"FROM medicine JOIN users u ON u.id_family = $1 AND medicine.id_user = u.id"

	rows, err := s.db.Query(sqlScript, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medicines := make([]*proto.GetMedicineData, 0)
	for rows.Next() {
		var medicine proto.GetMedicineData
		medicine.Medicine = &proto.Medicine{
			Image:     "",
			Name:      "",
			IsTablets: false,
			Count:     0,
		}
		if err = rows.Scan(&medicine.ID, &medicine.Medicine.Name, &medicine.Medicine.Count, &medicine.Medicine.Image, &medicine.Medicine.IsTablets); err != nil {
			return nil, err
		}
		medicine.Medicine.Image, err = images.GenerateFileURL(medicine.Medicine.Image, constants.MedicinesObjectsBucketName)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, &medicine)
	}

	return medicines, nil
}

func (s Storage) EditMedicine(data *proto.GetMedicineData) (string, error) {
	sqlScript := "SELECT name, count, image, is_tablets FROM medicine WHERE id=$1"

	var oldName, oldImage string
	var oldCount int64
	var oldIsTablets bool
	err := s.db.QueryRow(sqlScript, data.ID).Scan(&oldName, &oldCount, &oldImage, &oldIsTablets)
	if err != nil {
		return "", err
	}

	image := oldImage

	if data.Medicine.Name != oldName && len(data.Medicine.Name) != 0 {
		oldName = data.Medicine.Name
	}

	if data.Medicine.Image != oldImage && len(data.Medicine.Image) != 0 && data.Medicine.Image != constants.DefaultMedicine {
		oldImage = data.Medicine.Image
	}

	if data.Medicine.Count != oldCount && data.Medicine.Count != -1 {
		if !oldIsTablets {
			return "", errors.New("cant change count for not tablets")
		}
		oldCount = data.Medicine.Count
	}

	sqlScript = "UPDATE medicine SET name = $2, count = $3, image = $4 WHERE id = $1"

	_, err = s.db.Exec(sqlScript, data.ID, oldName, oldCount, oldImage)
	if err != nil {
		return "", err
	}
	return image, nil
}

func (s Storage) AddNotification(data *proto.NotificationData) error {
	sqlScript := "INSERT INTO notification_user(id_from, to_is_user, id_to_user, name_to, id_medicine, name_medicine, time, is_accepted) VALUES($1, $2, $3, $4, $5, $6, $7, false)"

	if _, err := s.db.Exec(sqlScript, data.IDFrom, data.IsUser, data.IDTo, data.NameTo, data.IDMedicine, data.NameMedicine, data.Time); err != nil {
		return err
	}

	return nil
}

func (s Storage) DeleteNotification(data *proto.DeleteNotificationData) error {
	sqlScript := "DELETE FROM notification_user WHERE id=$1"

	_, err := s.db.Exec(sqlScript, data.NotificationID)
	if err != nil {
		return err
	}
	return nil
}

func (s Storage) GetNotifications(userID int64) ([]*proto.GetNotificationData, error) {
	notifications := make([]*proto.GetNotificationData, 0)
	sqlScript := "SELECT notification_user.id, to_is_user, id_to_user, u.name, id_medicine, medicine.name, medicine.is_tablets, time, is_accepted FROM notification_user JOIN medicine ON medicine.id = id_medicine JOIN users u ON notification_user.id_to_user = u.id AND id_to_user = $1"

	rows, err := s.db.Query(sqlScript, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification proto.GetNotificationData
		notification.NotificationData = &proto.NotificationData{
			IsUser:       false,
			IDTo:         0,
			NameTo:       "",
			IDMedicine:   0,
			NameMedicine: "",
			IsTablets:    false,
			Time:         "",
			IsAccepted:   false,
		}
		if err = rows.Scan(&notification.ID, &notification.NotificationData.IsUser, &notification.NotificationData.IDTo, &notification.NotificationData.NameTo, &notification.NotificationData.IDMedicine, &notification.NotificationData.NameMedicine, &notification.NotificationData.IsTablets, &notification.NotificationData.Time, &notification.NotificationData.IsAccepted); err != nil {
			return nil, err
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (s Storage) GetNotificationsFamily(familyID int64) ([]*proto.GetNotificationData, error) {
	sqlScript := "SELECT notification_user.id, notification_user.to_is_user, notification_user.id_to_user, u.name, notification_user.id_medicine, medicine.name, medicine.is_tablets, notification_user.time, notification_user.is_accepted " +
		"FROM notification_user JOIN medicine ON medicine.id = id_medicine JOIN users u ON u.id_family = $1 AND notification_user.id_to_user = u.id"

	rows, err := s.db.Query(sqlScript, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]*proto.GetNotificationData, 0)
	for rows.Next() {
		var notification proto.GetNotificationData
		notification.NotificationData = &proto.NotificationData{
			IsUser:       false,
			IDTo:         0,
			NameTo:       "",
			IDMedicine:   0,
			NameMedicine: "",
			IsTablets:    false,
			Time:         "",
			IsAccepted:   false,
		}
		if err = rows.Scan(&notification.ID, &notification.NotificationData.IsUser, &notification.NotificationData.IDTo, &notification.NotificationData.NameTo, &notification.NotificationData.IDMedicine, &notification.NotificationData.NameMedicine, &notification.NotificationData.IsTablets, &notification.NotificationData.Time, &notification.NotificationData.IsAccepted); err != nil {
			return nil, err
		}

		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (s Storage) AcceptNotification(data *proto.Accept) (int64, error) {
	sqlScript := "UPDATE notification_user SET is_accepted = true WHERE id = $1 RETURNING id_medicine"

	var medicineID int64
	if err := s.db.QueryRow(sqlScript, data.ID).Scan(&medicineID); err != nil {
		return 0, err
	}

	return medicineID, nil
}

func (s Storage) Substruct(idMedicine, count int64) error {
	sqlScript := "DO $$ DECLARE input_id INT:= " + strconv.Itoa(int(idMedicine)) +
		"; input_count INT:= " + strconv.Itoa(int(count)) + "; BEGIN IF (SELECT count FROM medicine WHERE id = input_id) - input_count >= 0 " +
		"THEN UPDATE medicine SET count = count - input_count WHERE ID = input_id; " +
		"ELSE UPDATE medicine SET count = 0 WHERE ID = input_id; END IF; END $$;"

	_, err := s.db.Exec(sqlScript)
	if err != nil {
		return err
	}
	return nil
}
