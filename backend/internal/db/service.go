package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Service struct {
	*DB
	Queries *Queries
}

type ServiceConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewService(cfg ServiceConfig) (*Service, error) {
	// Create database connection
	database, err := NewConnection(Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DBName:   cfg.DBName,
		SSLMode:  cfg.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create queries instance
	queries := New(database.Pool)

	return &Service{
		DB:      database,
		Queries: queries,
	}, nil
}

func (s *Service) RunMigrations(ctx context.Context) error {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Construct migrations path
	migrationsPath := filepath.Join(wd, "internal", "migrations")

	// Construct database URL
	sslmode := "disable"
	if s.DB.Pool.Config().ConnConfig.TLSConfig != nil {
		sslmode = "require"
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		s.DB.Pool.Config().ConnConfig.User,
		s.DB.Pool.Config().ConnConfig.Password,
		s.DB.Pool.Config().ConnConfig.Host,
		s.DB.Pool.Config().ConnConfig.Port,
		s.DB.Pool.Config().ConnConfig.Database,
		sslmode)

	// Run migrations
	cfg := MigrateConfig{
		MigrationsPath: migrationsPath,
		DatabaseURL:    databaseURL,
	}

	return RunMigrations(cfg)
}

func (s *Service) Close() {
	s.DB.Close()
}

// Health check
func (s *Service) HealthCheck(ctx context.Context) error {
	return s.DB.Ping(ctx)
}
