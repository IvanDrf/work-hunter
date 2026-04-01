package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidJWT = errors.New("invalid jwt token")

type claims struct {
	UserID      string `json:"user_id"`
	Verificated bool   `json:"verificated"`

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

func (j *Jwt) CreateTokens(userID uuid.UUID, verificated bool) (string, string, error) {
	access, err := j.createToken(userID, verificated, j.accessTime)
	if err != nil {
		return "", "", err
	}

	refresh, err := j.createToken(userID, verificated, j.refreshTime)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (j *Jwt) createToken(userID uuid.UUID, verificated bool, duration time.Duration) (string, error) {
	claims := claims{
		UserID:      userID.String(),
		Verificated: verificated,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *Jwt) GetPayload(token string) (uuid.UUID, bool, error) {
	claims := &claims{}

	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidJWT
		}

		return j.secret, nil
	})

	if err != nil || !t.Valid {
		return uuid.UUID{}, false, err
	}

	id, err := uuid.Parse(claims.UserID)
	return id, claims.Verificated, err
}

func (j *Jwt) RefreshTokens(refresh string) (string, string, error) {
	userID, verificated, err := j.GetPayload(refresh)
	if err != nil {
		return "", "", err
	}

	return j.CreateTokens(userID, verificated)
}
