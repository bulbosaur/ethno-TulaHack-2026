package main

import (
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
}