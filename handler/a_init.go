package handler

import (
	"go-cloud/services"
)

type Handler struct {
	serv *services.Services
}

func NewHandler(serv *services.Services) *Handler {
	return &Handler{serv: serv}
}
