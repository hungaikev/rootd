package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hungaikev/rootd/backend/internal/db"
)

func main() {
	var (
		host     = flag.String("host", "localhost", "Database host")
		port     = flag.Int("port", 5432, "Database port")
		user     = flag.String("user", "postgres", "Database user")
		password = flag.String("password", "password", "Database password")
		dbname   = flag.String("dbname", "rootd", "Database name")
		sslmode  = flag.String("sslmode", "disable", "SSL mode")
		action   = flag.String("action", "up", "Migration action: up, down")
	)
	flag.Parse()

	cfg := db.ServiceConfig{
		Host:     *host,
		Port:     *port,
		User:     *user,
		Password: *password,
		DBName:   *dbname,
		SSLMode:  *sslmode,
	}

	service, err := db.NewService(cfg)
	if err != nil {
		log.Fatalf("Failed to create database service: %v", err)
	}
	defer service.Close()

	switch *action {
	case "up":
		if err := service.RunMigrations(nil); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		// Rollback logic would go here
		fmt.Println("Rollback not implemented yet")
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		os.Exit(1)
	}
}
