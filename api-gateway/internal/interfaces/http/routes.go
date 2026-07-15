package http

import "net/http"

func (s *server) registerRoutes() {
	mux := http.NewServeMux()

	s.registerHealthRoute(mux)
	s.registerAuthRoutes(mux)

	s.server.Handler = mux
}

func (s *server) registerHealthRoute(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.handlers.Health)
}

func (s *server) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/user/register", s.handlers.RegisterUser)
	mux.HandleFunc("POST /api/auth/user/login", s.handlers.LoginUser)
	mux.HandleFunc("POST /api/auth/user/password", s.handlers.ChangeUserPassword)
	mux.HandleFunc("DELETE /api/auth/user", s.handlers.DeleteUser)
	mux.HandleFunc("POST /api/auth/tokens", s.handlers.RefreshTokens)
}
