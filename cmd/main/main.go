package main

import (
	"log"
	api "main/internal/api/delivery"
	authMicroservice "main/internal/microservices/auth/proto"
	profileMicroservice "main/internal/microservices/profile/proto"
	"main/internal/middleware"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	echoServer := echo.New()

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, err := config.Build()
	if err != nil {
		log.Fatal("zap logger build error")
	}
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		err = prLogger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(prLogger)

	authConn, err := grpc.Dial(
		os.Getenv("AUTH_HOST")+":"+
			os.Getenv("AUTH_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		echoServer.Logger.Fatal("auth cant connect to grpc")
	}

	authManager := authMicroservice.NewAuthClient(authConn)

	authHandlers := api.NewAuthHandler(logger, authManager)
	authHandlers.Register(echoServer)

	profileConn, err := grpc.Dial(
		os.Getenv("PROFILE_HOST")+":"+os.Getenv("PROFILE_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		echoServer.Logger.Fatal("profile cant connect to grpc")
	}

	profileManager := profileMicroservice.NewProfileClient(profileConn)

	profileHandlers := api.NewProfileHandler(logger, profileManager)
	profileHandlers.Register(echoServer)

	middlewares := middleware.NewMiddleware(authManager, logger)
	middlewares.Register(echoServer)

	echoServer.Logger.Fatal(echoServer.Start(":1323"))
}
