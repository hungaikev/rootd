.PHONY: help build run test clean migrate-up migrate-down sqlc-generate deps install-tools postgres createdb dropdb stop logs

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  sqlc-generate  - Generate Go code from SQL queries"
	@echo "  deps           - Install dependencies"
	@echo "  install-tools  - Install tools"
	@echo "  postgres       - Start a PostgreSQL container"
	@echo "  createdb       - Create the database"
	@echo "  dropdb         - Drop the database"
	@echo "  stop           - Stop the PostgreSQL container"
	@echo "  logs           - View PostgreSQL container logs"

# Build the application
build:
	go build -o bin/server ./cmd/server/*.go

# Run the application
run:
	go run ./cmd/server/*.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations
migrate-up:
	go run ./cmd/migrate/*.go -action=up

# Rollback database migrations
migrate-down:
	go run ./cmd/migrate/*.go -action=down

# Generate Go code from SQL queries
sqlc-generate:
	sqlc generate

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install tools
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Dockerized PostgreSQL
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rootd
DB_PORT=5432

postgres:
	docker run --name rootd-postgres -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d postgres

createdb:
	docker exec -it rootd-postgres createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it rootd-postgres dropdb $(DB_NAME)

stop:
	docker stop rootd-postgres

logs:
	docker logs -f rootd-postgres
