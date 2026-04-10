package jwt

import (
	"errors"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidJWT = errors.New("invalid jwt token")

type claims struct {
	models.JwtPayload

	jwt.RegisteredClaims
}

type Jwt struct {
	secret []byte

	accessTime  time.Duration
	refreshTime time.Duration
}

func NewJwt(secret string, accessTime time.Duration, refreshTime time.Duration) *Jwt {
	return &Jwt{
		secret:      []byte(secret),
		accessTime:  accessTime,
		refreshTime: refreshTime,
	}
}

func (j *Jwt) CreateTokens(payload *models.JwtPayload) (string, string, error) {
	access, err := j.createToken(payload, j.accessTime)
	if err != nil {
		return "", "", err
	}

	refresh, err := j.createToken(payload, j.refreshTime)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (j *Jwt) createToken(payload *models.JwtPayload, duration time.Duration) (string, error) {
	claims := claims{
		JwtPayload: *payload,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *Jwt) GetPayload(token string) (*models.JwtPayload, error) {
	claims := &claims{}

	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWT
		}

		return j.secret, nil
	})

	if err != nil || !t.Valid {
		return nil, err
	}

	return &models.JwtPayload{
		UserID:      claims.UserID,
		Verificated: claims.Verificated,
		Role:        claims.Role,
	}, err
}
