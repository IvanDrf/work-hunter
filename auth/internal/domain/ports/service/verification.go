package service

import "context"

type VerificationService interface {
	SendVerificationEmail(ctx context.Context, email string) error
	VerifyEmailByToken(ctx context.Context, token string) (string, string, error)

	Close()
}
