package repository

import (
	"context"
	"database/sql"
	"ethno/internal/models"
	"fmt"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, dto models.CreateUserDTO) (*models.User, error) {
    query := `
        INSERT INTO users (email, password_hash, username, role) 
        VALUES ($1, $2, $3, 'user') 
        RETURNING id, created_at
    `

    var user models.User
    err := r.db.QueryRowContext(ctx, query, dto.Email, dto.PasswordHash, dto.Username).
        Scan(&user.ID, &user.CreatedAt)
    
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    user.Email = dto.Email
    user.Username = dto.Username
    user.Role = "user"

    return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.UserAuth, error) {
    query := `
        SELECT id, email, password_hash, username 
        FROM users 
        WHERE email = $1
    `

    var u models.UserAuth
    err := r.db.QueryRowContext(ctx, query, email).
        Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Username)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }

    return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    query := `
        SELECT id, email, username, role, created_at 
        FROM users 
        WHERE id = $1
    `

    var u models.User
    err := r.db.QueryRowContext(ctx, query, id).
        Scan(&u.ID, &u.Email, &u.Username, &u.Role, &u.CreatedAt)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return &u, nil
}
