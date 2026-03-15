package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/jwt"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/repo"
	"github.com/IvanDrf/work-hunter/auth/internal/domain/ports/service"
)

type VerificationService struct {
	emailProducer service.EmailProducer
	jwter         jwt.Jwter

	userRepo repo.UserRepo
}

func NewVerificationService(emailProducer service.EmailProducer, userRepo repo.UserRepo, jwter jwt.Jwter) *VerificationService {
	return &VerificationService{
		emailProducer: emailProducer,
		userRepo:      userRepo,
		jwter:         jwter,
	}
}

func (v *VerificationService) SendVerificationEmail(ctx context.Context, email string) error {
	err := v.emailProducer.SendEmailInQueue(ctx, &models.EmailMessage{
		Email: email,
		Token: v.createToken(),
	})

	if err != nil {
		return models.Error{
			Message: "can't send verification message",
			Code:    models.ErrCodeInternal,
		}
	}

	return nil
}

func (v *VerificationService) VerifyEmail(ctx context.Context, email string) (string, string, error) {
	user, err := v.userRepo.FindUser(ctx, email)
	if err != nil {
		return "", "", models.Error{
			Message: "can't find user with that email",
			Code:    models.ErrCodeInvalidEmail,
		}
	}

	err = v.userRepo.VerifyEmail(ctx, email)
	if err != nil {
		return "", "", models.Error{
			Message: "can't verify user email",
			Code:    models.ErrCodeInternal,
		}
	}

	access, refresh, err := v.jwter.CreateTokens(user.ID, true)
	if err != nil {
		return "", "", models.Error{
			Message: "can't create jwt tokens for user",
			Code:    models.ErrCodeInternal,
		}
	}

	return access, refresh, nil
}

func (v *VerificationService) createToken() string {
	buff := make([]byte, 32)

	rand.Read(buff)
	return hex.EncodeToString(buff)
}
