package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type tokenType string

const (
	ACCESS  tokenType = "ACCESS"
	REFRESH tokenType = "REFRESH"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: []byte(secret)}
}

func (j *JWTService) GenerateToken(userID string, tokenType tokenType) (string, error) {
	var exp int64

	switch tokenType {
	case ACCESS:
		exp = time.Now().Add(15 * time.Minute).Unix()
	case REFRESH:
		exp = time.Now().Add(24 * time.Hour).Unix()
	default:
		return "", fmt.Errorf("invalid token type")
	}

	claims := jwt.MapClaims{
		"sub":  userID,
		"type": string(tokenType),
		"exp":  exp,
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTService) Validate(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
}
