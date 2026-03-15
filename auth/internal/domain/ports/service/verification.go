package service

import "context"

type VerificationService interface {
	SendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, email string) error
}
