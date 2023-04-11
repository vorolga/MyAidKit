package delivery

import (
	"context"
	"main/internal/constants"
	"main/internal/csrf"
	profile "main/internal/microservices/profile/proto"
	"main/internal/models"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"time"

	"github.com/mailru/easyjson"

	"github.com/labstack/echo/v4"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type profileHandler struct {
	logger *zap.SugaredLogger

	profileMicroservice profile.ProfileClient
}

func NewProfileHandler(logger *zap.SugaredLogger, profile profile.ProfileClient) *profileHandler {
	return &profileHandler{profileMicroservice: profile, logger: logger}
}

func (p *profileHandler) Register(router *echo.Echo) {
	router.GET(constants.ProfileURL, p.GetUserProfile())
	router.PUT(constants.EditURL, p.EditProfile())
	router.PUT(constants.AvatarURL, p.EditAvatar())
	router.GET(constants.CsrfURL, p.GetCsrf())
	router.GET(constants.AcceptInvitationURL, p.AcceptInvitation())
	router.POST(constants.InviteURL, p.Invite())
	router.POST(constants.CreateFamilyURL, p.CreateFamily())
	router.DELETE(constants.DeleteFamilyURL, p.DeleteFamily())
	router.DELETE(constants.RemoveUserUrl, p.RemoveUser())
	router.DELETE(constants.RemoveMemberUrl, p.RemoveMember())
	router.POST(constants.AddMembersToFamilyURL, p.AddMember())
	router.GET(constants.GetFamilyURL, p.GetFamily())
	router.DELETE(constants.DeleteMedicine, p.DeleteMedicine())
	router.POST(constants.AddMedicineURL, p.AddMedicine())
	router.GET(constants.GetMedicineURL, p.GetMedicine())
}

func (p *profileHandler) ParseError(ctx echo.Context, requestID string, err error) error {
	if getErr, ok := status.FromError(err); ok {
		switch getErr.Code() {
		case codes.Internal:
			p.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: getErr.Message(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		case codes.Unavailable:
			p.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: getErr.Message(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		case codes.InvalidArgument:
			p.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusBadRequest),
			)
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusBadRequest,
				Message: getErr.Message(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusBadRequest, resp)
		case codes.PermissionDenied:
			p.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusForbidden),
			)
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusForbidden,
				Message: getErr.Message(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusForbidden, resp)
		case codes.AlreadyExists:
			p.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusBadRequest),
			)
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusBadRequest,
				Message: getErr.Message(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusBadRequest, resp)
		}
	}
	return nil
}

func (p *profileHandler) GetUserProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}
		data := &profile.UserID{ID: userID}
		userData, err := p.profileMicroservice.GetUserProfile(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
			userData,
		)

		profileData := models.ProfileUserDTO{
			ID:      userID,
			Name:    userData.Name,
			Surname: userData.Surname,
			Email:   userData.Email,
			Avatar:  userData.Avatar,
			Date:    userData.Date,
			Main:    userData.Main,
			Adult:   userData.Adult,
		}

		sanitizer := bluemonday.UGCPolicy()
		profileData.Name = sanitizer.Sanitize(profileData.Name)
		profileData.Surname = sanitizer.Sanitize(profileData.Surname)
		resp, err := easyjson.Marshal(&models.ResponseUserProfile{
			Status:   http.StatusOK,
			UserData: &profileData,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) EditAvatar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		file, err := ctx.FormFile("file")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		src, err := file.Open()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		buffer := make([]byte, file.Size)
		_, err = src.Read(buffer)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}
		err = src.Close()
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		file, err = ctx.FormFile("file")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}
		src, err = file.Open()
		defer func(src multipart.File) {
			err = src.Close()
			if err != nil {
				return
			}
		}(src)
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		fileType := http.DetectContentType(buffer)

		// Validate File Type
		if _, ex := constants.ImageTypes[fileType]; !ex {
			return constants.RespError(ctx, p.logger, requestID, constants.FileTypeIsNotSupported, http.StatusBadRequest)
		}

		uploadData := &profile.UploadInputFile{
			ID:          userID,
			File:        buffer,
			Size:        file.Size,
			ContentType: fileType,
			BucketName:  constants.UserObjectsBucketName,
		}

		fileName, err := p.profileMicroservice.UploadAvatar(context.Background(), uploadData)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		editData := &profile.EditAvatarData{
			ID:     userID,
			Avatar: fileName.Name,
		}

		_, err = p.profileMicroservice.EditAvatar(context.Background(), editData)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.ProfileIsEdited,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) EditProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		userData := models.EditProfileDTO{}

		if err = ctx.Bind(&userData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data := &profile.EditProfileData{
			ID:       userID,
			Name:     userData.Name,
			Surname:  userData.Surname,
			Password: userData.Password,
			Date:     userData.Date,
		}

		_, err = p.profileMicroservice.EditProfile(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.ProfileIsEdited,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) GetCsrf() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			return constants.RespError(ctx, p.logger, requestID, constants.NoRequestID, http.StatusInternalServerError)
		}

		cookie, err := ctx.Cookie("Session_cookie")
		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		token, err := csrf.Tokens.Create(cookie.Value, time.Now().Add(time.Hour).Unix())

		if err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)
		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: token,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) CreateFamily() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}
		data := &profile.UserID{ID: userID}
		_, err = p.profileMicroservice.CreateFamily(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.FamilyIsCreated,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) DeleteFamily() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}
		data := &profile.UserID{ID: userID}
		_, err = p.profileMicroservice.DeleteFamily(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.FamilyIsDeleted,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) RemoveUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}
		data := &profile.Delete{UserID: &profile.UserID{ID: userID}, UserToDelete: &profile.UserID{ID: 0}}

		userData := models.UserIDDTO{}

		if err = ctx.Bind(&userData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data.UserToDelete.ID = userData.ID
		_, err = p.profileMicroservice.DeleteFromFamily(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.UserIsDeleted,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) RemoveMember() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}
		data := &profile.Delete{UserID: &profile.UserID{ID: userID}, UserToDelete: &profile.UserID{ID: 0}}

		userData := models.UserIDDTO{}

		if err = ctx.Bind(&userData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data.UserToDelete.ID = userData.ID
		_, err = p.profileMicroservice.DeleteMember(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.MemberIsDeleted,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) AddMember() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)

		fileName := constants.DefaultImage
		file, err := ctx.FormFile("file")
		if err == nil {
			src, err := file.Open()
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}

			buffer := make([]byte, file.Size)
			_, err = src.Read(buffer)
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}
			err = src.Close()
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}

			file, err = ctx.FormFile("file")
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}
			src, err = file.Open()
			defer func(src multipart.File) {
				err = src.Close()
				if err != nil {
					return
				}
			}(src)
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}

			fileType := http.DetectContentType(buffer)

			// Validate File Type
			if _, ex := constants.ImageTypes[fileType]; !ex {
				return constants.RespError(ctx, p.logger, requestID, constants.FileTypeIsNotSupported, http.StatusBadRequest)
			}

			uploadData := &profile.UploadInputFile{
				ID:          userID,
				File:        buffer,
				Size:        file.Size,
				ContentType: fileType,
				BucketName:  constants.UserObjectsBucketName,
			}

			name, err := p.profileMicroservice.UploadAvatar(context.Background(), uploadData)
			if err != nil {
				return p.ParseError(ctx, requestID, err)
			}

			fileName = name.Name
		} else {
			if err.Error() != "http: no such file" {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}
		}

		data := &profile.MemberData{
			IDFamily:   0,
			IDMainUser: userID,
			Name:       "",
			Avatar:     fileName,
		}

		userData := models.MemberDTO{}

		if err = ctx.Bind(&userData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data.Name = userData.Name
		_, err = p.profileMicroservice.AddMember(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.MemberIsAdded,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) GetFamily() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}
		data := &profile.UserID{
			ID: userID,
		}

		family, err := p.profileMicroservice.GetFamily(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		membersResult := make([]models.Member, 0)
		for _, member := range family.ResponseMemberData {
			membersResult = append(membersResult, models.Member{
				ID:     member.ID,
				Name:   member.Name,
				Avatar: member.Avatar,
				Adult:  member.IsAdult,
				User:   member.IsUser,
			})
		}

		resp, err := easyjson.Marshal(&models.ResponseMembers{
			Status:  200,
			Members: membersResult,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) Invite() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		userData := models.InviteUserDTO{}

		if err = ctx.Bind(&userData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data := &profile.UserID{ID: userID}
		hasFamilyResp, err := p.profileMicroservice.HasFamily(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		if hasFamilyResp.IDMainUser != userID {
			resp, err := easyjson.Marshal(&models.Response{
				Status:  200,
				Message: constants.NotMainUser,
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusOK, resp)
		}

		dataProfile := &profile.UserID{ID: userID}
		profileData, err := p.profileMicroservice.GetUserProfile(context.Background(), dataProfile)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		if profileData.Email == userData.Email {
			resp, err := easyjson.Marshal(&models.Response{
				Status:  400,
				Message: constants.ErrWrongData.Error(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusBadRequest, resp)
		}

		email := &profile.EmailData{Email: userData.Email}
		exists, err := p.profileMicroservice.UserExists(context.Background(), email)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		if exists.Exists == false {
			resp, err := easyjson.Marshal(&models.Response{
				Status:  400,
				Message: constants.ErrWrongData.Error(),
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusBadRequest, resp)
		}

		from := "vorrovvorrov@gmail.com"
		password := os.Getenv("EMAILPASSWORD")

		toList := []string{userData.Email}

		host := "smtp.gmail.com"
		port := "587"

		msg := "Вам пришло приглашение в семью\r\n" +
			"Чтобы принять, перейдите по ссылке: " + "http://" + os.Getenv("HOST") +
			"/accept?family=" + strconv.Itoa(int(hasFamilyResp.IDFamily)) +
			"&email=" + userData.Email + "&adult=" + strconv.FormatBool(userData.Adult)

		body := []byte(msg)

		authSMTP := smtp.PlainAuth("", from, password, host)
		err = smtp.SendMail(host+":"+port, authSMTP, from, toList, body)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}

		resp, err := easyjson.Marshal(&models.Response{
			Status:  200,
			Message: constants.InvitationIsSent,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) AcceptInvitation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestID, ok := ctx.Get("REQUEST_ID").(string)
		if !ok {
			p.logger.Error(
				zap.String("ERROR", constants.NoRequestID),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			resp, err := easyjson.Marshal(&models.Response{
				Status:  http.StatusInternalServerError,
				Message: constants.NoRequestID,
			})
			if err != nil {
				return ctx.NoContent(http.StatusInternalServerError)
			}
			return ctx.JSONBlob(http.StatusInternalServerError, resp)
		}

		family := ctx.QueryParam("family")
		email := ctx.QueryParam("email")
		adult := ctx.QueryParam("adult")

		familyID, err := strconv.Atoi(family)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}

		isAdult, err := strconv.ParseBool(adult)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}

		data := &profile.AddToFamily{
			ID:      int64(familyID),
			Email:   email,
			IsAdult: isAdult,
		}
		_, err = p.profileMicroservice.AcceptInvitationToFamily(context.Background(), data)
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.InvitationIsAccepted,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) AddMedicine() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)

		fileName := constants.DefaultMedicine
		file, err := ctx.FormFile("file")
		if err == nil {
			src, err := file.Open()
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}

			buffer := make([]byte, file.Size)
			_, err = src.Read(buffer)
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}
			err = src.Close()
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}

			file, err = ctx.FormFile("file")
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}
			src, err = file.Open()
			defer func(src multipart.File) {
				err = src.Close()
				if err != nil {
					return
				}
			}(src)
			if err != nil {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}

			fileType := http.DetectContentType(buffer)

			// Validate File Type
			if _, ex := constants.ImageTypes[fileType]; !ex {
				return constants.RespError(ctx, p.logger, requestID, constants.FileTypeIsNotSupported, http.StatusBadRequest)
			}

			uploadData := &profile.UploadInputFile{
				ID:          userID,
				File:        buffer,
				Size:        file.Size,
				ContentType: fileType,
				BucketName:  constants.MedicinesObjectsBucketName,
			}

			name, err := p.profileMicroservice.UploadAvatar(context.Background(), uploadData)
			if err != nil {
				return p.ParseError(ctx, requestID, err)
			}

			fileName = name.Name
		} else {
			if err.Error() != "http: no such file" {
				return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusInternalServerError)
			}
		}

		medicineData := models.AddMedicineDTO{}

		if err = ctx.Bind(&medicineData); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data := &profile.AddMed{
			UserID: userID,
			Medicine: &profile.Medicine{
				Image:     fileName,
				Name:      medicineData.Name,
				IsTablets: medicineData.IsTablets,
				Count:     medicineData.Count,
			},
		}

		_, err = p.profileMicroservice.AddMedicine(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.MedicineIsAdded,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) DeleteMedicine() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		_, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		dataMedicine := models.MedecineIDDTO{}

		if err = ctx.Bind(&dataMedicine); err != nil {
			return constants.RespError(ctx, p.logger, requestID, err.Error(), http.StatusBadRequest)
		}

		data := &profile.DeleteMed{MedicineID: dataMedicine.ID}
		_, err = p.profileMicroservice.DeleteMedicine(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		p.logger.Info(
			zap.String("ID", requestID),
			zap.Int("ANSWER STATUS", http.StatusOK),
		)

		resp, err := easyjson.Marshal(&models.Response{
			Status:  http.StatusOK,
			Message: constants.MedicineIsDeleted,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}

func (p *profileHandler) GetMedicine() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, requestID, err := constants.DefaultUserChecks(ctx, p.logger)
		if err != nil {
			return err
		}

		data := &profile.UserID{
			ID: userID,
		}
		medicines, err := p.profileMicroservice.GetMedicine(context.Background(), data)
		if err != nil {
			return p.ParseError(ctx, requestID, err)
		}

		medicineResult := make([]models.Medicine, 0)
		for _, medicine := range medicines.MedicineArr {
			medicineResult = append(medicineResult, models.Medicine{
				ID:        medicine.ID,
				Name:      medicine.Medicine.Name,
				Image:     medicine.Medicine.Image,
				IsTablets: medicine.Medicine.IsTablets,
				Count:     medicine.Medicine.Count,
			})
		}

		resp, err := easyjson.Marshal(&models.ResponseMedicine{
			Status:   200,
			Medicine: medicineResult,
		})
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
		return ctx.JSONBlob(http.StatusOK, resp)
	}
}
