package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("Token in Invalid")
	ErrExpiredToken = errors.New("Token Is Expired")
)

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	payload := Payload{
		ID:       tokenId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}
	return &payload, nil
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, fmt.Errorf("Invalid Payload : %w", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
