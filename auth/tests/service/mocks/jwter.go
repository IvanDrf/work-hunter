package mocks

import (
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/auth/pkg"
)

var (
	// Secret - valid Secret for  auth jwt tokens, auth service use this Secret
	Secret = rules.GenerateToken()

	// InvalidSecret - invalid secret for auth jwt tokens, auth service doesn't use this secret
	InvalidSecret = rules.GenerateToken()

	Jwter        = pkg.NewJwt(Secret, 1*time.Minute, 2*time.Minute)
	InvalidJwter = pkg.NewJwt(InvalidSecret, 1*time.Minute, 2*time.Minute)
)
