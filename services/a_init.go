package services

import (
	"go-cloud/repository"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	repo        *repository.Repository
	redisClient *redis.Client
}

func NewServices(repo *repository.Repository, redisClient *redis.Client) *Services {
	return &Services{repo: repo, redisClient: redisClient}
}
