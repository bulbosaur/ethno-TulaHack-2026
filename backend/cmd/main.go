package main

import (
	"database/sql"
	"ethno/db"
	"ethno/internal/auth"
	"ethno/internal/config"
	"ethno/internal/repository"
	srv "ethno/internal/transport/http"
	"log"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	logger.WithFields(cfg.LogFields()).Info("Successful loading config")

	err = db.RunMigrations(cfg.Database.DSN())
	if err != nil {
		logger.Fatalf("Failed to make migrations: %v", err)
	}

	db, err := db.Connect(cfg.Database)
    if err != nil {
        logger.Fatalf("Failed to connect to database: %v", err)
    }

	_, err = db.Exec("SET NAMES 'UTF8'")
	if err != nil {
		log.Printf("Warning: could not set client encoding: %v", err)
	}
    defer db.Close()

	logger.Info("Successful connection to database")
	server := initServer(db, cfg.Server, logger)
	httpAddr := net.JoinHostPort(cfg.HTTP.Host, (cfg.HTTP.Port))
	logger.Fatal(server.Start(httpAddr))
}


func initServer(db *sql.DB, cfg config.ServerConfig, logger *logrus.Logger) *srv.Server {
	authProv := auth.NewService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.ExpiryHours)*time.Hour,
	)
	userRepo := repository.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo, authProv)
	folkRepo := repository.NewFolkRepository(db)
	questRepo := repository.NewQuestRepository(db)
    server := srv.NewServer(folkRepo, authService, &cfg, logger, questRepo)
	return server
}
