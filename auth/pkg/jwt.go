package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidJWT = errors.New("invalid jwt token")

type claims struct {
	UserID string `json:"user_id"`

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
	claims := claims{
		UserID: userID.String(),

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *Jwt) GetUserID(token string) (uuid.UUID, error) {
	claims := &claims{}

	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWT
		}

		return j.secret, nil
	})

	if err != nil || !t.Valid {
		return uuid.UUID{}, err
	}

	return uuid.Parse(claims.UserID)
}

func (j *Jwt) RefreshTokens(refresh string) (string, string, error) {
	userID, err := j.GetUserID(refresh)
	if err != nil {
		return "", "", err
	}

	return j.CreateTokens(userID)
}
