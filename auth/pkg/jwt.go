package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidJWT = errors.New("invalid jwt token")

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

func (j *Jwt) CreateTokens(userID uuid.UUID) (string, string, error) {
	access, err := j.createToken(userID, j.accessTime)
	if err != nil {
		return "", "", err
	}

	refresh, err := j.createToken(userID, j.refreshTime)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (j *Jwt) createToken(userID uuid.UUID, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *Jwt) GetTokenClaims(token string) (jwt.MapClaims, error) {
	t, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWT
		}

		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		return claims, nil
	}

	return nil, ErrInvalidJWT
}

func (j *Jwt) RefreshTokens(refresh string) (string, string, error) {
	claims, err := j.GetTokenClaims(refresh)
	if err != nil {
		return "", "", err
	}

	userID, ok := claims["user_id"].(uuid.UUID)
	if !ok {
		return "", "", ErrInvalidJWT
	}

	return j.CreateTokens(userID)
}
