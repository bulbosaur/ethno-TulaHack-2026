package main

import (
	"ethno/db"
	"ethno/internal/config"

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
    defer db.Close()

	logger.Info("Successful connection to database")
}