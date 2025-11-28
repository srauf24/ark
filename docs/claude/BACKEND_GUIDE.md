# ARK Backend Guide

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
