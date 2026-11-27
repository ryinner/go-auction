package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const accessTokenTTL = 15 * time.Minute

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type TokenIssuer struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

func NewTokenIssuer(secret string) *TokenIssuer {
	return &TokenIssuer{
		secret: []byte(secret),
		ttl:    accessTokenTTL,
		now:    time.Now,
	}
}

func (t *TokenIssuer) Generate(id int, email string) (string, error) {
	var err error
	if id <= 0 {
		err = errors.Join(err, errors.New("user id must be positive"))
	}
	if email == "" {
		err = errors.Join(err, errors.New("email is required"))
	}

	if err != nil {
		return "", err
	}

	now := t.now()
	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(id),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(t.ttl)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(t.secret)
}
