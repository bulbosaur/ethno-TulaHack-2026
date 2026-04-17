package db

import (
	"database/sql"
	"ethno/internal/config"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
    dsn := cfg.DSN()
    
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    
    return db, nil
}
