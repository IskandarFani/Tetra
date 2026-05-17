package services

import (
	"go-cloud/repository"
)

type Services struct {
	repo *repository.Repository
}

func NewServices(repo *repository.Repository) *Services {
	return &Services{repo: repo}
}
