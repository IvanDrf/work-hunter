package mocks

import (
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/jwt"
)

var (
	// Secret - valid Secret for  auth jwt tokens, auth service use this Secret
	Secret = rules.GenerateToken()

	// InvalidSecret - invalid secret for auth jwt tokens, auth service doesn't use this secret
	InvalidSecret = rules.GenerateToken()

	Jwter        = jwt.NewJwt(Secret, 1*time.Minute, 2*time.Minute)
	InvalidJwter = jwt.NewJwt(InvalidSecret, 1*time.Minute, 2*time.Minute)
)
