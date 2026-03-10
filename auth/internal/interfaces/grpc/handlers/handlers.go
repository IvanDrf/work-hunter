package handlers

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

type Handler struct {
	authService service.AuthService
	auth_api.UnimplementedAuthServer
}

func NewHandler(service service.AuthService) *Handler {
	return &Handler{
		authService: service,
	}
}
