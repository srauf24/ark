# Ark

A production-ready monorepo for homelab asset tracking and configuration log management with AI-powered search. Built with Go backend and TypeScript frontend using modern best practices and clean architecture.

## Features

- **Monorepo Structure**: Organized with Turborepo for efficient builds and development
- **Go Backend**: High-performance REST API with Echo framework for asset and log management
- **AI-Powered Search**: Natural language querying of configuration logs (planned)
- **Authentication**: Integrated Clerk SDK for secure user management
- **Database**: PostgreSQL with full-text search, migrations, and connection pooling
- **Background Jobs**: Redis-based async job processing with Asynq
- **Observability**: New Relic APM integration and structured logging
- **Email Service**: Transactional emails with Resend and HTML templates
- **Testing**: Comprehensive test infrastructure
- **API Documentation**: OpenAPI/Swagger specification
- **Security**: Rate limiting, CORS, secure headers, and JWT validation

## Project Structure

```
ark/
├── apps/
│   ├── backend/          # Go backend application
│   └── frontend/         # React frontend (in progress)
├── packages/             # Shared packages
├── package.json          # Monorepo configuration
├── turbo.json            # Turborepo configuration
└── README.md             # This file
```
## OpenApi Documentation 
<img width="1512" height="834" alt="image" src="https://github.com/user-attachments/assets/b95d92ab-72c9-47ce-943e-7de22f90492b" />


## Quick Start

### Prerequisites

- Go 1.24 or higher
- Node.js 22+ and Bun
- PostgreSQL 16+
- Redis 8+

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ark.git
cd ark
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
# PostgreSQL (local installation required)
brew services start postgresql@16  # macOS
sudo systemctl start postgresql    # Linux

# Redis (local installation required)
brew services start redis           # macOS
sudo systemctl start redis          # Linux
```

> [!NOTE]
> Docker and Docker Compose are not yet configured for this project. You'll need to install PostgreSQL 16+ and Redis 8+ locally.

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

The backend uses environment variables prefixed with `ARK_`. Key variables include:

- `ARK_DATABASE_*` - PostgreSQL connection settings
- `ARK_SERVER_*` - Server configuration
- `ARK_AUTH_*` - Authentication settings (Clerk)
- `ARK_REDIS_*` - Redis connection
- `ARK_INTEGRATION_*` - Email service configuration (Resend)
- `ARK_OBSERVABILITY_*` - Monitoring settings (New Relic)
- `ARK_OPENAI_*` - OpenAI API configuration (for AI features)

See `apps/backend/.env.example` for a complete list.

## Architecture

This project follows clean architecture principles:

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

# Run integration tests
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
# Connect to Redis CLI (local installation)
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
go build -o bin/ark ./cmd/ark

# Build with optimizations
go build -ldflags="-s -w" -o bin/ark ./cmd/ark
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

# Clear Redis cache
redis-cli FLUSHDB

# Check PostgreSQL connection
psql -h localhost -U postgres -d ark -c "SELECT 1;"

# Restart Redis (local)
brew services restart redis  # macOS
sudo systemctl restart redis # Linux

# Restart PostgreSQL (local)
brew services restart postgresql@16  # macOS
sudo systemctl restart postgresql    # Linux
```

### Useful Development Tips

```bash
# Check Go module dependencies
cd apps/backend
go mod graph

# View available tasks
task help
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
