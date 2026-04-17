package auth

import (
	"errors"
	"ethno/internal/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   string   `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
type Provider interface {
	GenerateHash(password string) (string, error)
	Compare(hash, password string) bool
	GenerateJWT(user *models.User) (string, error)
	ParseJWT(tokenString string) (*Claims, error)
}

func (s *Service) GenerateJWT(user *models.User) (string, error) {
    if s.SecretKey == "" {
        return "", errors.New("secret key is empty")
    }
    
    claims := &Claims{
        UserID:   fmt.Sprintf("%v", user.ID),
        Email:    user.Email,
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.TokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "hack-auth-service",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.SecretKey))
}

func (s *Service) ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
