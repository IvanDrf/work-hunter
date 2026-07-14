package http

import "net/http"

func (s *server) registerRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/register", s.handlers.RegisterUser)
	mux.HandleFunc("POST /api/login", s.handlers.LoginUser)
	mux.HandleFunc("POST /api/password", s.handlers.ChangeUserPassword)

	s.server.Handler = mux
}
