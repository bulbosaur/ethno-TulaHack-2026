package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Environment    string         `env:"ENV"`
	Database       DatabaseConfig `envPrefix:"DB_"`
	HTTP           HTTPConfig     `envPrefix:"HTTP_"`
	Server         ServerConfig
	MigrationsPath string `env:"MIGRATIONS_PATH"`
}

type DatabaseConfig struct {
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT,required"`
	Name     string `env:"NAME,required"`
}

type JWTConfig struct {
	SecretKey   string `env:"SECRET,required"`
	ExpiryHours int    `env:"EXPIRY_HOURS,required"`
}

type CookieConfig struct {
	Domain string
	Secure bool
}

type ServerConfig struct {
	JWT    JWTConfig    `envPrefix:"JWT_"`
	Cookie CookieConfig `envPrefix:"COOKIE_"`
}

type HTTPConfig struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (db DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.User, db.Password, db.Host, db.Port, db.Name,
	)
}


func Load() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        if !os.IsNotExist(err) {
            return nil, fmt.Errorf("failed to load .env: %w", err)
        }
    }

    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    return cfg, nil
}

func (c *Config) LogFields() logrus.Fields {
    return logrus.Fields{
        "env":             c.Environment,
        "db_user":         c.Database.User,
        "db_host":         c.Database.Host,
        "db_port":         c.Database.Port,
        "db_name":         c.Database.Name,
		"http_host": c.HTTP.Host,
		"http_port": c.HTTP.Port,
        "migrations_path": c.MigrationsPath,
    }
}
