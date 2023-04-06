package composites

import (
	"main/internal/microservices/auth"
	"main/internal/microservices/auth/repository"
	"main/internal/microservices/auth/usecase"
)

type AuthComposite struct {
	Storage auth.Storage
	Service *usecase.Service
}

func NewAuthComposite(postgresComposite *PostgresDBComposite, redisComposite *RedisComposite) (*AuthComposite, error) {
	storage := repository.NewStorage(postgresComposite.db, redisComposite.redis)
	service := usecase.NewService(storage)
	return &AuthComposite{
		Storage: storage,
		Service: service,
	}, nil
}
