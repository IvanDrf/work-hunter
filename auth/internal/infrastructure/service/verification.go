package service

import (
	"context"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
)

type VerificationService struct {
	emailProducer service.EmailProducer

	userRepo  repo.UserRepo
	tokenRepo repo.TokenRepo

	jwter jwt.Jwter
}

func NewVerificationService(emailProducer service.EmailProducer, userRepo repo.UserRepo, tokenRepo repo.TokenRepo, jwter jwt.Jwter) *VerificationService {
	return &VerificationService{
		emailProducer: emailProducer,
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		jwter:         jwter,
	}
}

func (v *VerificationService) Close() {
	v.emailProducer.Close()
	v.userRepo.Close()
}

func (v *VerificationService) SendVerificationEmail(ctx context.Context, email string) error {
	token := models.NewToken()

	err := v.tokenRepo.CreateToken(ctx, email, token)
	if err != nil {
		slog.Error("verif:SendVerifEmail service error", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't create token for user",
			Code:    models.ErrCodeInternal,
		}
	}

	err = v.emailProducer.SendEmailInQueue(ctx, &models.EmailMessage{
		Email: email,
		Token: token,
	})

	if err != nil {
		slog.Error("verif:SendVerifEmail service error", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't send verification message",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("verif:SendVerifEmail service success")
	return nil
}

func (v *VerificationService) VerifyEmailByToken(ctx context.Context, token string) (string, string, error) {
	email, exp, err := v.tokenRepo.FindEmailExpByToken(ctx, token)
	if err != nil {
		slog.Error("verif:VerifyEmailByToken error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't find user with that token",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	if !(&models.Token{Exp: exp}).IsTokenValid() {
		slog.Info("verif:VerifyEmailByToken token is expired")
		return "", "", models.Error{
			Message: "token is outdated",
			Code:    models.ErrOutdatedToken,
		}
	}

	user, err := v.userRepo.FindUser(ctx, email)
	if err != nil {
		slog.Error("verif:VerifyEmailByToken error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't find user with that email",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	err = v.userRepo.VerifyEmail(ctx, email)
	if err != nil {
		slog.Error("verif:VerifyEmailByToken error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't verify user email",
			Code:    models.ErrCodeInternal,
		}
	}

	err = v.tokenRepo.DeleteToken(ctx, token)
	if err != nil {
		slog.Error("verif:VerifyByEmailToken internal error, can't delete token", slog.String("error", err.Error()))
	}

	access, refresh, err := v.jwter.CreateTokens(user.ID, true)
	if err != nil {
		slog.Error("verif:VerifyEmailByToken error", slog.String("error", err.Error()))
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	slog.Info("verif:VerifyEmailByToken service success")
	return access, refresh, nil
}
