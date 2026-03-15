package handlers

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

type Handler struct {
	authService         service.AuthService
	verificationService service.VerificationService

	auth_api.UnimplementedAuthServer
}

func NewHandler(service service.AuthService, verificationService service.VerificationService) *Handler {
	return &Handler{
		authService:         service,
		verificationService: verificationService,
	}
}
