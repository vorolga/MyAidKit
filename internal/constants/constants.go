package constants

import (
	"errors"
	"main/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

var (
	ErrWrongData             = errors.New("wrong data")
	ErrLetter                = errors.New("at least one letter is required")
	ErrNum                   = errors.New("at least one digit is required")
	ErrCount                 = errors.New("at least eight characters long is required")
	ErrBan                   = errors.New("password uses unavailable symbols")
	ErrEmailAlreadyConfirmed = errors.New("email is already confirmed")
	ErrEmailIsNotUnique      = errors.New("email is not unique")
	ErrNoEmailLink           = errors.New("no email link")
	ErrEmailIsNotConfirmed   = errors.New("email is not confirmed")

	ErrFamilyAlreadyExists   = errors.New("family already exists")
	ErrNoFamily              = errors.New("no family")
	ErrNotMainUser           = errors.New("not main user")
	ErrMainUser              = errors.New("main user")
	ErrNotAvailableForDelete = errors.New("not available for delete")
	ErrNotAvailableForAdd    = errors.New("not available for add")
)

const (
	DefaultImage               = "default_avatar.webp"
	DefaultMedicine            = "default_medicine.webp"
	UserObjectsBucketName      = "avatars"
	MedicinesObjectsBucketName = "medicines"
	SessionRequired            = "Session required"
	UserIsUnauthorized         = "User is unauthorized"
	NoRequestID                = "No RequestID in context"
	UserIsLoggedOut            = "User is logged out"
	UserCanBeLoggedIn          = "User can be logged in"
	EmailConfirmed             = "Email confirmed"
	FileTypeIsNotSupported     = "File type is not supported"
	ProfileIsEdited            = "Profile is edited"
	FamilyIsCreated            = "Family is created"
	FamilyIsDeleted            = "Family is deleted"
	UserIsDeleted              = "User is deleted"
	MemberIsDeleted            = "Member is deleted"
	MedicineIsDeleted          = "Medicine is deleted"
	MemberIsAdded              = "Member is added"
	MedicineIsAdded            = "Medicine is added"
	NotMainUser                = "Not main user"
	InvitationIsSent           = "Invitation is sent"
	InvitationIsAccepted       = "Invitation is accepted"
	MedicineIsEdited           = "Medicine is eddited"
)

const (
	SignupURL             = "/api/v1/signup"
	LoginURL              = "/api/v1/login"
	LogoutURL             = "/api/v1/logout"
	ConfirmEmailURL       = "/confirm"
	ProfileURL            = "/api/v1/profile"
	EditURL               = "/api/v1/edit"
	AvatarURL             = "/api/v1/avatar"
	CsrfURL               = "/api/v1/csrf"
	AcceptInvitationURL   = "/accept"
	InviteURL             = "/api/v1/invite"
	CreateFamilyURL       = "/api/v1/create"
	DeleteFamilyURL       = "/api/v1/delete"
	RemoveMemberUrl       = "/api/v1/remove/member"
	RemoveUserUrl         = "/api/v1/remove/user"
	AddMembersToFamilyURL = "/api/v1/add"
	GetFamilyURL          = "/api/v1/family"
	DeleteMedicine        = "/api/v1/remove/medicine"
	AddMedicineURL        = "/api/v1/add/medicine"
	GetMedicineURL        = "/api/v1/medicine"
	EditMedicineURL       = "/api/v1/edit/medicine"
)

var (
	ImageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}
)

func RespError(ctx echo.Context, logger *zap.SugaredLogger, requestID, errorMsg string, status int) error {
	logger.Error(
		zap.String("ID", requestID),
		zap.String("ERROR", errorMsg),
		zap.Int("ANSWER STATUS", status),
	)
	resp, err := easyjson.Marshal(&models.Response{
		Status:  status,
		Message: errorMsg,
	})
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSONBlob(status, resp)
}

func DefaultUserChecks(ctx echo.Context, logger *zap.SugaredLogger) (int64, string, error) {
	requestID, ok := ctx.Get("REQUEST_ID").(string)
	if !ok {
		err := RespError(ctx, logger, requestID, NoRequestID, http.StatusInternalServerError)
		if err != nil {
			return 0, "", err
		}
		return 0, "", errors.New("")
	}

	userID, ok := ctx.Get("USER_ID").(int64)
	if !ok {
		err := RespError(ctx, logger, requestID, SessionRequired, http.StatusBadRequest)
		if err != nil {
			return 0, "", err
		}
		return 0, "", errors.New("")
	}

	if userID == -1 {
		err := RespError(ctx, logger, requestID, UserIsUnauthorized, http.StatusUnauthorized)
		if err != nil {
			return 0, "", err
		}
		return userID, "", errors.New("")
	}
	return userID, requestID, nil
}
