package service

import "context"

type VerificationService interface {
	SendVerificationEmail(ctx context.Context, email string) error
	ResendVerificationEmail(ctx context.Context, access string) error
	VerifyEmailByToken(ctx context.Context, token string) (string, string, error)

	Close()
}
