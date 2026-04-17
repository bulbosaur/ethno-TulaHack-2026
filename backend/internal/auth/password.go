package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	SecretKey     string
	TokenDuration time.Duration
}

func NewService(secretKey string, tokenDuration time.Duration) *Service {
	return &Service{
		SecretKey:     secretKey,
		TokenDuration: tokenDuration * time.Hour,
	}
}

func (s *Service) GenerateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

func (s *Service) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
