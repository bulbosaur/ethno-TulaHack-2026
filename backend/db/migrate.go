package db

import (
	"fmt"
	"log"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("init migrate failed: %w", err)
	}
	defer m.Close()

	log.Println("Running migrations...")

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Database is already up to date.")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	version, _, _ := m.Version()
	log.Printf("Migrations applied successfully! Version: %d", version)
	return nil
}
