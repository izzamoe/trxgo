# Makefile for Transaction Management API

# Variables
BINARY_NAME=server
MIGRATE_BINARY=migrate
SETUP_BINARY=setup
GO_BUILD_FLAGS=-v
DB_URL="mysql://root:password@localhost:3306/interview_db"

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Setup commands
.PHONY: setup
setup: deps build-setup db-setup ## Complete initial setup (install deps, build setup tool, setup database)
	@echo "✅ Setup completed successfully!"

.PHONY: setup-simple
setup-simple: deps build-migrate db-migrate ## Simple setup using migration tool only
	@echo "✅ Simple setup completed!"

.PHONY: deps
deps: ## Install Go dependencies
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Build commands
.PHONY: build
build: ## Build the server binary
	@echo "🏗️  Building server..."
	go build $(GO_BUILD_FLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

.PHONY: build-migrate
build-migrate: ## Build the migration binary
	@echo "🏗️  Building migration tool..."
	go build $(GO_BUILD_FLAGS) -o bin/$(MIGRATE_BINARY) cmd/migrate/main.go

.PHONY: build-setup
build-setup: ## Build the setup binary
	@echo "🏗️  Building setup tool..."
	go build $(GO_BUILD_FLAGS) -o bin/$(SETUP_BINARY) cmd/setup/main.go

# Database commands
.PHONY: db-migrate
db-migrate: build-migrate ## Run database migrations (up)
	@echo "🚀 Running database migrations..."
	./bin/$(MIGRATE_BINARY) -action=up

.PHONY: db-migrate-down
db-migrate-down: build-migrate ## Rollback database migrations
	@echo "📉 Rolling back database migrations..."
	./bin/$(MIGRATE_BINARY) -action=down

.PHONY: db-reset
db-reset: build-migrate ## Reset database (drop and recreate all tables)
	@echo "🔄 Resetting database..."
	./bin/$(MIGRATE_BINARY) -action=reset

.PHONY: db-status
db-status: build-migrate ## Check database migration status
	@echo "📊 Checking database status..."
	./bin/$(MIGRATE_BINARY) -action=status

.PHONY: db-setup
db-setup: build-setup ## Run complete database setup (create DB + migrate)
	@echo "🚀 Running complete database setup..."
	./bin/$(SETUP_BINARY)

# Development commands
.PHONY: run
run: build ## Build and run the server
	@echo "🚀 Starting server..."
	./bin/$(BINARY_NAME)

.PHONY: dev
dev: ## Run server in development mode with auto-reload
	@echo "🔧 Starting development server..."
	go run cmd/server/main.go

.PHONY: test
test: ## Run all tests
	@echo "🧪 Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

.PHONY: test-cmd
test-cmd: ## Run tests for cmd packages only
	@echo "🧪 Running cmd package tests..."
	go test -v ./cmd/...

.PHONY: test-internal
test-internal: ## Run tests for internal packages only
	@echo "🧪 Running internal package tests..."
	go test -v ./internal/...

.PHONY: coverage
coverage: ## Show coverage summary only
	@echo "📊 Running coverage analysis..."
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@echo "📈 Coverage Summary:"
	@echo "==================="
	@go tool cover -func=coverage.out | grep -E "(cmd|internal|pkg)" | awk '{printf "%-40s %s\n", $$1, $$3}' | sort
	@echo "==================="
	@go tool cover -func=coverage.out | tail -n 1 | awk '{printf "🎯 TOTAL COVERAGE: %s\n", $$3}'
	@echo "==================="

.PHONY: coverage-html
coverage-html: ## Generate HTML coverage report only
	@echo "📊 Generating HTML coverage report..."
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ HTML coverage report: coverage.html"

.PHONY: coverage-detail
coverage-detail: ## Show detailed coverage with line counts
	@echo "📊 Detailed Coverage Analysis..."
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@echo "📈 Detailed Coverage Report:"
	@echo "============================"
	@go tool cover -func=coverage.out

# Lint and format
.PHONY: fmt
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	go fmt ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "🔍 Running linter..."
	golangci-lint run

# Clean commands
.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean ## Clean everything including vendor and cache
	@echo "🧹 Deep cleaning..."
	go clean -cache -modcache

# Docker commands (optional)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t interview-api .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 --env-file .env interview-api

.PHONY: docker-run-bg
docker-run-bg: ## Run Docker container in background
	@echo "🐳 Running Docker container in background..."
	docker run -d -p 8080:8080 --env-file .env --name interview-api interview-api

# Quick start commands
.PHONY: start
start: setup run ## Quick start: setup and run the server

.PHONY: restart
restart: clean start ## Clean, setup and restart the server

# Environment setup
.PHONY: env
env: ## Create .env file template
	@echo "📝 Creating .env template..."
	@echo "# Database Configuration" > .env
	@echo "DB_HOST=localhost" >> .env
	@echo "DB_PORT=3306" >> .env
	@echo "DB_USER=root" >> .env
	@echo "DB_PASSWORD=password" >> .env
	@echo "DB_NAME=interview_db" >> .env
	@echo "" >> .env
	@echo "# Server Configuration" >> .env
	@echo "SERVER_HOST=localhost" >> .env
	@echo "SERVER_PORT=8080" >> .env
	@echo "" >> .env
	@echo "# Log Configuration" >> .env
	@echo "LOG_LEVEL=info" >> .env
	@echo "✅ .env file created. Please update with your values."

# Database setup helpers
.PHONY: db-create
db-create: ## Create database (requires MySQL running)
	@echo "🗄️  Creating database..."
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS interview_db;"

.PHONY: db-drop
db-drop: ## Drop database (requires MySQL running)
	@echo "🗑️  Dropping database..."
	mysql -u root -p -e "DROP DATABASE IF EXISTS interview_db;"

# Watch mode for development
.PHONY: watch
watch: ## Watch for changes and restart server (requires entr)
	@echo "👀 Watching for changes..."
	find . -name "*.go" | entr -r make dev
