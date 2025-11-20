# Ark Backend

A production-ready Go backend service for homelab asset tracking and configuration log management, built with Echo framework, PostgreSQL, and Clerk authentication. Features clean architecture, comprehensive testing, and structured logging.

## Quick Start

```bash
# 1. Clone and navigate to backend
cd apps/backend

# 2. Install dependencies
go mod download

# 3. Set up environment
cp .env.sample .env
# Edit .env with your Clerk credentials and database settings

# 4. Start PostgreSQL and Redis
# (Make sure they're running on localhost)

# 5. Run database migrations
task migrations:up

# 6. Start the server
task run
# Server will be available at http://localhost:8080
```

See [Getting Started](#getting-started) for detailed instructions.

## Architecture Overview

This backend follows clean architecture principles with clear separation of concerns:

```
backend/
├── cmd/ark/                  # Application entry point
├── internal/                  # Private application code
│   ├── config/               # Configuration management (Koanf-based)
│   ├── database/             # Database connections and migrations (pgx/v5)
│   ├── handler/              # HTTP request handlers (Echo)
│   │   ├── asset.go             # Asset CRUD operations
│   │   └── log.go                # Log CRUD operations
│   ├── service/              # Business logic layer
│   ├── repository/           # Data access layer (PostgreSQL)
│   ├── model/                # Domain models and DTOs
│   ├── middleware/           # HTTP middleware
│   │   ├── auth.go              # Two-phase authentication (Clerk)
│   │   ├── global.go            # Error handling, logging, CORS
│   │   └── middleware.go        # Middleware aggregator
│   ├── router/               # Route registration
│   │   └── v1/                  # API v1 routes
│   ├── lib/                  # Shared libraries
│   │   ├── jwt/                 # JWT verification helpers
│   │   └── errs/                # Custom error types
│   ├── validation/           # Request validation
│   └── server/               # Server configuration
├── tests/                    # Test suites
│   ├── integration/             # Integration tests (httptest)
│   └── manual/                  # Manual test scripts (.http files)
├── static/                   # Static files (OpenAPI spec)
├── templates/                # Email templates
└── Taskfile.yml              # Task automation
```

## Features

### Core Framework
- **Echo v4**: High-performance, minimalist web framework
- **Clean Architecture**: Handlers → Services → Repositories → Models
- **Dependency Injection**: Constructor-based DI for testability

### Database
- **PostgreSQL**: Primary database with pgx/v5 driver
- **Migration System**: Tern for schema versioning
- **Connection Pooling**: Optimized for production workloads
- **Transaction Support**: ACID compliance for critical operations

### Authentication & Security
- **Two-Phase Authentication Pattern**:
  - **Phase 1 - ClerkAuthMiddleware**: Global JWT verification on all `/api/v1/*` routes
  - **Phase 2 - RequireAuth**: User claims extraction and context setup per resource
- **Clerk SDK v2 Integration**: Modern authentication with RS256 JWT verification
- **Secure Token Validation**: Issuer verification, expiration checking, signature validation
- **User-Scoped Data Access**: All database queries automatically filtered by user_id
- **Custom Error Handling**: Structured errors without exposing internal details
- **Security Headers**: CORS, XSS protection, and secure defaults

### Observability
- **New Relic APM**: Application performance monitoring
- **Structured Logging**: JSON logs with Zerolog
- **Request Tracing**: Distributed tracing support
- **Health Checks**: Readiness and liveness endpoints
- **Custom Metrics**: Business-specific monitoring

### Background Jobs
- **Asynq**: Redis-based distributed task queue
- **Priority Queues**: Critical, default, and low priority
- **Job Scheduling**: Cron-like task scheduling
- **Retry Logic**: Exponential backoff for failed jobs
- **Job Monitoring**: Real-time job status tracking

### Email Service
- **Resend Integration**: Reliable email delivery
- **HTML Templates**: Beautiful transactional emails
- **Preview Mode**: Test emails in development
- **Batch Sending**: Efficient bulk operations

### API Documentation
- **OpenAPI 3.0**: Complete API specification
- **Swagger UI**: Interactive API explorer
- **Auto-generation**: Code-first approach

## Getting Started

### Prerequisites
- Go 1.24+
- PostgreSQL 16+
- Redis 8+
- Task (taskfile.dev)

### Installation

1. Install dependencies:
```bash
go mod download
```

2. Set up environment:
```bash
cp .env.sample .env
# Edit .env with your configuration:
# - Database credentials (PostgreSQL)
# - Clerk authentication keys (from Clerk Dashboard)
# - Redis connection (for background jobs)
# - Optional: New Relic license key, Resend API key
```

3. Run migrations:
```bash
task migrations:up
```

4. Start the server:
```bash
task run
```

## Configuration

Configuration is managed through environment variables with the `ARK_` prefix:

### Required Environment Variables

```bash
# Server Configuration
ARK_SERVER.PORT="8080"
ARK_SERVER.READ_TIMEOUT="30"
ARK_SERVER.WRITE_TIMEOUT="30"

# Database Configuration
ARK_DATABASE.HOST="localhost"
ARK_DATABASE.PORT="5432"
ARK_DATABASE.USER="your_db_user"
ARK_DATABASE.PASSWORD="your_db_password"
ARK_DATABASE.NAME="ark"
ARK_DATABASE.SSL_MODE="disable"  # Use "require" in production

# Clerk Authentication (Required)
ARK_AUTH.CLERK.SECRET_KEY="sk_test_..."  # From Clerk Dashboard → API Keys
ARK_AUTH.CLERK.JWT_ISSUER="https://your-app.clerk.accounts.dev"  # From Clerk Dashboard → JWT Templates

# Redis (for background jobs)
ARK_REDIS.ADDRESS="localhost:6379"

# Logging & Observability
ARK_OBSERVABILITY.LOGGING.LEVEL="debug"  # debug, info, warn, error
ARK_OBSERVABILITY.LOGGING.FORMAT="console"  # console or json

# Optional: New Relic APM
ARK_OBSERVABILITY.NEW_RELIC.LICENSE_KEY="your_license_key"
ARK_OBSERVABILITY.NEW_RELIC.APP_LOG_FORWARDING_ENABLED="true"

# Optional: Resend Email Service
ARK_INTEGRATION.RESEND_API_KEY="re_..."
```

See `.env.sample` for complete configuration options.

## Development

### Available Tasks

```bash
task help                    # Show all available tasks
task run                     # Run the application
task test                    # Run tests
task migrations:new name=X   # Create new migration
task migrations:up           # Apply migrations
task migrations:down         # Rollback last migration
task tidy                    # Format and tidy dependencies
```

### Project Structure

#### Handlers (`internal/handler/`)
HTTP request handlers that:
- Parse and validate requests
- Call appropriate services
- Format responses
- Handle HTTP-specific concerns

#### Services (`internal/service/`)
Business logic layer that:
- Implements use cases
- Orchestrates operations
- Enforces business rules
- Handles transactions

#### Repositories (`internal/repository/`)
Data access layer that:
- Encapsulates database queries
- Provides data mapping
- Handles database-specific logic
- Supports multiple data sources

#### Models (`internal/model/`)
Domain entities that:
- Define core business objects
- Include validation rules
- Remain database-agnostic

#### Middleware (`internal/middleware/`)
Cross-cutting concerns implemented as Echo middleware:

**Authentication** (`auth.go`):
- `ClerkAuthMiddleware`: Verifies JWT tokens using Clerk SDK, stores claims in context
- `RequireAuth`: Extracts user metadata (user_id, role, permissions) from claims

**Global Middleware** (`global.go`):
- `GlobalErrorHandler`: Centralized error handling with structured logging
- `RequestLogger`: Logs all requests with duration, status, and request_id
- `Recover`: Panic recovery with error logging
- `CORS`: Configurable cross-origin resource sharing

**Middleware Chain**:
```
Request → Global → ClerkAuth → RequireAuth → Validation → Handler → ErrorHandler → Response
```

### Testing

The project follows a comprehensive testing strategy with unit, integration, and manual tests.

#### Unit Tests
Run all unit tests:
```bash
task test
# or
go test ./...
```

**Test Coverage**:
- Configuration loading (`internal/config/config_test.go`)
- JWT verification helpers (`internal/lib/jwt/clerk_test.go`)
- Middleware behavior (`internal/middleware/auth_test.go`)
- Service layer business logic
- Repository data access patterns

**Testing Best Practices**:
- Use table-driven tests for multiple scenarios
- Mock external dependencies (database, Clerk SDK)
- Test both success and error paths
- Verify error messages and status codes
- Use `testify/assert` and `testify/require` for assertions

#### Integration Tests
Located in `tests/integration/`:
```bash
go test ./tests/integration/...
```

**Integration test features**:
- Full middleware chain testing with `httptest`
- Route-level authentication verification
- End-to-end request/response validation
- Tests all HTTP methods (GET, POST, PUT, DELETE)
- Verifies proper error handling and status codes

#### Manual Testing
Located in `tests/manual/`:

**REST Client Testing** (`.http` files):
- Use with Bruno, HTTPie, Postman, or VS Code REST Client
- Comprehensive authentication test scenarios
- Pre-configured requests for all endpoints
- Includes troubleshooting guides

Example:
```bash
# Using httpie
http GET http://localhost:8080/api/v1/assets \
  Authorization:"Bearer YOUR_JWT_TOKEN"

# Using curl
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/assets
```

**Getting a Fresh JWT Token**:
1. Open your Ark frontend in a browser
2. Open DevTools Console (F12)
3. Run: `await window.Clerk.session.getToken({ template: "api-test" })`
4. Copy the token for manual testing

## API Endpoints

All API endpoints require authentication via JWT token in the `Authorization` header.

### Available Endpoints

**Assets** (`/api/v1/assets`):
- `GET /api/v1/assets` - List all assets for authenticated user (paginated)
- `POST /api/v1/assets` - Create a new asset
- `GET /api/v1/assets/:id` - Get a specific asset by ID
- `PATCH /api/v1/assets/:id` - Update an asset
- `DELETE /api/v1/assets/:id` - Delete an asset

**Logs** (`/api/v1/logs`):
- `GET /api/v1/assets/:id/logs` - List all logs for a specific asset (paginated)
- `POST /api/v1/assets/:id/logs` - Create a new log for an asset
- `GET /api/v1/logs/:id` - Get a specific log by ID
- `PATCH /api/v1/logs/:id` - Update a log
- `DELETE /api/v1/logs/:id` - Delete a log

**Authentication**:
All requests must include:
```
Authorization: Bearer <your-jwt-token>
```

**Response Format**:
```json
{
  "data": [...],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 20
  }
}
```

**Error Response Format**:
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

See `static/openapi.json` for complete API specification.

## Troubleshooting

### Port Already in Use
```bash
# Error: listen tcp :8080: bind: address already in use
lsof -ti:8080 | xargs kill -9
task run
```

### Database Permission Denied
```bash
# Error: permission denied for table assets
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_db_user;"
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO your_db_user;"
```

### 401 Unauthorized Errors
**Common causes**:
1. Missing Authorization header
2. Invalid token format (must be `Bearer <token>`)
3. Expired JWT token
4. Wrong Clerk credentials in .env
5. JWT issuer mismatch

**Debugging steps**:
- Check server logs for detailed error messages
- Decode token at https://jwt.io to verify expiration
- Verify `iss` claim matches `ARK_AUTH.CLERK.JWT_ISSUER`
- Generate a fresh token from the frontend
- Confirm Clerk Secret Key is correct

### Database Connection Failed
```bash
# Test database connection
psql -U your_db_user -d ark -c "SELECT 1"

# Check if database exists
psql -U postgres -l | grep ark

# Create database if it doesn't exist
psql -U postgres -c "CREATE DATABASE ark;"
```

### Redis Connection Failed
```bash
# Check if Redis is running
redis-cli ping
# Expected output: PONG

# Start Redis (macOS with Homebrew)
brew services start redis

# Start Redis (Linux)
sudo systemctl start redis
```

## Logging

Structured logging with Zerolog:

```go
log.Info().
    Str("user_id", userID).
    Str("action", "login").
    Msg("User logged in successfully")
```

Log levels:
- `debug`: Detailed debugging information
- `info`: General informational messages
- `warn`: Warning messages
- `error`: Error messages
- `fatal`: Fatal errors that cause shutdown

### Production Checklist

- [ ] Set production environment variables
- [ ] Enable SSL/TLS
- [ ] Configure production database
- [ ] Set up monitoring alerts
- [ ] Configure log aggregation
- [ ] Enable rate limiting
- [ ] Set up backup strategy
- [ ] Configure auto-scaling
- [ ] Implement graceful shutdown
- [ ] Set up CI/CD pipeline

## Performance Optimization

### Database
- Connection pooling configured
- Prepared statements for frequent queries
- Indexes on commonly queried fields
- Query optimization with EXPLAIN ANALYZE

### Caching
- Redis for session storage
- In-memory caching for hot data
- HTTP caching headers

### Concurrency
- Goroutine pools for parallel processing
- Context-based cancellation
- Proper mutex usage

## Security Best Practices

1. **Input Validation**: All inputs validated and sanitized
2. **SQL Injection**: Parameterized queries only
3. **XSS Protection**: Output encoding and CSP headers
4. **CSRF Protection**: Token-based protection
5. **Rate Limiting**: Per-IP and per-user limits
6. **Secrets Management**: Environment variables, never in code
7. **HTTPS Only**: Enforce TLS in production
8. **Dependency Scanning**: Regular vulnerability checks

## Contributing

### Development Guidelines

1. **Follow Go Best Practices**
   - Use `gofmt` for code formatting
   - Follow effective Go guidelines
   - Use meaningful variable and function names
   - Add comments for exported functions

2. **Testing Requirements**
   - Write unit tests for new functions (`*_test.go`)
   - Add integration tests for new routes
   - Ensure all tests pass: `task test`
   - Aim for meaningful test coverage, not just high percentages

3. **Incremental Development**
   - Break large features into small, testable steps
   - Commit after each completed step
   - Include tests in the same commit as implementation
   - Write clear, descriptive commit messages

4. **Git Workflow**
   - Create feature branches: `git checkout -b feature/your-feature`
   - Make atomic commits with clear messages
   - Push regularly to remote branches
   - Use conventional commit format: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`

5. **Documentation**
   - Update README.md for new features
   - Update claude.md with architectural changes
   - Add comments to complex logic
   - Update OpenAPI spec for API changes

6. **Code Review**
   - Run `task tidy` before committing
   - Ensure no linter errors
   - Test manually with `.http` files
   - Verify error handling and logging

## License

See the parent project's LICENSE file.