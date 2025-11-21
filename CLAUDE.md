# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Ark is a production-ready monorepo for homelab asset tracking and configuration log management with AI-powered search capabilities. The project combines a Go backend (Echo framework) with a TypeScript/React frontend, organized using Turborepo for efficient builds and development.

**Core Use Case**: Track homelab infrastructure (servers, VMs, containers, network equipment) and maintain searchable logs of configuration changes. An AI assistant helps query logs using natural language.

## Monorepo Structure

```
ark/
├── apps/
│   ├── backend/          # Go REST API (Echo framework)
│   └── frontend/         # React application (Vite + TypeScript)
├── packages/
│   ├── zod/              # Shared Zod schemas and types
│   ├── openapi/          # OpenAPI spec generation
│   └── emails/           # Email templates
├── package.json          # Workspace root
├── turbo.json            # Turborepo configuration
└── README.md
```

## Development Commands

### Monorepo-wide (from root directory)
```bash
# Install all dependencies
bun install

# Start all development servers (backend + frontend)
bun dev

# Build all packages
bun build

# Lint all packages
bun lint

# Fix linting issues
bun lint:fix

# Type check all TypeScript packages
bun typecheck

# Clean all build artifacts and node_modules
bun clean
```

### Backend (from apps/backend directory)
```bash
# Run the application
task run

# Run database migrations
task migrations:up

# Create a new migration
task migrations:new name=migration_name

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...

# Run a single test
go test -v ./internal/service -run TestAssetService_Create

# Format and tidy code
task tidy

# Show all available tasks
task help
```

### OpenAPI Documentation (from packages/openapi)
```bash
# Generate OpenAPI specification (outputs to apps/backend/static/openapi.json)
bun run gen

# Build TypeScript contracts
bun run build
```

## Backend Architecture

The backend follows **clean architecture** principles with clear separation of concerns:

### Layer Responsibilities

1. **Handlers** (`internal/handler/`): HTTP request/response handling
   - Parse requests, validate inputs
   - Delegate to service layer
   - Format responses
   - Return appropriate HTTP status codes

2. **Services** (`internal/service/`): Business logic implementation
   - Coordinate operations across repositories
   - Enforce business rules
   - Transaction management
   - Domain-specific validation

3. **Repositories** (`internal/repository/`): Data access layer
   - Database queries and mutations
   - Full-text search implementation
   - Data mapping to domain models
   - PostgreSQL-specific operations

4. **Models** (`internal/model/`): Domain entities and DTOs
   - Asset and AssetLog domain models
   - Request/response DTOs
   - Pagination structures
   - API response wrappers

5. **Middleware** (`internal/middleware/`): Request processing pipeline
   - Two-phase authentication (Clerk JWT)
   - Request ID generation
   - Distributed tracing (New Relic)
   - Global error handling
   - CORS configuration
   - Rate limiting

### Key Architectural Patterns

**Two-Phase Authentication**:
- Phase 1: `ClerkAuthMiddleware` - Applied globally to all `/api/v1/*` routes
  - Validates JWT from `Authorization: Bearer <token>`
  - Verifies with Clerk SDK v2
  - Stores session claims in Echo context
- Phase 2: `RequireAuth` - Applied per route group
  - Extracts user_id from verified claims
  - Sets user context for handlers
  - **All database queries MUST be scoped to user_id** (critical for multi-tenancy)

**Error Handling**:
- Custom error types in `internal/errs/`
  - `NotFoundError` → 404
  - `ValidationError` → 400
  - `AuthError` → 401
  - Generic errors → 500
- Global error handler in `internal/middleware/global.go`
  - Structured logging with request context
  - New Relic error tracking integration
  - Stack traces hidden in production

**Database Layer**:
- PostgreSQL 16+ with connection pooling (pgxpool)
- Full-text search using tsvector and GIN indexes
- Trigram matching for fuzzy search
- Automatic `updated_at` timestamps via triggers
- Migrations managed with tern

## Domain Models

### Asset
Represents homelab infrastructure (server, VM, container, NAS, network equipment)

| Field | Type | Notes |
|-------|------|-------|
| id | UUID | Primary key |
| user_id | string | Clerk user ID (multi-tenancy) |
| name | string | Required, max 100 chars |
| type | string? | server, vm, nas, container, network, other |
| hostname | string? | Max 255 chars |
| metadata | JSON | Flexible specs (CPU, RAM, IP, etc.) |
| created_at | timestamp | Auto-set |
| updated_at | timestamp | Auto-updated by trigger |

### AssetLog
Configuration changes or troubleshooting logs

| Field | Type | Notes |
|-------|------|-------|
| id | UUID | Primary key |
| asset_id | UUID | Foreign key (cascade delete) |
| user_id | string | Denormalized for performance |
| content | string | Required, 2-10,000 chars |
| tags | string[] | Optional, max 20 tags, 50 chars each |
| content_vector | tsvector | Generated column for FTS |
| created_at | timestamp | Auto-set |
| updated_at | timestamp | Auto-updated by trigger |

## API Endpoints

**Assets**:
```
GET    /api/v1/assets              # List (paginated)
POST   /api/v1/assets              # Create
GET    /api/v1/assets/:id          # Get single
PATCH  /api/v1/assets/:id          # Update
DELETE /api/v1/assets/:id          # Delete (cascades to logs)
```

**Logs** (nested for asset context):
```
POST   /api/v1/assets/:id/logs     # Create log for asset
GET    /api/v1/assets/:id/logs     # List logs for asset (paginated)
```

**Logs** (flat for direct access):
```
GET    /api/v1/logs/:id            # Get single log
PATCH  /api/v1/logs/:id            # Update log
DELETE /api/v1/logs/:id            # Delete log
```

**System**:
```
GET    /health                     # Health check (database, redis)
GET    /openapi.json               # OpenAPI specification
```

## Configuration

The backend uses environment variables prefixed with `ARK_`. Configuration is managed via Koanf with hierarchical loading (env vars override defaults).

**Critical Variables**:
```bash
# Server
ARK_SERVER.PORT="8080"
ARK_SERVER.CORS_ALLOWED_ORIGINS="http://localhost:3000"

# Database (PostgreSQL 16+)
ARK_DATABASE.HOST="localhost"
ARK_DATABASE.NAME="ark"
ARK_DATABASE.USER="ark_user"
ARK_DATABASE.PASSWORD="your_password"
ARK_DATABASE.MAX_OPEN_CONNS="25"

# Clerk Authentication
ARK_AUTH.CLERK.SECRET_KEY="sk_test_..."
ARK_AUTH.CLERK.JWT_ISSUER="https://your-app.clerk.accounts.dev"

# Redis (background jobs)
ARK_REDIS.ADDRESS="localhost:6379"

# Resend (transactional emails)
ARK_INTEGRATION.RESEND_API_KEY="re_..."

# OpenAI (AI features - planned)
ARK_OPENAI.API_KEY="sk-..."
ARK_OPENAI.MODEL="gpt-4o-mini"

# New Relic (optional but recommended)
ARK_OBSERVABILITY.NEW_RELIC.LICENSE_KEY="..."
ARK_OBSERVABILITY.LOGGING.LEVEL="debug"  # or "info" for production
```

See `apps/backend/.env.example` for complete configuration options.

## Production Infrastructure

### Observability Stack

**Structured Logging (Zerolog)**:
- Request-scoped context (request_id, trace_id, span_id, user_id)
- Configurable formats: console (dev) or JSON (production)
- Slow query logging with configurable threshold
- Integration with New Relic log forwarding

**New Relic APM Integration**:
- Distributed tracing across services
- Database query performance monitoring
- Custom transaction naming per endpoint
- Error tracking with stack traces
- Application log forwarding
- Performance dashboards

**Health Checks**:
- Endpoint: `GET /health`
- Validates database and Redis connectivity
- Used by load balancers and monitoring systems

### Background Job Processing

**Asynq (Redis-based job queue)**:
- Async task processing with retries (exponential backoff)
- Scheduled/delayed jobs
- Job prioritization and worker pools
- Current use cases: email sending, welcome emails
- Future: report generation, data exports, cleanup tasks

**Job Service** (`internal/lib/job/`):
- Centralized job enqueueing
- Worker lifecycle management
- Error handling and monitoring

### Email Integration

**Resend API** (`internal/lib/email/`):
- Transactional email sending
- HTML template rendering (templates in `templates/emails/`)
- Email tracking and high deliverability
- Current templates: welcome email
- Future: password reset, notifications, digests

## Testing Strategy

### Integration Tests (`tests/integration/`)
- Full middleware chain testing
- Auth verification and error scenarios
- OpenAPI spec validation
- Uses `httptest` for request/response simulation
- New Relic transaction tracking validation

### Unit Tests
Run unit tests for specific packages:
```bash
# Test specific service
go test ./internal/service/...

# Test with verbose output
go test -v ./internal/handler/...

# Test with coverage report
go test -cover ./internal/repository/...
```

**Coverage expectations**:
- Handlers: request parsing, validation, error responses
- Services: business logic, transaction handling
- Repositories: data access, query construction
- Middleware: auth flow, error handling, context management
- Models: validation, serialization

### Manual Tests (`tests/manual/`)
- `test_auth.http` - Authentication testing
- HTTP files for asset and log CRUD operations
- Can be used with REST client extensions in IDEs

## Database Management

### Migrations
```bash
# Navigate to backend directory
cd apps/backend

# Create a new migration
task migrations:new name=add_users_table

# Apply all pending migrations (requires confirmation)
task migrations:up

# Rollback the last migration
task migrations:down

# Check migration status
task migrations:status
```

**Migration Guidelines**:
- Always create reversible migrations (both up and down)
- Include data migrations when schema changes affect existing data
- Test migrations on a copy of production data before deployment
- Migrations are located in `internal/database/migrations/`

### Database Features
- **Full-text search**: `content_vector` tsvector with GIN index on asset_logs
- **Trigram matching**: Fuzzy search on asset names
- **Indexes**: user_id, asset_id, created_at, tags (GIN), content_vector (GIN)
- **Triggers**: Auto-update `updated_at` on row changes
- **Multi-tenancy**: All queries scoped to user_id

## Common Development Tasks

### Adding a new API endpoint
1. Define Zod schema in `packages/zod/` for request/response validation
2. Add handler method in `internal/handler/`
3. Implement business logic in `internal/service/`
4. Add data access in `internal/repository/` if needed
5. Register route in `internal/router/v1/v1.go`
6. Update OpenAPI contracts and regenerate spec
7. Write unit tests for handler, service, and repository
8. Add integration test in `tests/integration/`

### API Consistency Guidelines
**Critical**: Ensure backend route registration matches OpenAPI contracts exactly.
- **Prefix Consistency**:
  - If registering on a group (e.g., `v1 := api.Group("/v1")`), the contract path should be relative (e.g., `/assets`).
  - If registering on the root Echo instance (e.g., `e.GET(...)`), the contract path must include the full path (e.g., `/api/status`).
- **Verification**:
  - Always run `bun gen` after changing contracts.
  - Run `go test -v ./internal/handler -run TestOpenAPI` to verify the spec matches expected endpoints.

### Working with the OpenAPI spec
1. Update Zod schemas in `packages/zod/`
2. Update contracts in `packages/openapi/src/contracts/`
3. Run `bun run gen` from `packages/openapi/` to regenerate spec
4. Generated spec is written to `apps/backend/static/openapi.json`
5. Verify spec at `http://localhost:8080/openapi.json` when backend is running

### Running tests during development
```bash
# Run tests on file save (use with entr or similar)
ls **/*.go | entr go test ./...

# Run specific test with verbose output
go test -v ./internal/service -run TestAssetService_Create

# Run tests excluding integration tests
go test ./... -short
```

## Security Considerations

**Multi-tenancy enforcement**:
- All database queries MUST filter by `user_id`
- Asset ownership verified before operations
- Log operations verify asset ownership transitively

**Authentication**:
- JWT validation on every `/api/v1/*` request
- Clerk SDK v2 handles token verification
- Session claims stored in Echo context
- No manual token parsing

**Input validation**:
- Zod schemas enforce max lengths, required fields
- SQL injection protection via parameterized queries
- XSS protection via Echo's built-in sanitization
- CORS configured for allowed origins only

**Secrets management**:
- Never commit `.env` files (in `.gitignore`)
- Use environment variables for all secrets
- API keys, JWT secrets, database passwords in env vars only

## Troubleshooting

### Port already in use
```bash
lsof -ti:8080 | xargs kill -9
```

### Database connection issues
```bash
# Test PostgreSQL connection
psql -h localhost -U postgres -d ark -c "SELECT 1;"

# Grant permissions if needed
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ark_user;"
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ark_user;"
```

### Redis connection issues
```bash
# Test Redis connection
redis-cli ping  # Should return PONG

# Start Redis (macOS)
brew services start redis

# Start Redis (Linux)
sudo systemctl start redis-server
```

### JWT authentication issues
- Verify token at https://jwt.io
- Check `iss` claim matches `ARK_AUTH.CLERK.JWT_ISSUER`
- Ensure Clerk secret key is valid
- Generate fresh token from browser console: `await window.Clerk.session.getToken({ template: "api-test" })`

### Clean build issues
```bash
# Clean Go build cache
go clean -cache
go clean -modcache

# Rebuild frontend packages
cd packages/zod && bun run build
cd ../openapi && bun run build

# Clean and reinstall dependencies
rm -rf node_modules bun.lockb
bun install
```

## Migration History

**Note**: This project was migrated from `garden_journal` codebase. All plant/observation domain code has been replaced with asset/log equivalents. Module name is now `ark`, database name is `ark`.

**What was reused** (production-ready infrastructure):
- Authentication system (Clerk middleware, JWT verification)
- Observability stack (New Relic, Zerolog, health checks)
- Background job processing (Asynq, Redis)
- Email system (Resend, templates)
- Error handling and middleware pipeline
- Configuration management (Koanf)

**What was adapted**:
- Domain models (Asset, AssetLog instead of Plant, Observation)
- Handlers, services, repositories for new domain
- API endpoints and routes
- OpenAPI specifications

## Prerequisites

- **Go**: 1.24 or higher
- **Node.js**: 22+ (for frontend and package builds)
- **Bun**: Latest version (package manager)
- **PostgreSQL**: 16+ (local installation required)
- **Redis**: 8+ (local installation required)
- **Task**: Task runner for Go backend (install via `brew install go-task`)
- **Tern**: Database migration tool (installed via Go modules)

**Note**: Docker and Docker Compose are not yet configured. You must install PostgreSQL and Redis locally.
