package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type Services interface {
	HashPassword(password string) (string, error)
	Verify(password, hashedPassword string) error
}

type services struct {
}

var _ Services = &services{}

func NewService() *services {
	return &services{}
}

func (s *services) HashPassword(password string) (string, error) {
	b := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *services) Verify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
