# Ark Project Context

> [!NOTE]
> This project was migrated from the garden_journal codebase. The migration is complete - all plant/observation code has been replaced with asset/log equivalents. Module name updated to `ark`, database name is `ark`.

# Ark Backend Guide

## Overview
The backend for **Ark** is built with Go (Golang) using the Echo framework. It provides a RESTful API for the frontend.
Ark is a homelab asset tracking and configuration log management application with AI-powered search. Built with Go (backend) and TypeScript/React (frontend) in a Turborepo monorepo.
Core Use Case: Track servers, VMs, containers, and network equipment while maintaining searchable logs of configuration changes. AI assistant helps query logs in natural language.

Tech Stack
Backend:

Go 1.24+, Echo framework
PostgreSQL 16+ (connection pooling, full-text search, trigram matching)
Redis 8+ for background jobs and caching
Clerk SDK v2 for authentication
OpenAI API (gpt-4o-mini) for AI queries
Zerolog for structured logging
New Relic for APM (optional but recommended)
Resend for transactional emails

Frontend:

React 19.1.0, TypeScript 5.8.2, Vite 7.0.4
TanStack Query for data fetching
Clerk for authentication
Tailwind CSS, React Router


Architecture
Backend Structure
bashapps/backend/
├── cmd/ark/                   # Main application entry point (was cmd/gardenjournal)
├── internal/
│   ├── config/                # Koanf-based configuration
│   ├── database/migrations/   # SQL migrations (tern)
│   ├── handler/               # HTTP handlers (Echo)
│   │   ├── asset.go              # Asset CRUD operations
│   │   ├── log.go                # Log CRUD operations
│   │   ├── ai.go                 # AI query endpoint (PLANNED - not yet implemented)
│   │   ├── health.go             # Health check endpoint
│   │   ├── openapi.go            # OpenAPI spec endpoint
│   │   ├── base.go               # Base handler utilities
│   │   └── handlers.go           # Handler aggregator
│   ├── service/               # Business logic
│   │   ├── asset_service.go      # Asset business logic
│   │   ├── log_service.go        # Log business logic
│   │   ├── ai_service.go         # RAG implementation (PLANNED - not yet implemented)
│   │   ├── auth.go               # Auth service interface
│   │   ├── services.go           # Service aggregator
│   │   └── test_helpers.go       # Testing utilities
│   ├── repository/            # Data access (PostgreSQL)
│   │   ├── asset_repository.go   # Asset data access
│   │   ├── log_repository.go     # Log data access with FTS methods
│   │   └── repositories.go       # Repository aggregator
│   ├── model/                 # Domain models & DTOs
│   │   ├── asset.go              # Asset domain model
│   │   ├── log.go                # Log domain model
│   │   ├── ai.go                 # AI query models (PLANNED - not yet implemented)
│   │   ├── pagination.go         # Pagination models
│   │   ├── response.go           # API response models
│   │   ├── base.go               # Base model utilities
│   │   └── weathersnapshot/      # Weather models (legacy, kept for reference)
│   ├── middleware/            # HTTP middleware (REUSE)
│   │   ├── auth.go               # Two-phase Clerk authentication
│   │   ├── global.go             # Error handling, logging, CORS
│   │   ├── context.go            # Request context management
│   │   ├── tracing.go            # New Relic tracing
│   │   └── middleware.go         # Middleware aggregator
│   ├── router/v1/             # Route registration
│   ├── lib/                   # Shared utilities
│   │   ├── jwt/                  # JWT verification
│   │   ├── llm/                  # LLM client (OpenAI) (PLANNED - not yet implemented)
│   │   ├── email/                # Email client (Resend)
│   │   ├── job/                  # Background job processing
│   │   ├── utils/                # General utilities
│   │   └── weather/              # Weather API integration (legacy, kept for reference)
│   ├── logger/                # Logging setup
│   ├── server/                # Server config
│   ├── validation/            # Request validation
│   ├── errs/                  # Error types
│   ├── sqlerr/                # SQL error handling
│   └── testing/               # Testing utilities
├── templates/              # Email templates
│   └── emails/
│       └── welcome.html        # Welcome email template
├── static/                 # Static files
│   └── openapi.json           # OpenAPI specification
└── tests/
    ├── integration/           # Integration tests
    └── manual/                # Manual test files
    └── manual/                # Manual test files
        ├── test_auth.http        # Authentication & Asset/Log CRUD testing
        ├── ai.http               # AI queries (TO BE CREATED)
        └── e2e_ai_flow.http      # E2E AI flow (TO BE CREATED)
Frontend Structure
bashapps/web/src/
├── components/
│   ├── assets/                # AssetList, AssetCard, AssetForm (replaces plants/)
│   ├── logs/                  # LogList, LogCard, LogForm (replaces observations/)
│   ├── ai/                    # AIQueryForm, AIResponse (NEW)
│   └── layout/                # Navbar, Layout (reuse, update branding)
├── hooks/                     # useAssets, useLogs, useAIQuery (replace usePlants, useObservations)
├── pages/                     # Dashboard, AssetDetailPage (replace PlantDashboard, etc.)
├── lib/                       # api.ts (Axios + auth), clerk.ts (reuse)
├── types/                     # TypeScript interfaces (update for Ark domain)
└── App.tsx                    # Routing (update routes)
```

---

## Domain Models

### Asset
Homelab asset (server, VM, container, NAS, network equipment, other)

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
Configuration change or troubleshooting log

| Field | Type | Notes |
|-------|------|-------|
| id | UUID | Primary key |
| asset_id | UUID | Foreign key (cascade delete) |
| user_id | string | Denormalized for performance |
| content | string | Required, 2-10,000 chars |
| tags | string[] | Optional, max 20 tags, 50 chars each |
| content_vector | tsvector | Generated column for FTS (hidden from API) |
| created_at | timestamp | Auto-set |
| updated_at | timestamp | Auto-updated by trigger |

**Example Log:**
```
Content: "Fixed nginx by updating /etc/nginx/nginx.conf and restarting with systemctl restart nginx"
Tags: ["nginx", "web-server", "fix"]
```

---

## API Endpoints

**Assets:**
```
GET    /api/v1/assets              # List (paginated)
POST   /api/v1/assets              # Create
GET    /api/v1/assets/:id          # Get single
PATCH  /api/v1/assets/:id          # Update
DELETE /api/v1/assets/:id          # Delete (cascades to logs)
```

**Logs:**
```
GET    /api/v1/assets/:id/logs     # List for asset (paginated)
POST   /api/v1/assets/:id/logs     # Create for asset
GET    /api/v1/logs/:id            # Get single
PATCH  /api/v1/logs/:id            # Update
DELETE /api/v1/logs/:id            # Delete
```

**AI:**
```
POST   /api/v1/ai/query            # Query logs with natural language
```

**System:**
```
GET    /health                     # Health check (database, redis)
GET    /openapi.json               # OpenAPI specification
```

Request/Response Examples:
json// AI Query Request
{
  "asset_id": "550e8400-e29b-41d4-a716-446655440000",
  "query": "How did I fix nginx?"
}

// AI Query Response
{
  "answer": "You fixed nginx on 2024-03-15 by updating /etc/nginx/nginx.conf and restarting the service.",
  "sources": ["660e8400-e29b-41d4-a716-446655440001"],
  "method": "recent"  // MVP: "recent", V1: "fts", V2: "vector"
}

// Health Check Response
{
  "status": "healthy",
  "checks": {
    "database": "ok",
    "redis": "ok"
  },
  "timestamp": "2024-03-15T14:30:00Z"
}
```

---

## Authentication Flow

**Two-Phase Pattern:**

1. **ClerkAuthMiddleware** (global on `/api/v1/*`)
   - Extracts JWT from `Authorization: Bearer <token>` header
   - Verifies with Clerk SDK v2
   - Stores `SessionClaims` in context
   - Returns 401 on failure

2. **RequireAuth** (per route group)
   - Extracts user_id from claims
   - Sets user context for handlers
   - All queries scoped to user_id (security-critical)

**Request Flow:**
```
Request → CORS/Logging → New Relic Tracing → ClerkAuth → RequireAuth → Validation → Handler → ErrorHandler → Response
Note: Authentication middleware from garden_journal can be reused as-is. Update route groups to apply to /assets and /logs instead of /plants and /observations.

Configuration
Required Environment Variables:
bash# Server
ARK_SERVER.PORT="8080"
ARK_SERVER.READ_TIMEOUT="30"
ARK_SERVER.WRITE_TIMEOUT="30"
ARK_SERVER.IDLE_TIMEOUT="60"
ARK_SERVER.CORS_ALLOWED_ORIGINS="http://localhost:3000"

# Database (PostgreSQL 16+)
ARK_DATABASE.HOST="localhost"
ARK_DATABASE.PORT="5432"
ARK_DATABASE.USER="ark_user"
ARK_DATABASE.PASSWORD="your_password"
ARK_DATABASE.NAME="ark"
ARK_DATABASE.SSL_MODE="disable"
ARK_DATABASE.MAX_OPEN_CONNS="25"
ARK_DATABASE.MAX_IDLE_CONNS="25"
ARK_DATABASE.CONN_MAX_LIFETIME="300"
ARK_DATABASE.CONN_MAX_IDLE_TIME="300"

# Clerk Authentication
ARK_AUTH.CLERK.SECRET_KEY="sk_test_..."
ARK_AUTH.CLERK.JWT_ISSUER="https://your-app.clerk.accounts.dev"
ARK_AUTH.CLERK.PEM_PUBLIC_KEY=""  # Optional for manual verification

# OpenAI (NEW)
ARK_OPENAI.API_KEY="sk-..."
ARK_OPENAI.MODEL="gpt-4o-mini"

# Redis (for background jobs)
ARK_REDIS.ADDRESS="localhost:6379"

# Resend (for emails)
ARK_INTEGRATION.RESEND_API_KEY="re_..."

# Observability
ARK_OBSERVABILITY.SERVICE_NAME="ark"
ARK_OBSERVABILITY.ENVIRONMENT="development"

# Logging
ARK_OBSERVABILITY.LOGGING.LEVEL="debug"
ARK_OBSERVABILITY.LOGGING.FORMAT="console"  # or "json" for production
ARK_OBSERVABILITY.LOGGING.SLOW_QUERY_THRESHOLD="100ms"

# New Relic (Optional but recommended for production)
ARK_OBSERVABILITY.NEW_RELIC.LICENSE_KEY="..."
ARK_OBSERVABILITY.NEW_RELIC.APP_LOG_FORWARDING_ENABLED="true"
ARK_OBSERVABILITY.NEW_RELIC.DISTRIBUTED_TRACING_ENABLED="true"

# Health Checks
ARK_OBSERVABILITY.HEALTH_CHECKS.ENABLED="true"
ARK_OBSERVABILITY.HEALTH_CHECKS.INTERVAL="30s"
ARK_OBSERVABILITY.HEALTH_CHECKS.TIMEOUT="5s"
ARK_OBSERVABILITY.HEALTH_CHECKS.CHECKS="database,redis"
Frontend .env:
bashVITE_CLERK_PUBLISHABLE_KEY=pk_test_...
VITE_API_URL=http://localhost:8080/api/v1
Migration Note: ARK-1 ticket renamed all GARDENJOURNAL_* prefixes to ARK_*.

Development Workflow
Backend
bashcd apps/backend
go mod download
cp .env.sample .env  # Edit with your credentials

task migrations:up   # Apply migrations
task run            # Start server (port 8080)
task test           # Run tests
task tidy           # Format code
Frontend
bashcd apps/web
bun install
bun dev             # Start dev server
Database Setup
bash# Create NEW database (not gardenjournal)
createdb -U postgres ark

# Grant permissions
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ark_user;"
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ark_user;"

# Apply NEW migrations (assets & asset_logs tables)
task migrations:up
Migration Strategy:

Starting fresh: Create new ark database, apply Ark migrations
Have garden_journal data: Optionally export plants/observations, transform to assets/logs, import (not required for MVP)

Getting JWT Tokens
Browser Console:
javascriptawait window.Clerk.session.getToken({ template: "api-test" })
```

---

## Production Features

### Observability

**Structured Logging (Zerolog):**
- Request-scoped context (request_id, trace_id, span_id)
- User-scoped logging (user_id in all logs)
- Different log levels: DEBUG, INFO, WARN, ERR, FTL
- Configurable formats: console (dev) or JSON (production)
- Slow query logging (configurable threshold)

**New Relic APM Integration:**
- Distributed tracing across services
- Database query monitoring
- Custom transaction naming
- Error tracking with stack traces
- Application log forwarding
- Performance metrics and dashboards

**Health Checks:**
- Endpoint: `GET /health`
- Checks database connectivity
- Checks Redis connectivity
- Configurable interval and timeout
- Used by load balancers and monitoring

### Background Job Processing

**Redis-based Job Queue (Asynq):**
- Async task processing
- Job retries with exponential backoff
- Scheduled/delayed jobs
- Job prioritization
- Worker pools

**Current Job Types:**
- Email sending (welcome emails, notifications)
- Future: Report generation, data exports, cleanup tasks

**Job Service** (`internal/lib/job/`):
- Job enqueueing
- Worker management
- Error handling
- Monitoring

### Email Integration

**Resend API:**
- Transactional email sending
- HTML email templates
- Email tracking
- High deliverability

**Email Client** (`internal/lib/email/`):
- Template rendering
- Email sending
- Error handling
- Logging

**Email Templates** (`templates/emails/`):
- Welcome email
- Future: Password reset, notifications, digests

**Template Updates Needed:**
- Change branding from "garden_journal" to "Ark"
- Update sender name/email
- Update template content for homelab context

---

## AI Implementation (RAG)

**MVP Approach (Current):**
1. User asks question about asset
2. Verify asset ownership (security)
3. Retrieve **recent 10 logs** from asset
4. Build prompt: system instructions + asset context + logs + question
5. Call OpenAI API (30s timeout)
6. Return answer with source log IDs

**Prompt Structure:**
```
You are a homelab assistant. Answer based ONLY on the following logs.

Asset: Homelab Server

Recent Logs:
[1] [2024-03-15] Fixed nginx by updating /etc/nginx/nginx.conf...
    Tags: nginx, fix

[2] [2024-03-10] Nginx throwing 502 errors...
    Tags: nginx, error

Question: How did I fix nginx?

Provide a concise answer with specific dates when relevant.
Cost: ~$0.0003 per query (1000 queries ≈ $0.30)
V1 Upgrade: Replace recent logs with FTS keyword search
V2 Upgrade: Add vector embeddings for semantic search

Database Features
Tables (NEW for Ark):

assets - Replaces plants table
asset_logs - Replaces observations table

Full-Text Search:

content_vector tsvector generated automatically on INSERT/UPDATE
GIN index for fast FTS queries
English language stemming

Indexes:

user_id - Security-critical (all queries scoped)
asset_id - Foreign key joins
created_at - Chronological ordering
content_vector - FTS performance
tags - GIN index for array operations
name (trigram) - Fuzzy asset name search

Triggers:

Auto-update updated_at timestamp on changes


## Testing

**Manual Tests** (`tests/manual/`):
- `test_auth.http` - Authentication testing (EXISTS)
- `migration_test_data.sql` - Test data for migration validation (EXISTS)
- `migration_validation.sql` - Migration validation queries (EXISTS)
- `asset.http` - Asset CRUD with error cases (TO BE CREATED)
- `log.http` - Log CRUD with tags validation (TO BE CREATED)
- `ai.http` - AI queries with various questions (TO BE CREATED)
- `e2e_ai_flow.http` - Complete flow: create asset → add logs → query AI (TO BE CREATED)

**Integration Tests:**
- Full middleware chain
- Auth verification
- Error handling (404, 400, 401, 504)
- Uses httptest for request/response
- New Relic transaction tracking

**Unit Tests:**
- Config loading
- JWT verification
- Middleware behavior
- Service layer logic
- Repository data access
- Model validation
- LLM client (with mocks, when implemented)

## Common Issues

**Port in use:**
```bash
lsof -ti:8080 | xargs kill -9
```

**Database permissions:**
```bash
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ark_user;"
```

**JWT invalid/expired:**
- Verify at https://jwt.io
- Check `iss` matches `CLERK.JWT_ISSUER`
- Generate fresh token from browser console

**OpenAI errors:**
- 401: Check `OPENAI.API_KEY`
- 429: Rate limit, wait or upgrade plan
- 504: Timeout, try simpler question

**AI returns "no logs":**
- Verify logs exist for asset
- Check `asset_id` matches in request

**Redis connection errors:**
- Verify Redis is running: `redis-cli ping`
- Check `REDIS.ADDRESS` in .env

**Email sending failures:**
- Verify Resend API key is valid
- Check Resend dashboard for errors
- Review email service logs


Error Handling
Custom Types (internal/lib/errs/):

NotFoundError → 404
ValidationError → 400
AuthError → 401
Generic errors → 500

Global Handler (internal/middleware/global.go):

Logs errors with context (request_id, user_id, trace_id, span_id)
Returns JSON errors
Hides stack traces in production
Integrates with New Relic error tracking

Note: Error handling from garden_journal can be reused as-is. No changes needed to middleware/global.go.

What to Reuse from Garden Journal
✅ KEEP (production-ready infrastructure):

Observability:

New Relic integration (internal/logger/logger.go)
Structured logging with Zerolog
Health check system
Request tracing


Background Jobs:

Redis integration (internal/lib/job/)
Asynq job queue
Worker management
Job handlers


Email System:

Resend client (internal/lib/email/)
Email templates (templates/emails/)
Email sending infrastructure


Authentication & Security:

Clerk middleware (internal/middleware/auth.go)
JWT verification (internal/lib/jwt/)
Two-phase auth pattern


Infrastructure:

Error types (internal/lib/errs/)
Global middleware (CORS, logging, recovery)
Server configuration
Database connection pooling
Migration system
Validation package
Config management (Koanf)



🔄 ADAPT (update for Ark domain):

Handlers (plant → asset, observation → log)
Services (business logic for assets/logs)
Repositories (data access for assets/logs)
Models (Asset, AssetLog instead of Plant, Observation)
Routes (update endpoints)
Frontend components (plants → assets, observations → logs)
API client (update endpoints)
TanStack Query hooks (update query keys)
Email templates (update branding from garden_journal to Ark)
OpenAPI spec (update for Ark endpoints)

➕ ADD (new for Ark):

LLM client (internal/lib/llm/)
AI service (internal/service/ai_service.go)
AI handler (internal/handler/ai.go)
AI DTOs (internal/model/ai.go)
AI frontend components (components/ai/)
Full-text search in log repository
OpenAI configuration
AI query endpoint and routes

❌ DELETE (domain-specific):

Plant-related files (handlers, services, repos, models)
Observation-related files (handlers, services, repos, models)
Plant/observation routes
Plant/observation frontend components
Old test files for plants/observations
Weather API integration (not needed for Ark, but keep code for reference)


Future Roadmap
V1 (Planned - 4-6 weeks)

FTS Search: Replace "recent logs" with keyword-based full-text search
AI Enhancements: Query history, copy answer, regenerate, suggested questions
UI Polish: Dark mode, markdown rendering, tag autocomplete, export
Performance: Redis caching for AI responses, rate limiting, pagination controls
Background Jobs: Scheduled log backups, data exports

V2 (Future - 8-10 weeks)

Vector Search: Semantic search with OpenAI embeddings + pgvector
Multi-Asset Queries: Search across all assets ("show me all nginx fixes")
Collaboration: Team workspaces, shared assets, activity feeds
Integrations: Webhooks, Slack/Discord, API automation
Analytics: Usage dashboards, common issues, resolution tracking
Mobile: PWA, offline mode, quick log entry
Email Notifications: Asset alerts, weekly digests

V3+ (Future)

Advanced AI: Streaming responses, custom instructions, AI-suggested tags
Automation: Auto-tagging, pattern recognition, anomaly detection
Extended Integrations: Monitoring tools (Prometheus, Grafana), ticketing systems
Advanced Email: Rich templates, personalization, A/B testing


Security Notes

All queries scoped by user_id (multi-tenancy)
Asset ownership verified before AI queries
JWT validation on every request
CORS configured for allowed origins
Input validation (max lengths, required fields)
API keys never committed (use .env, add to .gitignore)
Rate limiting planned for V1
SQL injection protection (parameterized queries)
XSS protection (Echo's built-in)


Performance Considerations
Current (MVP):

Recent 10 logs: <100ms query time
AI response: 5-15s typical (LLM latency)
No caching (comes in V1)

Optimizations (V1):

Redis for AI response caching
FTS with ranked results
Connection pooling (pgxpool)
Query result pagination
Background job processing for expensive operations

Scaling (V2):

Vector search with pgvector
Hybrid search (FTS + vector)
Horizontal scaling with load balancer
Database read replicas
CDN for static assets


Monitoring & Observability
New Relic Dashboards:

Transaction throughput and response times
Database query performance
Error rates and types
AI query latency and cost tracking
Background job success/failure rates

Key Metrics to Track:

API endpoint latency (p50, p95, p99)
Database connection pool utilization
Redis job queue depth
OpenAI API latency and token usage
Email delivery rates
User authentication success/failure
Health check status

Alerts (Recommended):

API error rate > 5%
Database connection pool exhausted
Redis job queue backlog > 1000
Health check failures
OpenAI API errors
Email delivery failures