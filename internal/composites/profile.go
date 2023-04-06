package composites

import (
	"main/internal/microservices/profile"
	"main/internal/microservices/profile/repository"
	"main/internal/microservices/profile/usecase"
)

type ProfileComposite struct {
	Storage profile.Storage
	Service *usecase.Service
}

func NewProfileComposite(postgresComposite *PostgresDBComposite, minioComposite *MinioComposite,
	redisComposite *RedisComposite) (*ProfileComposite, error) {
	storage := repository.NewStorage(postgresComposite.db, minioComposite.client, redisComposite.redis)
	service := usecase.NewService(storage)
	return &ProfileComposite{
		Storage: storage,
		Service: service,
	}, nil
}
