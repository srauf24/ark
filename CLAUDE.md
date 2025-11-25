# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Ark Project Guide

## Project Overview
**Ark** is a self-hosted asset management system for homelabs. It allows users to track their hardware (servers, VMs, containers) and maintain a history of changes and logs.

## Core Features
- **Asset Management**: Track servers, VMs, containers, and network devices.
- **Log History**: Maintain a chronological log of changes and maintenance tasks.
- **Self-Hosted**: Designed to run in a homelab environment.

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

### Frontend (from apps/frontend directory)
```bash
# Start development server
bun dev

# Build for production
bun build

# Testing
bun test              # Run unit/component tests (Vitest)
bun test:e2e          # Run E2E tests (Playwright, headless)
bun test:e2e:ui       # Run E2E tests (Playwright UI mode)
bun test:e2e:debug    # Run E2E tests (Playwright debug mode)
bun test:e2e:report   # View E2E test report

# Linting and formatting
bun lint              # Check for linting issues
bun lint:fix          # Auto-fix linting issues
bun format            # Check formatting
bun format:fix        # Auto-fix formatting

# Type checking
bun typecheck

# Clean build artifacts
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

### Shared Packages

#### OpenAPI Package (from packages/openapi)
```bash
# Generate OpenAPI specification (outputs to apps/backend/static/openapi.json)
bun run gen

# Build TypeScript contracts
bun run build
```

#### Zod Package (from packages/zod)
```bash
# Build shared Zod schemas and TypeScript types
bun run build
```

**Important**: The frontend dev server waits for OpenAPI contracts to build before starting (`wait-on ../../packages/openapi/dist/index.js`). If the frontend won't start, ensure packages are built first.

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
- **CRITICAL**: Clerk SDK must be initialized with `clerk.SetKey(cfg.Auth.Clerk.SecretKey)` in `internal/server/server.go` during server startup. Without this, JWT verification will fail with 401 errors.
- Phase 1: `ClerkAuthMiddleware` - Applied globally to all `/api/v1/*` routes
  - Validates JWT from `Authorization: Bearer <token>`
  - Verifies with Clerk SDK v2
  - Stores session claims in Echo context using key `clerk_session_claims`
- Phase 2: `RequireAuth` - Applied globally to all `/api/v1/*` routes (MUST be applied after ClerkAuthMiddleware)
  - Retrieves verified claims from context
  - Extracts user_id from claims.Subject
  - Sets `user_id` in context for handlers
  - **All database queries MUST be scoped to user_id** (critical for multi-tenancy)
- **Common Issue**: If Phase 2 is missing from router configuration, handlers will fail with "unauthorized: user not authenticated" even though tokens are valid.

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

## Frontend Architecture

### Tech Stack
- **React 19**: Latest React with concurrent features
- **TypeScript**: Strict mode enabled for type safety
- **Vite 7**: Ultra-fast build tool and dev server
- **TailwindCSS v4**: Utility-first CSS with Vite plugin
- **Clerk**: Authentication provider (`@clerk/clerk-react`)
- **TanStack Query**: Server state management and data fetching
- **React Router v7**: Client-side routing
- **ts-rest**: Type-safe API client from OpenAPI contracts
- **shadcn/ui**: Accessible component library (New York style, Lucide icons)
- **Playwright**: E2E testing framework

### Environment Variables (Frontend)

**Critical**: Vite has special handling for environment variables:

1. **Only variables prefixed with `VITE_` are exposed to the browser**
   - This prevents accidental exposure of server secrets
   - Example: `VITE_API_URL`, `VITE_CLERK_PUBLISHABLE_KEY`

2. **Access via `import.meta.env`, NOT `process.env`**
   - `process.env` is Node.js only and will cause `ReferenceError` in browser
   - Use `import.meta.env.VITE_API_URL` in frontend code

3. **Environment files loaded by Vite**:
   - `.env.local` (loaded in development, gitignored)
   - `.env.development` (loaded in development)
   - `.env.production` (loaded in production)
   - Plain `.env` files are NOT automatically loaded in dev mode

4. **Always restart `bun dev` after changing `.env` files**

**Required Frontend Variables** (in `.env.local`):
```bash
# Clerk Authentication (PUBLISHABLE key, not secret!)
VITE_CLERK_PUBLISHABLE_KEY=pk_test_...

# Backend API URL
VITE_API_URL=http://localhost:8080

# Environment
VITE_ENV=local  # or "development" | "production"
```

**Common Mistake**: Using Clerk Secret Key (`sk_test_...`) in frontend instead of Publishable Key (`pk_test_...`). The secret key is BACKEND ONLY.

### Type-Safe API Client

The frontend uses ts-rest to create a fully type-safe API client from backend contracts:

```typescript
import { useApiClient } from "@/api";

function MyComponent() {
  const apiClient = useApiClient();

  // Fully typed API call - response types inferred from backend
  const response = await apiClient.assets.getAll({
    query: { page: 1, limit: 10 }
  });

  if (response.status === 200) {
    // response.body.data is typed as Asset[]
    const assets = response.body.data;
  }
}
```

**Features**:
- Automatic JWT injection from Clerk (`Authorization: Bearer <token>`)
- Custom JWT template named **"api-test"** (configured in Clerk dashboard)
- Retry logic for 401 errors (up to 2 retries for token refresh)
- Support for blob responses (file downloads)

### Implemented Features

**Asset List View** (`/assets`):
- Component: `AssetList` with `AssetCard` children
- Displays paginated grid of user's assets
- Features:
  - Dynamic icons based on asset type (Server, HardDrive, Container, Network, Box)
  - Last updated timestamp (formatted with `date-fns`)
  - Loading, error, and empty states
  - Click to navigate to detail view
- Data fetching: TanStack Query with `useApiClient`
- Tests: Full coverage in `AssetList.test.tsx` and `AssetCard.test.tsx`

**Asset Detail View** (`/assets/:id`):
- Component: `AssetDetailPage`
- Displays full asset information:
  - Name, type, hostname
  - Formatted JSON metadata viewer
  - Created/updated timestamps
  - Back navigation to list
- Placeholder sections for future features (Logs, Actions)
- Tests: Full coverage in `AssetDetailPage.test.tsx`

**Testing Notes**:
- All tests use `happy-dom` environment (specified via `// @vitest-environment happy-dom`)
- Run with: `TZ=UTC VITE_CLERK_PUBLISHABLE_KEY=pk_test_mock bun x vitest run`
- Timezone must be UTC to ensure consistent date formatting across environments
- Full type safety from backend contracts

### Frontend Testing

#### Unit/Component Tests (Vitest)
```bash
bun test              # Run all tests
bun test --watch      # Watch mode
bun test --coverage   # Coverage report
```

**Testing setup**:
- Vitest as test runner (Vite-native)
- @testing-library/react for component testing
- happy-dom for DOM simulation
- Mock API calls using TanStack Query testing utilities

#### E2E Tests (Playwright)
```bash
bun test:e2e          # Headless mode (CI)
bun test:e2e:ui       # Interactive UI mode
bun test:e2e:debug    # Step-through debugging
bun test:e2e:report   # View HTML report
```

**Playwright configuration** (`playwright.config.ts`):
- Tests in `e2e/` directory
- Runs on Chromium, Firefox, and WebKit
- Automatic dev server startup (port 3000)
- Parallel execution enabled
- Screenshots on failure
- Traces on first retry

**Writing E2E tests**:
```typescript
import { test, expect } from "@playwright/test";

test("should authenticate and view assets", async ({ page }) => {
  await page.goto("/");

  // Use semantic selectors (role, label, text)
  await page.getByRole("button", { name: "Sign in" }).click();

  // Assertions
  await expect(page.getByRole("heading", { name: "Assets" })).toBeVisible();
});
```

## Troubleshooting

### Authentication 401 Errors

**Symptom**: Frontend shows "Failed to load assets" or backend logs show "unauthorized: user not authenticated" even though user is logged in.

**Root Causes and Fixes**:

1. **Clerk SDK not initialized** (Backend)
   - **Check**: Look for `clerk.SetKey(cfg.Auth.Clerk.SecretKey)` in `apps/backend/internal/server/server.go`
   - **Fix**: Add initialization in `server.New()` function before any middleware setup
   - **Verify**: Backend logs should show "token verification successful" when requests arrive

2. **RequireAuth middleware missing** (Backend)
   - **Check**: Verify `internal/router/v1/v1.go` has BOTH middlewares:
     ```go
     v1.Use(m.Auth.ClerkAuthMiddleware)  // Phase 1: Verify token
     v1.Use(m.Auth.RequireAuth)          // Phase 2: Set user_id
     ```
   - **Symptom**: Logs show "JWT verified and claims stored" but then "user not authenticated"
   - **Fix**: Add `v1.Use(m.Auth.RequireAuth)` after ClerkAuthMiddleware

3. **JWT template mismatch** (Frontend/Backend)
   - **Check**: Frontend uses template "api-test" in `apps/frontend/src/api/index.ts`
   - **Check**: Clerk dashboard has JWT template named "api-test" configured
   - **Fix**: Create template in Clerk dashboard or update frontend to match existing template name

4. **Clerk configuration mismatch** (Backend)
   - **Check**: `ARK_AUTH.CLERK.JWT_ISSUER` matches your Clerk instance
   - **Example**: `https://ace-dinosaur-39.clerk.accounts.dev`
   - **Verify**: Decode JWT token (jwt.io) and check `iss` claim matches backend config

5. **Publishable key mismatch** (Frontend)
   - **Check**: `VITE_CLERK_PUBLISHABLE_KEY` in `.env.local` matches backend's Clerk instance
   - **Example**: `pk_test_YWNlLWRpbm9zYXVyLTM5...` should correspond to same Clerk app as backend secret key

**Debug Steps**:
```bash
# 1. Check backend logs for authentication flow
tail -f apps/backend/logs/app.log | grep -E "token verification|JWT verified|user not authenticated"

# 2. Test token generation in browser console
await window.Clerk.session.getToken({ template: "api-test" })

# 3. Verify backend config
grep -E "ARK_AUTH.CLERK" apps/backend/.env

# 4. Check frontend config  
cat apps/frontend/.env.local
```

### Frontend won't start or shows blank page

**Check 1: Environment variables**
```bash
# Ensure .env.local exists and has required variables
cat apps/frontend/.env.local

# Should contain:
# VITE_CLERK_PUBLISHABLE_KEY=pk_test_...
# VITE_API_URL=http://localhost:8080
# VITE_ENV=local

# After changes, always restart:
bun dev
```

**Check 2: Browser console errors**
- Open DevTools (F12) → Console
- `ReferenceError: process is not defined` → Using `process.env` instead of `import.meta.env`
- `Invalid environment variables` → Missing or incorrect `.env.local`

**Check 3: OpenAPI contracts**
```bash
# Rebuild contracts if needed
cd packages/openapi && bun run build
cd packages/zod && bun run build
```

### Clerk authentication errors

**Invalid publishable key**:
```bash
# Verify you're using pk_test_... (publishable) not sk_test_... (secret)
echo $VITE_CLERK_PUBLISHABLE_KEY

# Get token in browser console to test:
await window.Clerk.session.getToken({ template: "custom" })
```

**Backend JWT verification fails**:
- Ensure backend has matching Clerk configuration
- Check `ARK_AUTH.CLERK.SECRET_KEY` (backend) matches your Clerk app
- Check `ARK_AUTH.CLERK.JWT_ISSUER` matches Clerk issuer
- JWT template name must match on both sides (currently "custom", migrating to "api-test")

### TypeScript errors after backend changes

```bash
# Rebuild shared packages
cd packages/zod && bun run build
cd ../openapi && bun run build

# Type check frontend
cd ../apps/frontend
bun typecheck
```

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


- `package.json` correctly uses `@ark/openapi` and `@ark/zod`
- **Action needed**: Update vite.config.ts aliases to match package.json dependencies

These inconsistencies don't currently break functionality but should be addressed for clarity.

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
- **Playwright** (for E2E tests): Installed via `bun install`, browsers via `bunx playwright install`

**Note**: Docker and Docker Compose are not yet configured. You must install PostgreSQL and Redis locally.
