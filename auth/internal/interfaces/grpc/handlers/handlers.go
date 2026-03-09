package handlers

import "github.com/IvanDrf/work-hunter/auth/internal/domain/service"

type Handler struct {
	authService service.AuthService
}

func NewHandler(service service.AuthService) Handler {
	return Handler{
		authService: service,
	}
}
