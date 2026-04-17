package auth

import (
	"context"
	"errors"
	"ethno/internal/models"
	"ethno/internal/repository"
	"fmt"
)

type AuthService struct {
    userRepo *repository.UserRepository
    AuthProv Provider
}

func NewAuthService(userRepo *repository.UserRepository, authProv Provider) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        AuthProv: authProv,
    }
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
    existing, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, errors.New("user already exists")
    }

    hash, err := s.AuthProv.GenerateHash(req.Password)
    if err != nil {
        return nil, err
    }

    user, err := s.userRepo.Create(ctx, models.CreateUserDTO{
        Email:        req.Email,
        Username:     req.Username,
        PasswordHash: hash,
    })
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.User, string, error) {
    userAuth, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, "", err
    }
    if userAuth == nil {
        return nil, "", errors.New("invalid credentials")
    }

    if !s.AuthProv.Compare(userAuth.PasswordHash, req.Password) {
        return nil, "", errors.New("invalid credentials")
    }

    user := &models.User{
        ID:        userAuth.ID,
        Email:     userAuth.Email,
        Username:  userAuth.Username,
        Role:      "user",
    }

    token, err := s.AuthProv.GenerateJWT(user)
    if err != nil {
        return nil, "", fmt.Errorf("failed to generate token: %w", err)
    }

    return user, token, nil
}