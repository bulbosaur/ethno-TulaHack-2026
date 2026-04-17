package main

import (
	"database/sql"
	"ethno/db"
	"ethno/internal/config"
	srv "ethno/internal/transport/http"
	"net"

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
	server := initServer(db, cfg.Server, logger)
	httpAddr := net.JoinHostPort(cfg.HTTP.Host, (cfg.HTTP.Port))
	logger.Fatal(server.Start(httpAddr))
}


func initServer(db *sql.DB, cfg config.ServerConfig, logger *logrus.Logger) *srv.Server {
    server := srv.NewServer(&cfg, logger)
	return server
}
