# ARK Architecture Guide

## Project Overview

**Ark** is a self-hosted asset management system designed for homelab enthusiasts. It allows users to track their hardware infrastructure (servers, VMs, containers, network devices) and maintain a searchable history of configuration changes and maintenance logs.

### Project Goals

**Primary Purpose**: Create a production-ready tool for the homelab community while building resume-worthy achievements that demonstrate end-to-end development capabilities.

**Target Audience**: Homelab enthusiasts who value control over their infrastructure and want a self-hosted solution for asset tracking and log management.

**Key Differentiators**:
- **Self-hosted**: Full data ownership and control
- **AI-powered search**: Natural language querying of configuration logs (planned with RAG)
- **Cost-optimized**: Designed to run on minimal infrastructure ($15-25/month budget)
- **Production-ready**: Enterprise patterns (observability, CI/CD, clean architecture)

### Current Status

**Live Production**: https://arkcore.dev (deployed on AWS EC2)
- Backend API: https://api.arkcore.dev
- Frontend: React SPA with Clerk authentication
- Database: PostgreSQL 16 with pgvector extension
- Infrastructure: Docker Compose on t3.micro instance

**Development Phase**: MVP complete, V2 planning
- Asset and log CRUD operations: ✅ Complete
- Authentication and multi-tenancy: ✅ Complete
- Frontend asset views: ✅ Complete
- AI-powered semantic search: 🚧 Planned for V2
- Vector database integration: 🚧 Planned for V2

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
├── docs/
│   └── claude/           # Domain-specific documentation
├── package.json          # Workspace root
├── turbo.json            # Turborepo configuration
└── README.md
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
- **CRITICAL**: Clerk SDK must be initialized with `clerk.SetKey(cfg.Auth.Clerk.SecretKey)` in `internal/server/server.go` during server startup
- Phase 1: `ClerkAuthMiddleware` - Applied globally to all `/api/v1/*` routes
  - Validates JWT from `Authorization: Bearer <token>`
  - Verifies with Clerk SDK v2
  - Stores session claims in Echo context using key `clerk_session_claims`
- Phase 2: `RequireAuth` - Applied globally to all `/api/v1/*` routes (MUST be applied after ClerkAuthMiddleware)
  - Retrieves verified claims from context
  - Extracts user_id from claims.Subject
  - Sets `user_id` in context for handlers
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
