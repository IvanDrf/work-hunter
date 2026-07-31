package handlers

import (
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

type Handler struct {
	auth_api.UnimplementedAuthServer

	authService         service.AuthService
	verificationService service.VerificationService
}

func NewHandler(service service.AuthService, verificationService service.VerificationService) *Handler {
	return &Handler{
		authService:         service,
		verificationService: verificationService,
	}
}

func (h *Handler) Close() {
	h.authService.Close()
	h.verificationService.Close()
}
