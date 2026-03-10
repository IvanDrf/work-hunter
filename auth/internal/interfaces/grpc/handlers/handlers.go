package handlers

import "github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"

type Handler struct {
	authService service.AuthService
}

func NewHandler(service service.AuthService) Handler {
	return Handler{
		authService: service,
	}
}
