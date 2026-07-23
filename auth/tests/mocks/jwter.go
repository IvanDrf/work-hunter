package mocks

import (
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/rules"
	"github.com/IvanDrf/work-hunter/auth/internal/infrastructure/jwt"
)

const (
	accessTime  = 1 * time.Minute
	refreshTime = 2 * time.Minute
)

var (
	// Secret - valid Secret for  auth jwt tokens, auth service use this Secret.
	Secret = rules.GenerateToken()

	// InvalidSecret - invalid secret for auth jwt tokens, auth service doesn't use this secret.
	InvalidSecret = rules.GenerateToken()

	Jwter        = jwt.NewJwt(Secret, accessTime, refreshTime)
	InvalidJwter = jwt.NewJwt(InvalidSecret, accessTime, refreshTime)
)
