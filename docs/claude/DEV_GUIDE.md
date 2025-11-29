# ARK Development Guide

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

## Prerequisites

### Local Development
- **Go**: 1.24 or higher
- **Node.js**: 22+ (for frontend and package builds)
- **Bun**: Latest version (package manager)
- **PostgreSQL**: 16+ (local installation required for development)
- **Redis**: 8+ (local installation required for development)
- **Task**: Task runner for Go backend (install via `brew install go-task`)
- **Tern**: Database migration tool (installed via Go modules)
- **Playwright** (for E2E tests): Installed via `bun install`, browsers via `bunx playwright install`
