package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseType int

const (
	_ DatabaseType = iota
	POSTGRES
	SQLITE
)

type Config struct {
	DatabaseType DatabaseType
	SqlitePath   string
	PostgresDsn  string
	ApiPort      string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	cfg := Config{}
	dt := os.Getenv("DATABASE_TYPE")
	switch dt {
	case "POSTGRES":
		cfg.DatabaseType = POSTGRES
		dsn := os.Getenv("POSTGRES_DSN")
		if dsn == "" {
			return nil, fmt.Errorf("missing POSTGRES_DSN in env")
		}
		cfg.PostgresDsn = dsn
	case "SQLITE":
		cfg.DatabaseType = SQLITE
		path := os.Getenv("SQLITE_DATABASE_PATH")
		if path == "" {
			return nil, fmt.Errorf("missing SQLITE_DATABASE_PATH in env")
		}
		cfg.SqlitePath = path
	default:
		return nil, fmt.Errorf("invalid database type; must be 'POSTGRES' or 'SQLITE'")
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		return nil, fmt.Errorf("missing API_PORT in env")
	}
	cfg.ApiPort = port
	return &cfg, nil
}
