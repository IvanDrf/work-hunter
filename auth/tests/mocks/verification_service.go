package mocks

import "context"

const (
	Access  = "access"
	Refresh = "refresh"
)

type VerificationService struct {
	messages map[string]struct{}
	repo     *UserRepo
}

func NewVerificationService() *VerificationService {
	return &VerificationService{
		messages: map[string]struct{}{},
	}
}

func (v *VerificationService) SendVerificationEmail(ctx context.Context, email string) error {
	v.messages[email] = struct{}{}
	return nil
}

func (v *VerificationService) VerifyEmailByToken(ctx context.Context, token string) (string, string, error) {
	return Access, Refresh, nil
}

func (v *VerificationService) Close() {
	v.messages = nil
	v.repo.Close()
}
