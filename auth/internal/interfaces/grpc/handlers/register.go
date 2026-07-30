package handlers

import (
	"context"
	"log/slog"

	auth_api "github.com/IvanDrf/work-hunter/pkg/auth-api"
)

func (h *Handler) Register(ctx context.Context, user *auth_api.User) (*auth_api.JwtTokens, error) {
	slog.Info("Register got request")

	access, refresh, err := h.authService.RegisterUser(ctx, user.GetEmail(), user.GetPassword(), user.GetRole().String())
	if err != nil {
		return nil, handleError(err, "Register error")
	}

	const attempts = 5
	err = retry(attempts, func() error {
		err = h.verificationService.SendVerificationEmail(ctx, user.GetEmail())
		if err != nil {
			slog.Error("Register error", slog.String("error", err.Error()))
			return err
		}

		return nil
	})
	if err != nil {
		slog.Error("Cant send verification email after all attempts", slog.String("error", err.Error()))
	}

	slog.Info("Register successful response")
	return &auth_api.JwtTokens{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func retry(attempts int, fn func() error) error {
	var err error

	for range attempts + 1 {
		if err = fn(); err == nil {
			return nil
		}
	}

	return err
}
