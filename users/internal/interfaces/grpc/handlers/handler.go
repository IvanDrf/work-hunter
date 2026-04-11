package handlers

import (
	user_api "github.com/IvanDrf/work-hunter/pkg/user-api"
	service "github.com/IvanDrf/work-hunter/users/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/users/internal/logger"
)

type Handler struct {
	service.UserService
	user_api.UnimplementedUserServer

	log *logger.Logger
}

func NewHandler(service service.UserService) *Handler {
	return &Handler{
		UserService: service,
	}
}

func (h *Handler) Close() {
	h.UserService.Close()
}
