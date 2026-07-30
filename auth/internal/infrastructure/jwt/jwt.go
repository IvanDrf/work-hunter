package jwt

import (
	"errors"
	"time"

	"github.com/IvanDrf/work-hunter/auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

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

var errInvalidJwt = errors.New(models.InvalidJWTMsg)

func (j *Jwt) GetPayload(token string) (*models.JwtPayload, error) {
	data := &claims{}

	t, err := jwt.ParseWithClaims(token, data, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidJwt
		}

		return j.secret, nil
	})

	if err != nil || !t.Valid {
		return nil, err
	}

	return &models.JwtPayload{
		UserID:      data.UserID,
		Verificated: data.Verificated,
		Role:        data.Role,
	}, err
}

func (j *Jwt) createToken(payload *models.JwtPayload, duration time.Duration) (string, error) {
	data := claims{
		JwtPayload: *payload,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(j.secret)
}
