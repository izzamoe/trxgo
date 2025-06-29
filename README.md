# Transaction Management API

A RESTful API for transaction management with dashboard analytics built using clean architecture pattern.

## 🚀 Features

- ✅ **Precise Decimal Calculations** - Uses `decimal.Decimal` for accurate monetary operations
- ✅ **CRUD Operations** - Complete transaction management
- ✅ **Dashboard Analytics** - Real-time transaction summaries
- ✅ **Advanced Filtering** - User-based and status-based filtering
- ✅ **Clean Architecture** - Separation of concerns with dependency injection
- ✅ **Comprehensive Testing** - 78.4% overall test coverage
- ✅ **Migration Tools** - Database migration and setup tools
- ✅ **Production Ready** - Error handling, logging, and validation

## 🛠 Technology Stack

- **Language**: Go 1.24+
- **Framework**: Gin Web Framework
- **Database**: MySQL 8.0+ with GORM ORM
- **Testing**: Go testing + Testify + SQLite (for unit tests)
- **Logging**: Logrus structured logging
- **Decimal Precision**: shopspring/decimal for monetary calculations
- **Validation**: go-playground/validator with custom rules
- **Migration**: GORM auto-migration with custom tools

## 🏗 Architecture

Clean architecture with clear separation of concerns:

```
HTTP → Router → Handler → Service → Repository → MySQL
```

### Directory Structure

```
├── cmd/
│   ├── server/                        # Application entry point
│   ├── migrate/                       # Database migration tool
│   └── setup/                         # Database setup tool
├── internal/
│   ├── config/                        # Configuration management
│   ├── middleware/                    # HTTP middleware
│   ├── handlers/                      # HTTP handlers
│   ├── services/                      # Business logic
│   ├── repositories/                  # Database operations
│   └── models/                        # Data models
├── pkg/utils/                         # Utility packages
├── tests/                             # Test files
├── docs/                              # Documentation
├── Makefile                           # Build automation
├── .env.example                       # Environment variables template
└── .env                               # Environment variables (create from .env.example)
```

## 📋 Prerequisites

- Go 1.24 or higher
- MySQL 8.0 or higher
- Git

## ⚡ Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/izzamoe/trxgo.git
cd trxgo
```

### 2. Setup environment

Create `.env` file:
```bash
make env
```

Edit `.env` file with your database credentials:
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=interview_db

# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# Log Configuration
LOG_LEVEL=info
```

### 3. Complete automated setup

Run complete setup (install dependencies + database setup + migrations):
```bash
make setup
```

### 4. Start the server

```bash
make run
```

The server will start on `http://localhost:8080`

## 🛠️ Migration System

This project uses **GORM auto-migration** for database schema management with custom migration tools.

### Migration Tools

#### 1. Setup Tool (`cmd/setup/main.go`)
One-command database setup:
```bash
make db-setup       # Auto-creates database + runs migrations
./bin/setup         # Direct binary execution
```

#### 2. Migration Tool (`cmd/migrate/main.go`)
Fine-grained migration control:
```bash
make db-migrate     # Run migrations (up)
make db-migrate-down # Rollback migrations (down)  
make db-reset       # Reset database (drop + recreate)
make db-status      # Check migration status

# Direct binary usage
./bin/migrate -action=up
./bin/migrate -action=down  
./bin/migrate -action=reset
./bin/migrate -action=status -verbose
```

### Quick Commands

```bash
make setup          # Complete setup (recommended for first time)
make db-migrate     # Run migrations only
make db-status      # Check migration status  
make db-reset       # Reset database (drop + recreate)
```

### Manual Setup (Alternative)

If you prefer manual control:

```bash
# Install dependencies
make deps

# Create database manually
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS interview_db;"

# Run migrations
make db-migrate

# Start server
make run
```

## 📡 API Endpoints

### Transaction Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/transactions` | Create new transaction |
| GET | `/api/transactions` | Get all transactions (with filters) |
| GET | `/api/transactions/:id` | Get transaction by ID |
| PUT | `/api/transactions/:id` | Update transaction status |
| DELETE | `/api/transactions/:id` | Delete transaction |

### Dashboard

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/dashboard/summary` | Get dashboard analytics |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Server health status |

## 📊 Database Schema

GORM auto-migration creates this schema from models:

```sql
CREATE TABLE transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    status VARCHAR(191) NOT NULL DEFAULT 'pending',
    created_at DATETIME(3) DEFAULT NULL,
    updated_at DATETIME(3) DEFAULT NULL,
    
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_transactions_user_id (user_id),
    INDEX idx_transactions_status (status)
);
```

### GORM Model

```go
type Transaction struct {
    ID        uint            `json:"id" gorm:"primaryKey"`
    UserID    uint            `json:"user_id" gorm:"not null;index"`
    Amount    decimal.Decimal `json:"amount" gorm:"not null;type:decimal(15,2)"`
    Status    string          `json:"status" gorm:"not null;default:'pending';index"`
    CreatedAt time.Time       `json:"created_at"`
    UpdatedAt time.Time       `json:"updated_at"`
}
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
make test-coverage

# Run specific test packages
make test-cmd       # Test migration and setup tools
make test-internal  # Test business logic layers
go test ./internal/handlers
go test ./internal/services
go test ./internal/repositories
```

### Coverage Analysis

```bash
# Quick coverage summary
make coverage

# Generate HTML coverage report
make coverage-html

# Detailed coverage analysis
make coverage-detail
```

### Current Test Coverage: **78.4%**

| Package Category | Coverage | Status |
|------------------|----------|--------|
| **Core Business Logic** | 98-100% | ✅ Excellent |
| **HTTP Handlers** | 100% | ✅ Perfect |
| **Services** | 100% | ✅ Perfect |
| **Repositories** | 98% | ✅ Excellent |
| **Middleware** | 100% | ✅ Perfect |
| **Migration Tools** | 75% | ✅ Good |
| **Setup Tools** | 44% | ✅ Adequate |

### Test Coverage Details

The project maintains high test coverage across all critical components:

**✅ Production-Ready Coverage (78.4% Total)**
- **Business Logic**: 100% (handlers, services)
- **Database Layer**: 98% (repositories) 
- **Configuration**: 100% (config, middleware)
- **Migration Tools**: 75% (setup and migration tools)
- **Utilities**: 100% (response helpers)

**Test Types:**
- **Unit Tests**: Each layer tested in isolation
- **Integration Tests**: Database operations with real MySQL
- **Tool Tests**: Migration and setup functionality
- **Error Testing**: Comprehensive error scenario coverage

**Testing Philosophy:**
- Core business logic: 100% coverage requirement
- Tool functions: Adequate coverage with edge cases
- Error scenarios: All failure paths tested
- Database integration: Real database testing for reliability

## 📖 Usage Examples

### Create a transaction

```bash
curl -X POST http://localhost:8080/api/transactions \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "amount": 100.50}'
```

### Get transactions with filters

```bash
# Get all transactions
curl http://localhost:8080/api/transactions

# Filter by user ID
curl "http://localhost:8080/api/transactions?user_id=1"

# Filter by status
curl "http://localhost:8080/api/transactions?status=pending"

# Pagination
curl "http://localhost:8080/api/transactions?limit=10&offset=0"
```

### Update transaction status

```bash
curl -X PUT http://localhost:8080/api/transactions/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "success"}'
```

### Get dashboard summary

```bash
curl http://localhost:8080/api/dashboard/summary
```

## 🔧 Development

### Available Commands

```bash
# Setup and Database
make setup          # Complete setup (deps + DB + migrations)
make env            # Create .env template
make deps           # Install dependencies only

# Database Migration
make db-migrate     # Run migrations
make db-migrate-down # Rollback migrations  
make db-reset       # Reset database (drop + recreate)
make db-status      # Check migration status
make db-setup       # Complete database setup

# Build and Run
make build          # Build server binary
make run            # Build and run server
make dev            # Development mode with auto-reload

# Testing
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make test-cmd       # Run tests for cmd packages only
make test-internal  # Run tests for internal packages only

# Coverage Analysis
make coverage       # Quick coverage summary
make coverage-html  # Generate HTML coverage report
make coverage-detail # Detailed coverage with line counts

# Code Quality
make fmt            # Format Go code
make lint           # Run linter
make clean          # Clean build artifacts
```

### Project Structure Guidelines

1. **Models**: Define data structures and validation rules with GORM tags
2. **Repository**: Handle database operations only
3. **Service**: Implement business logic and orchestration
4. **Handler**: Handle HTTP requests/responses and validation
5. **Middleware**: Handle cross-cutting concerns (logging, CORS, recovery)
6. **Migration**: Use GORM auto-migration for schema changes

### Code Guidelines

- Follow Go best practices
- Use meaningful variable and function names
- Add comments for public functions
- Implement proper error handling
- Write tests for all layers
- Use dependency injection

### Adding New Features

1. Define models in `internal/models/` with proper GORM tags
2. Add model to migration list in `cmd/migrate/main.go`
3. Create repository interface and implementation
4. Implement service layer with business logic
5. Create HTTP handlers
6. Add routes in main.go
7. Run `make db-migrate` to update database schema
8. Write tests for all layers
9. Update API documentation

### Database Schema Changes

GORM auto-migration is **additive only**:
- ✅ Add new tables
- ✅ Add new columns  
- ✅ Add new indexes
- ❌ Drop columns (manual SQL required)
- ❌ Change column types (manual SQL required)

For destructive changes, add custom SQL in `cmd/migrate/main.go`

## 📝 Configuration

The application uses environment variables for configuration. All variables can be set in the `.env` file or as system environment variables.

### Available Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_USER` | Database username | `root` |
| `DB_PASSWORD` | Database password | `root` |
| `DB_NAME` | Database name | `interview_db` |
| `SERVER_HOST` | Server host | `localhost` |
| `SERVER_PORT` | Server port | `8080` |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |

## 🔍 Monitoring and Logging

The application includes comprehensive logging using Logrus:

- Request/response logging
- Error logging with stack traces
- Structured JSON logging
- Configurable log levels

Logs include:
- HTTP request details
- Database operation logs
- Error traces
- Performance metrics

## 🐳 Docker Support

The project includes a multi-stage Dockerfile optimized for production:

```dockerfile
# Multi-stage build for efficiency
FROM golang:1.24-alpine AS builder

# Install required dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates mysql-client

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy .env file if exists
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]
```

### Build and Run with Docker

Using Makefile commands:

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Or run in background
make docker-run-bg
```

Manual Docker commands:

```bash
# Build image
docker build -t transaction-api .

# Run container
docker run -p 8080:8080 --env-file .env transaction-api

# Run in background
docker run -d -p 8080:8080 --env-file .env --name transaction-api transaction-api
```

## 🚀 Deployment

### Production Checklist

- [ ] Set up production database
- [ ] Configure environment variables
- [ ] Set up reverse proxy (nginx)
- [ ] Configure SSL/TLS
- [ ] Set up monitoring and alerting
- [ ] Configure log aggregation
- [ ] Set up backup strategy
- [ ] Performance testing
- [ ] Security audit

### Environment-specific Configuration

Use different `.env` files for different environments:

- `.env.development`
- `.env.staging`
- `.env.production`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and ensure they pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 📚 Additional Resources

- [API Documentation](docs/api.md)
- [Gin Framework Documentation](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [Go Testing](https://golang.org/pkg/testing/)

## 🐛 Troubleshooting

### Common Issues

1. **Database connection failed**
   - Check MySQL server is running
   - Verify connection parameters in `.env`
   - Ensure database exists or run `make setup`

2. **Migration failed**
   - Check database permissions
   - Run `make db-status` to check current state
   - Try `make db-reset` to reset database

3. **Port already in use**
   - Change `SERVER_PORT` in `.env`
   - Kill process using the port: `lsof -ti:8080 | xargs kill`

4. **Build failed**
   - Run `make deps` to install dependencies
   - Check Go version (requires 1.21+)

5. **Test failures**
   - Run `make coverage` to check test coverage
   - Use `make test-cmd` to test migration tools
   - Check database connection for integration tests

6. **Coverage issues**
   - Run `make coverage-detail` for detailed analysis
   - Use `make coverage-html` to generate visual report

### Getting Help

- Check the logs for detailed error messages
- Review the API documentation  
- Run `make coverage` to check test coverage
- Use `make help` to see all available commands
- Create an issue in the repository
