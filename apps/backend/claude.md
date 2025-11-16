# Ark Project Context

## ⚠️ Project Transition Notice

**This repository was previously the `garden_journal` project and is being repurposed for Ark.**

**What this means:**
- Any references to "garden_journal", "plants", "observations" should be **replaced** with Ark equivalents (assets, logs)
- Old code/files are being **migrated or removed** - see MVP tickets (ARK-1 through ARK-22)
- If you encounter garden_journal code, it should be **updated to Ark** or **deleted**
- Configuration may still reference old names - these need updating

**Key Changes:**
| Garden Journal | Ark Equivalent |
|----------------|----------------|
| Plants | Assets |
| Observations | Logs (AssetLogs) |
| Plant species/variety | Asset type/hostname |
| Plant notes | Log content |
| - | AI query (new feature) |

**During Development:**
- Prioritize Ark implementation over garden_journal compatibility
- Remove old plant/observation routes, handlers, models as you build Ark equivalents
- Update import paths from `garden_journal` to `ark`
- Rename module in `go.mod` from `garden_journal` to `ark`
- Update database name from `gardenjournal` to `ark`

---

## Overview
**Ark** is a homelab asset tracking and configuration log management application with AI-powered search. Built with **Go (backend)** and **TypeScript/React (frontend)** in a **Turborepo** monorepo.

**Core Use Case:** Track servers, VMs, containers, and network equipment while maintaining searchable logs of configuration changes. AI assistant helps query logs in natural language.

---

## Tech Stack

**Backend:**
- Go 1.24+, Echo framework
- PostgreSQL 16+ (connection pooling, full-text search, trigram matching)
- Clerk SDK v2 for authentication
- OpenAI API (gpt-4o-mini) for AI queries
- Zerolog for structured logging

**Frontend:**
- React 19.1.0, TypeScript 5.8.2, Vite 7.0.4
- TanStack Query for data fetching
- Clerk for authentication
- Tailwind CSS, React Router

---

## Architecture

### Backend Structure
```bash
apps/backend/
├── cmd/ark/                   # Main application entry point (was cmd/gardenjournal)
├── internal/
│   ├── config/                # Koanf-based configuration
│   ├── database/migrations/   # SQL migrations (tern)
│   ├── handler/               # HTTP handlers (Echo)
│   │   ├── asset.go              # Asset CRUD (replaces plant.go)
│   │   ├── log.go                # Log CRUD (replaces observation.go)
│   │   └── ai.go                 # AI query endpoint (NEW)
│   ├── service/               # Business logic
│   │   ├── asset_service.go      # (replaces plant_service.go)
│   │   ├── log_service.go        # (replaces observation_service.go)
│   │   └── ai_service.go         # RAG implementation (NEW)
│   ├── repository/            # Data access (PostgreSQL)
│   │   ├── asset_repository.go   # (replaces plant_repository.go)
│   │   └── log_repository.go     # Includes FTS methods (replaces observation_repository.go)
│   ├── model/                 # Domain models & DTOs
│   │   ├── asset.go              # (replaces plant.go)
│   │   ├── log.go                # (replaces observation.go)
│   │   └── ai.go                 # (NEW)
│   ├── middleware/            # Auth, CORS, logging, errors (reuse from garden_journal)
│   ├── router/v1/             # Route registration
│   ├── lib/
│   │   ├── jwt/                  # JWT verification (reuse)
│   │   ├── errs/                 # Custom error types (reuse)
│   │   └── llm/                  # LLM client (OpenAI) (NEW)
│   └── server/                # Server config (reuse)
└── tests/
    ├── integration/           # Integration tests
    └── manual/                # .http files (asset, log, ai, e2e)
```

### Frontend Structure
```bash
apps/web/src/
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

**OLD Endpoints (TO BE REMOVED):**
```
❌ /api/v1/plants                  # DELETE - replaced by /assets
❌ /api/v1/observations            # DELETE - replaced by /logs
```

**Request/Response Examples:**
```json
// AI Query Request
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
Request → CORS/Logging → ClerkAuth → RequireAuth → Validation → Handler → ErrorHandler → Response
```

**Note:** Authentication middleware from garden_journal can be **reused as-is**. Update route groups to apply to `/assets` and `/logs` instead of `/plants` and `/observations`.

---

## Configuration

**Required Environment Variables:**
```bash
# Server
ARK_SERVER.PORT="8080"                    # Was GARDENJOURNAL_SERVER.PORT

# Database (PostgreSQL 16+)
ARK_DATABASE.HOST="localhost"             # Was GARDENJOURNAL_DATABASE.HOST
ARK_DATABASE.USER="ark_user"              # Was GARDENJOURNAL_DATABASE.USER
ARK_DATABASE.PASSWORD="your_password"     # Was GARDENJOURNAL_DATABASE.PASSWORD
ARK_DATABASE.NAME="ark"                   # Was GARDENJOURNAL_DATABASE.NAME="gardenjournal"

# Clerk Authentication (REUSE - no change)
ARK_AUTH.CLERK.SECRET_KEY="sk_test_..."
ARK_AUTH.CLERK.JWT_ISSUER="https://your-app.clerk.accounts.dev"

# OpenAI (NEW)
ARK_OPENAI.API_KEY="sk-..."
ARK_OPENAI.MODEL="gpt-4o-mini"
```

**Frontend `.env`:**
```bash
VITE_CLERK_PUBLISHABLE_KEY=pk_test_...
VITE_API_URL=http://localhost:8080/api/v1
```

**Migration Note:** Update all `GARDENJOURNAL_*` prefixes to `ARK_*` in your `.env` file.

---

## Development Workflow

### Backend
```bash
cd apps/backend
go mod download
cp .env.sample .env  # Edit with your credentials

task migrations:up   # Apply migrations
task run            # Start server (port 8080)
task test           # Run tests
task tidy           # Format code
```

### Frontend
```bash
cd apps/web
bun install
bun dev             # Start dev server
```

### Database Setup
```bash
# Create NEW database (not gardenjournal)
createdb -U postgres ark

# Grant permissions
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ark_user;"
psql -U postgres -d ark -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ark_user;"

# Apply NEW migrations (assets & asset_logs tables)
task migrations:up
```

**Migration Strategy:**
- **Starting fresh:** Create new `ark` database, apply Ark migrations
- **Have garden_journal data:** Optionally export plants/observations, transform to assets/logs, import (not required for MVP)

### Getting JWT Tokens
**Browser Console:**
```javascript
await window.Clerk.session.getToken({ template: "api-test" })
```

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
```

**Cost:** ~$0.0003 per query (1000 queries ≈ $0.30)

**V1 Upgrade:** Replace recent logs with FTS keyword search  
**V2 Upgrade:** Add vector embeddings for semantic search

---

## Database Features

**Tables (NEW for Ark):**
- `assets` - Replaces `plants` table
- `asset_logs` - Replaces `observations` table

**Full-Text Search:**
- `content_vector` tsvector generated automatically on INSERT/UPDATE
- GIN index for fast FTS queries
- English language stemming

**Indexes:**
- `user_id` - Security-critical (all queries scoped)
- `asset_id` - Foreign key joins
- `created_at` - Chronological ordering
- `content_vector` - FTS performance
- `tags` - GIN index for array operations
- `name` (trigram) - Fuzzy asset name search

**Triggers:**
- Auto-update `updated_at` timestamp on changes

---

## Testing

**Manual Tests** (`tests/manual/*.http`):
- `asset.http` - Asset CRUD with error cases (replaces plant.http)
- `log.http` - Log CRUD with tags validation (replaces observation.http)
- `ai.http` - AI queries with various questions (NEW)
- `e2e_ai_flow.http` - Complete flow: create asset → add logs → query AI (NEW)

**Integration Tests:**
- Full middleware chain
- Auth verification
- Error handling (404, 400, 401, 504)
- Uses `httptest` for request/response

**Unit Tests:**
- Config loading
- JWT verification (reuse from garden_journal)
- Middleware behavior (reuse from garden_journal)
- LLM client (NEW, with mocks)

**Note:** Garden_journal test patterns can be reused. Update test data from plants/observations to assets/logs.

---

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

**Module name errors:**
- If you see import errors like `cannot find package "garden_journal"`, run ARK-1 ticket to rename module in `go.mod`

**Old routes still registered:**
- Remove plant/observation route registration from `router/v1/routes.go`
- Delete old handler files after replacing with Ark equivalents

---

## Error Handling

**Custom Types** (`internal/lib/errs/`):
- `NotFoundError` → 404
- `ValidationError` → 400
- `AuthError` → 401
- Generic errors → 500

**Global Handler:**
- Logs with context (request_id, user_id)
- Returns JSON errors
- Hides stack traces in production

**Note:** Error handling from garden_journal can be **reused as-is**. No changes needed to middleware/global.go.

---

## What to Reuse from Garden Journal

**✅ REUSE (no changes needed):**
- Authentication middleware (`internal/middleware/auth.go`)
- JWT verification (`internal/lib/jwt/`)
- Error types (`internal/lib/errs/`)
- Global middleware (CORS, logging, recovery)
- Server configuration
- Clerk integration
- Database connection pooling
- Migration system structure (just new migration files)
- Validation package
- Config management pattern (Koanf)

**🔄 ADAPT (update for Ark domain):**
- Handlers (plant → asset, observation → log)
- Services (business logic for assets/logs)
- Repositories (data access for assets/logs)
- Models (Asset, AssetLog instead of Plant, Observation)
- Routes (update endpoints)
- Frontend components (plants → assets, observations → logs)
- API client (update endpoints)
- TanStack Query hooks (update query keys)

**➕ ADD (new for Ark):**
- LLM client (`internal/lib/llm/`)
- AI service (`internal/service/ai_service.go`)
- AI handler (`internal/handler/ai.go`)
- AI DTOs (`internal/model/ai.go`)
- AI frontend components (`components/ai/`)
- Full-text search in log repository
- OpenAI configuration
- AI query endpoint and routes

**❌ DELETE (no longer needed):**
- Plant-related files (handlers, services, repos, models)
- Observation-related files (handlers, services, repos, models)
- Plant/observation routes
- Plant/observation frontend components
- Old test files for plants/observations
---

## Future Roadmap

### V1 (Planned - 4-6 weeks)
- **FTS Search:** Replace "recent logs" with keyword-based full-text search
- **AI Enhancements:** Query history, copy answer, regenerate, suggested questions
- **UI Polish:** Dark mode, markdown rendering, tag autocomplete, export
- **Performance:** Redis caching, rate limiting, pagination controls

### V2 (Future - 8-10 weeks)
- **Vector Search:** Semantic search with OpenAI embeddings + pgvector
- **Multi-Asset Queries:** Search across all assets ("show me all nginx fixes")
- **Collaboration:** Team workspaces, shared assets, activity feeds
- **Integrations:** Webhooks, Slack/Discord, API automation
- **Analytics:** Usage dashboards, common issues, resolution tracking
- **Mobile:** PWA, offline mode, quick log entry

### V3+ (Future)
- **Advanced AI:** Streaming responses, custom instructions, AI-suggested tags
- **Automation:** Auto-tagging, pattern recognition, anomaly detection
- **Extended Integrations:** Monitoring tools (Prometheus, Grafana), ticketing systems

---

## Security Notes

- All queries scoped by `user_id` (multi-tenancy)
- Asset ownership verified before AI queries
- JWT validation on every request
- CORS configured for allowed origins
- Input validation (max lengths, required fields)
- API keys never committed (use .env, add to .gitignore)
- Rate limiting planned for V1

---

## Performance Considerations

**Current (MVP):**
- Recent 10 logs: <100ms query time
- AI response: 5-15s typical (LLM latency)
- No caching (comes in V1)

**Optimizations (V1):**
- Redis for response caching
- FTS with ranked results
- Connection pooling (pgxpool)
- Query result pagination

**Scaling (V2):**
- Vector search with pgvector
- Hybrid search (FTS + vector)
- Background job processing for expensive queries