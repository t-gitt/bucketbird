package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

// Generate creates a new JWT token for the given user ID
func (tm *TokenManager) Generate(userID uuid.UUID) (string, time.Time, error) {
	expires := time.Now().Add(tm.ttl)
	claims := Claims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(tm.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expires, nil
}

// Validate verifies a JWT token and returns the user ID
func (tm *TokenManager) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return tm.secret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}
	
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}
	
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}
	
	return userID, nil
}
