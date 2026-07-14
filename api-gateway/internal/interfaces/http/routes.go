package http

import "net/http"

func (s *server) registerRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/auth/register", s.handlers.RegisterUser)
	mux.HandleFunc("POST /api/auth/login", s.handlers.LoginUser)
	mux.HandleFunc("POST /api/auth/password", s.handlers.ChangeUserPassword)
	mux.HandleFunc("POST /api/auth/tokens", s.handlers.RefreshTokens)

	s.server.Handler = mux
}
