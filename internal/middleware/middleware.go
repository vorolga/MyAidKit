package middleware

import (
	"context"
	"main/internal/csrf"
	auth "main/internal/microservices/auth/proto"
	"main/internal/models"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Middleware struct {
	authMicroservice auth.AuthClient
	logger           *zap.SugaredLogger
}

func NewMiddleware(authMicroservice auth.AuthClient, logger *zap.SugaredLogger) *Middleware {
	return &Middleware{
		authMicroservice: authMicroservice,
		logger:           logger,
	}
}

func (m Middleware) Register(router *echo.Echo) {
	router.Use(m.CheckAuthorization())
	router.Use(m.CORS())
	router.Use(m.AccessLog())
	router.Use(m.CSRF())
}

func (m Middleware) CheckAuthorization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cookie, err := ctx.Cookie("Session_cookie")
			var userID int64
			userID = -1
			if err == nil {
				data := &auth.Cookie{Cookie: cookie.Value}
				id, err := m.authMicroservice.CheckAuthorization(context.Background(), data)
				if err != nil {
					cookie = &http.Cookie{
						Name:       "",
						Value:      "",
						Path:       "",
						Domain:     "",
						Expires:    time.Now().AddDate(0, 0, -1),
						RawExpires: "",
						MaxAge:     0,
						Secure:     false,
						HttpOnly:   false,
						SameSite:   0,
						Raw:        "",
						Unparsed:   nil,
					}
					ctx.SetCookie(cookie)
					ctx.Set("USER_ID", int64(-1))
					return next(ctx)
				}
				userID = id.ID
			}
			if err != nil {
				cookie = &http.Cookie{
					Name:       "",
					Value:      "",
					Path:       "",
					Domain:     "",
					Expires:    time.Now().AddDate(0, 0, -1),
					RawExpires: "",
					MaxAge:     0,
					Secure:     false,
					HttpOnly:   false,
					SameSite:   0,
					Raw:        "",
					Unparsed:   nil,
				}
				ctx.SetCookie(cookie)
			}

			ctx.Set("USER_ID", userID)

			return next(ctx)
		}
	}
}

func (m Middleware) CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          nil,
		AllowOrigins:     []string{"http://localhost:8080", "http://94.139.246.97:8080", "http://myaidkit.ru:8080", "http://myaidkit.ru", "myaidkit.ru:8080", "https://myaidkit.ru"},
		AllowOriginFunc:  nil,
		AllowMethods:     nil,
		AllowHeaders:     []string{"Accept", "Cache-Control", "Content-Type", "X-Requested-With", "csrf-token", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		ExposeHeaders:    nil,
		MaxAge:           84600,
	})
}

func (m Middleware) AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newUUID, _ := uuid.NewV4()

			start := time.Now()
			ctx.Set("REQUEST_ID", newUUID.String())

			m.logger.Info(
				zap.String("ID", newUUID.String()),
				zap.String("URL", ctx.Request().URL.Path),
				zap.String("METHOD", ctx.Request().Method),
			)

			err := next(ctx)

			responseTime := time.Since(start)
			m.logger.Info(
				zap.String("ID", newUUID.String()),
				zap.Duration("TIME FOR ANSWER", responseTime),
			)

			return err
		}
	}
}

func (m Middleware) CSRF() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.Request().Method == "PUT" {
				cookie, err := ctx.Cookie("Session_cookie")
				if err != nil {
					m.logger.Debug(
						zap.String("COOKIE", err.Error()),
						zap.Int("ANSWER STATUS", http.StatusInternalServerError),
					)

					return ctx.JSON(http.StatusInternalServerError, &models.Response{
						Status:  http.StatusInternalServerError,
						Message: err.Error(),
					})
				}

				GetToken := ctx.Request().Header.Get("csrf-token")

				isValidCsrf, err := csrf.Tokens.Check(cookie.Value, GetToken)
				if err != nil {
					return ctx.JSON(http.StatusInternalServerError, &models.Response{
						Status:  http.StatusInternalServerError,
						Message: err.Error(),
					})
				}

				if !isValidCsrf {
					return ctx.JSON(http.StatusForbidden, &models.Response{
						Status:  http.StatusForbidden,
						Message: "Wrong csrf token",
					})
				}
			}
			return next(ctx)
		}
	}
}
