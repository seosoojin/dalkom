package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CtxKey string

const (
	CTXJWTKEY CtxKey = "JWTTOKEN"
)

type JWTService interface {
	GenerateToken(claims any) (string, error)
	VerifyToken(tokenString string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey []byte
}

var _ JWTService = &jwtService{}

func NewJWTService(secretKey []byte) *jwtService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (s *jwtService) GenerateToken(claims any) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": claims,
		"iss":  "dalkom",
		"exp":  time.Now().Add(time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *jwtService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
