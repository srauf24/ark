# Go gardenjournal

A production-ready monorepo template for building scalable web applications with Go backend and TypeScript frontend. Built with modern best practices, clean architecture, and comprehensive tooling.

## Features

- **Monorepo Structure**: Organized with Turborepo for efficient builds and development
- **Go Backend**: High-performance REST API with Echo framework
- **Authentication**: Integrated Clerk SDK for secure user management
- **Database**: PostgreSQL with migrations and connection pooling
- **Background Jobs**: Redis-based async job processing with Asynq
- **Observability**: New Relic APM integration and structured logging
- **Email Service**: Transactional emails with Resend and HTML templates
- **Testing**: Comprehensive test infrastructure with Testcontainers
- **API Documentation**: OpenAPI/Swagger specification
- **Security**: Rate limiting, CORS, secure headers, and JWT validation

## Project Structure

```
gardenjournal/
├── apps/backend/          # Go backend application
├── packages/         # Frontend packages (React, Vue, etc.)
├── package.json      # Monorepo configuration
├── turbo.json        # Turborepo configuration
└── README.md         # This file
```

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Node.js 22+ and Bun
- PostgreSQL 16+
- Redis 8+

### Installation

1. Clone the repository:
```bash
git clone https://github.com/sriniously/gardenjournal.git
cd gardenjournal
```

2. Install dependencies:
```bash
# Install frontend dependencies
bun install

# Install backend dependencies
cd apps/backend
go mod download
```

3. Set up environment variables:
```bash
cp apps/backend/.env.example apps/backend/.env
# Edit apps/backend/.env with your configuration
```

4. Start PostgreSQL and Redis:
```bash
# Option 1: Using Docker Compose (recommended)
docker-compose up -d postgres redis

# Option 2: Using Docker directly
docker run -d --name gardenjournal-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gardenjournal \
  -p 5432:5432 \
  postgres:16

docker run -d --name gardenjournal-redis \
  -p 6379:6379 \
  redis:8-alpine

# Option 3: Using local installations
# PostgreSQL (if installed locally)
brew services start postgresql@16  # macOS
sudo systemctl start postgresql    # Linux

# Redis (if installed locally)
brew services start redis           # macOS
sudo systemctl start redis          # Linux
```

5. Run database migrations:
```bash
cd apps/backend
task migrations:up
```

6. Start the development server:
```bash
# From root directory
bun dev

# Or just the backend
cd apps/backend
task run
```

The API will be available at `http://localhost:8080`

## Development

### Available Commands

```bash
# Backend commands (from backend/ directory)
task help              # Show all available tasks
task run               # Run the application
task migrations:new    # Create a new migration
task migrations:up     # Apply migrations
task test              # Run tests
task tidy              # Format code and manage dependencies

# Frontend commands (from root directory)
bun dev                # Start development servers
bun build              # Build all packages
bun lint               # Lint all packages
```

### Environment Variables

The backend uses environment variables prefixed with `gardenjournal_`. Key variables include:

- `gardenjournal_DATABASE_*` - PostgreSQL connection settings
- `gardenjournal_SERVER_*` - Server configuration
- `gardenjournal_AUTH_*` - Authentication settings
- `gardenjournal_REDIS_*` - Redis connection
- `gardenjournal_EMAIL_*` - Email service configuration
- `gardenjournal_OBSERVABILITY_*` - Monitoring settings

See `apps/backend/.env.example` for a complete list.

## Architecture

This gardenjournal follows clean architecture principles:

- **Handlers**: HTTP request/response handling
- **Services**: Business logic implementation
- **Repositories**: Data access layer
- **Models**: Domain entities
- **Infrastructure**: External services (database, cache, email)

## Testing

```bash
# Run backend tests
cd apps/backend
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (requires Docker)
go test -tags=integration ./...
```

### Production Considerations

1. Use environment-specific configuration
2. Enable production logging levels
3. Configure proper database connection pooling
4. Set up monitoring and alerting
5. Use a reverse proxy (nginx, Caddy)
6. Enable rate limiting and security headers
7. Configure CORS for your domains

## Common Development Commands

### Database Migrations

```bash
# Navigate to backend directory
cd apps/backend

# Create a new migration
task migrations:new name=add_users_table

# Apply all pending migrations
task migrations:up

# Rollback the last migration
task migrations:down

# Check migration status
task migrations:status

# Rollback to a specific version
task migrations:goto version=001
```

### Redis Management

```bash
# Start Redis (Docker)
docker start gardenjournal-redis

# Stop Redis (Docker)
docker stop gardenjournal-redis

# View Redis logs
docker logs -f gardenjournal-redis

# Connect to Redis CLI
docker exec -it gardenjournal-redis redis-cli

# Or if Redis is installed locally
redis-cli

# Common Redis CLI commands
PING                    # Test connection
INFO                    # Server information
KEYS *                  # List all keys (use with caution in production)
FLUSHDB                 # Clear current database
MONITOR                 # Watch commands in real-time
CLIENT LIST             # List connected clients

# Check Redis connection from terminal
redis-cli ping          # Should return PONG

# Start Redis (local installation)
brew services start redis           # macOS
sudo systemctl start redis-server   # Linux

# Stop Redis (local installation)
brew services stop redis            # macOS
sudo systemctl stop redis-server    # Linux

# Check Redis status
brew services list                  # macOS
sudo systemctl status redis-server  # Linux
```

### OpenAPI Documentation

```bash
# Navigate to openapi package
cd packages/openapi

# Build the TypeScript contracts
bun run build

# Generate OpenAPI specification
bun run gen

# The generated spec will be at: packages/openapi/openapi.json
```

### Running the Application

```bash
# Run everything (backend + frontend) from root
bun dev

# Run backend only
cd apps/backend
task run

# Run backend with hot reload (if using air)
cd apps/backend
air

# Run frontend only
cd apps/frontend
bun dev
```

### Building for Production

```bash
# Build all packages from root
bun build

# Build backend binary
cd apps/backend
go build -o bin/gardenjournal ./cmd/gardenjournal

# Build with optimizations
go build -ldflags="-s -w" -o bin/gardenjournal ./cmd/gardenjournal

# Build Zod package
cd packages/zod
bun run build

# Build OpenAPI package
cd packages/openapi
bun run build
```

### Testing

```bash
# Run all Go tests
cd apps/backend
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/service/...

# Run tests with verbose output
go test -v ./...

# Run integration tests (requires Docker)
go test -tags=integration ./...

# Run frontend tests
cd apps/frontend
bun test
```

### Code Quality

```bash
# Format and tidy Go code
cd apps/backend
task tidy

# Run Go linter
golangci-lint run

# Format TypeScript/JavaScript
bun lint

# Fix linting issues
bun lint:fix
```

### Package Management

```bash
# Install new Go dependency
cd apps/backend
go get github.com/example/package

# Update Go dependencies
go get -u ./...
go mod tidy

# Install new package with bun
cd packages/zod
bun add <package-name>

# Update dependencies with bun
bun update
```

### Troubleshooting

```bash
# Clean Go build cache
go clean -cache

# Clean Go module cache
go clean -modcache

# Remove node_modules and reinstall
rm -rf node_modules bun.lockb
bun install

# Rebuild all TypeScript packages
cd packages/zod && bun run build
cd ../openapi && bun run build

# Reset database (WARNING: destroys all data)
cd apps/backend
task migrations:down  # Roll back all migrations
task migrations:up    # Reapply migrations

# Check Redis connection
redis-cli ping
# or for Docker
docker exec -it gardenjournal-redis redis-cli ping

# Clear Redis cache
redis-cli FLUSHDB
# or for Docker
docker exec -it gardenjournal-redis redis-cli FLUSHDB

# Check PostgreSQL connection
psql -h localhost -U postgres -d gardenjournal -c "SELECT 1;"
# or for Docker
docker exec -it gardenjournal-postgres psql -U postgres -d gardenjournal -c "SELECT 1;"

# View Redis and PostgreSQL container status
docker ps | grep gardenjournal

# Restart Redis
docker restart gardenjournal-redis
# or local
brew services restart redis  # macOS
sudo systemctl restart redis # Linux

# Restart PostgreSQL
docker restart gardenjournal-postgres
# or local
brew services restart postgresql@16  # macOS
sudo systemctl restart postgresql    # Linux
```

### Useful Development Tips

```bash
# Watch OpenAPI changes and regenerate
cd packages/openapi
bun run gen:watch  # If available

# Check Go module dependencies
cd apps/backend
go mod graph

# View available tasks
task help

# Check API routes
cd apps/backend
go run cmd/gardenjournal/main.go --routes  # If implemented
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
