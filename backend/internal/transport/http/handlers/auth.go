package handler

import (
	"ethno/internal/auth"
	"ethno/internal/config"

	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
    authService *auth.AuthService
    config *config.ServerConfig
    logger *logrus.Logger
}

func NewAuthHandler(authService *auth.AuthService, cfg *config.ServerConfig) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        config: cfg,
        logger: logrus.New(),
    }
}
