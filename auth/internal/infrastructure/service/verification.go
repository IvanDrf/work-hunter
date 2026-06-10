package service

import (
	"context"
	"log/slog"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/google/uuid"
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
	v.tokenRepo.Close()
}

func (v *VerificationService) SendVerificationEmail(ctx context.Context, email string) error {
	user, err := v.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return models.Error{
			Message: "can't find user with that email",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	return v.sendVerificationEmailForUser(ctx, user)
}

func (v *VerificationService) ResendVerificationEmail(ctx context.Context, access string) error {
	payload, err := v.jwter.GetPayload(access)
	if err != nil {
		return models.Error{
			Message: "invalid jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return models.Error{
			Message: "invalid userID in jwt token",
			Code:    models.ErrCodeInvalidJWT,
		}
	}

	user, err := v.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return models.Error{
			Message: "can't find user with that userID",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	return v.sendVerificationEmailForUser(ctx, user)
}

func (v *VerificationService) sendVerificationEmailForUser(ctx context.Context, user *models.User) error {
	if user.Verificated {
		return models.Error{
			Message: "user already verificated",
			Code:    models.ErrCodeUserAlreadyVerificated,
		}
	}

	token := rules.GenerateToken()

	err := v.tokenRepo.CreateToken(ctx, user.Email, token, rules.TokenTTL)
	if err != nil {
		slog.Error("verif:SendVerifEmail service error", slog.String("error", err.Error()))
		return models.Error{
			Message: "can't create token for user",
			Code:    models.ErrCodeInternal,
		}
	}

	err = v.emailProducer.SendEmailInQueue(ctx, &models.EmailMessage{
		Email: user.Email,
		Token: token,
		Exp:   rules.NewExpTime(),
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
	email := v.tokenRepo.FindEmailByToken(ctx, token)
	if email == "" {
		slog.Error("verif:VerifyEmailByToken error")
		return "", "", models.Error{
			Message: "can't find user with that token or token is outdated",
			Code:    models.ErrCodeUserNotFound,
		}
	}

	user, err := v.userRepo.FindUserByEmail(ctx, email)
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

	access, refresh, err := v.jwter.CreateTokens(&models.JwtPayload{
		UserID:      user.ID.String(),
		Verificated: true,
		Role:        user.Role,
	})
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
