package models

import (
	"time"
)

type User struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type UserAuth struct {
    ID           string
    Email        string
    PasswordHash string
    Username     string
}

type CreateUserDTO struct {
    Email        string
    PasswordHash string
    Username     string
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Username string `json:"username"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
