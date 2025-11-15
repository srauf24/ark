#  Garden Journal Project Context

##  Project Overview
**Garden Journal** is a modern web application for plant care and garden management built with **Go (backend)** and **TypeScript/React (frontend)**.  
It follows a **monorepo architecture** using **Turborepo** for efficient builds and development workflow.

---

## ⚙ Technical Stack

### 🐹 Backend (Go)
- **Go 1.24+**
- **Echo** framework for REST API
- **PostgreSQL 16+** with connection pooling
- **Redis 8+** for background jobs
- **Clerk SDK** for authentication
- **New Relic** for APM
- **Resend** for email services

### ⚛ Frontend (TypeScript / React)
- **React 19.1.0**
- **TypeScript 5.8.2**
- **Vite 7.0.4**
- **TanStack Query** for data fetching
- **Clerk** for authentication
- **Tailwind CSS** for styling
- **React Router** for navigation

---

## 🏗 Architecture

###  Backend Structure
```bash
apps/backend/
├── cmd/                    # Application entry points
│   └── gardenjournal/         # Main application
├── internal/               # Private application code
│   ├── config/                # Configuration management (Koanf-based)
│   ├── database/              # Database connections and migrations
│   ├── handler/               # HTTP request handlers (Echo handlers)
│   │   ├── plant.go              # Plant CRUD operations
│   │   └── observation.go        # Observation CRUD operations
│   ├── service/               # Business logic layer
│   ├── repository/            # Data access layer (PostgreSQL)
│   ├── model/                 # Domain models and DTOs
│   ├── middleware/            # HTTP middleware
│   │   ├── auth.go               # Two-phase authentication
│   │   ├── global.go             # Error handling, logging, CORS
│   │   └── middleware.go         # Middleware aggregator
│   ├── router/                # Route registration
│   │   └── v1/                   # API v1 routes
│   ├── validation/            # Request validation
│   ├── lib/                   # Shared utilities
│   │   ├── jwt/                  # JWT verification helpers
│   │   └── errs/                 # Custom error types
│   └── server/                # Server configuration
├── templates/              # Email templates
├── static/                 # Static files
└── tests/                  # Test suites
    ├── integration/           # Integration tests
    └── manual/                # Manual test scripts (.http files)
````

###  Key Backend Features

1. **Configuration Management**

   * Environment-based configuration using Koanf
   * Structured validation
   * Support for multiple environments

2. **Database Layer**

   * PostgreSQL with connection pooling
   * Migration system using `tern`
   * Configurable connection settings

3. **Authentication & Security**

   **Two-Phase Authentication Pattern:**

   * **Phase 1 - ClerkAuthMiddleware** (`internal/middleware/auth.go`)
     - Applied globally to all `/api/v1/*` routes
     - Extracts JWT from `Authorization: Bearer <token>` header
     - Verifies token using Clerk SDK v2 (`clerk/clerk-sdk-go/v2`)
     - Stores validated `SessionClaims` in Echo context
     - Returns 401 for missing, invalid, or expired tokens

   * **Phase 2 - RequireAuth** (`internal/middleware/auth.go`)
     - Applied to individual route groups (plants, observations)
     - Retrieves `SessionClaims` from context
     - Extracts user metadata (user_id, role, permissions)
     - Sets user data in context for downstream handlers

   **JWT Verification** (`internal/lib/jwt/clerk.go`):
   - RS256 signature validation
   - Issuer verification against configured Clerk domain
   - Expiration checking with detailed error messages
   - Bearer token extraction with format validation

   **Configuration** (`internal/config/config.go`):
   - Clerk Secret Key (from Clerk Dashboard)
   - JWT Issuer URL (e.g., `https://your-app.clerk.accounts.dev`)
   - Optional PEM public key for manual verification

   **Security Features:**
   - CORS with configurable allowed origins
   - Custom error handling (no stack traces in production)
   - Structured logging with request IDs
   - User-scoped data access (all queries filtered by user_id)

4. **Background Processing**

   * Redis-based job queue
   * Async task processing
   * Email notifications

5. **Observability**

   * New Relic APM integration
   * Structured logging (`zerolog`)
   * Health checks
   * Performance monitoring

---

###  Frontend Structure

```bash
apps/frontend/
├── src/
│   ├── components/   # Reusable UI components
│   ├── features/     # Feature-specific code
│   ├── hooks/        # Custom React hooks
│   ├── pages/        # Route pages
│   ├── api/          # API integration
│   ├── utils/        # Utility functions
│   └── styles/       # Global styles
└── tests/            # Frontend tests
```

---

##  Development Workflow

###  Backend Development

1. **Environment Setup**

   ```bash
   cd apps/backend
   go mod download
   cp .env.sample .env
   # Edit .env with your Clerk credentials and database settings
   ```

2. **Database Management**

   ```bash
   task migrations:new name=<migration_name>  # Create migration
   task migrations:up                         # Apply migrations

   # Grant database permissions (if needed)
   psql -U postgres -d gardenjournal -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO <db_user>;"
   ```

3. **Running the Server**

   ```bash
   task run    # Start server on port 8080
   task test   # Run all tests
   task tidy   # Format and tidy code
   ```

4. **Development Best Practices**

   **Incremental Development:**
   - Break large features into small, testable steps
   - Write implementation + tests for each step
   - Verify tests pass before proceeding
   - Commit after each completed step with clear messages

   **Testing Workflow:**
   - Write unit tests for new functions (`*_test.go` files)
   - Add integration tests for route changes
   - Use manual `.http` files to verify with real data
   - Run `task test` to ensure all tests pass

   **Git Workflow:**
   - Create feature branches for new work
   - Make atomic commits with descriptive messages
   - Include test files in the same commit as implementation
   - Push commits regularly to remote branch
   - Create a pull requests once all features in ticket are complete

   **Debugging:**
   - Check server logs for detailed error messages
   - Use structured logging to trace request flow
   - Verify JWT tokens at https://jwt.io
   - Test database permissions with `psql` commands
   - Kill orphaned processes: `lsof -ti:8080 | xargs kill -9`

5. **Configuration Management**

   **Required Environment Variables:**
   ```bash
   # Server
   GARDENJOURNAL_SERVER.PORT="8080"

   # Database
   GARDENJOURNAL_DATABASE.HOST="localhost"
   GARDENJOURNAL_DATABASE.USER="your_db_user"
   GARDENJOURNAL_DATABASE.PASSWORD="your_db_password"
   GARDENJOURNAL_DATABASE.NAME="gardenjournal"

   # Clerk Authentication
   GARDENJOURNAL_AUTH.CLERK.SECRET_KEY="sk_test_..."
   GARDENJOURNAL_AUTH.CLERK.JWT_ISSUER="https://your-app.clerk.accounts.dev"
   ```

6. **Getting Fresh JWT Tokens**

   **Method 1 - Browser Console:**
   ```javascript
   await window.Clerk.session.getToken({ template: "api-test" })
   ```

   **Method 2 - Network Tab:**
   - Open DevTools Network tab
   - Make authenticated request in frontend
   - Copy Bearer token from Authorization header

### ⚛ Frontend Development

1. **Setup**

   ```bash
   bun install
   ```

2. **Development**

   ```bash
   bun dev     # Start dev server
   bun build   # Production build
   bun lint    # Run linter
   ```

---

##  API Structure

**Endpoint Pattern:**
```
/api/v1/<resource>
```

**Available Resources:**
- `/api/v1/plants` - Plant management (CRUD)
- `/api/v1/observations` - Observation tracking (CRUD)

**Middleware Chain:**
```
Request
  ↓
Global Middleware (CORS, Logging, Recovery)
  ↓
ClerkAuthMiddleware (JWT Verification) [/api/v1/*]
  ↓
RequireAuth (Claims Extraction) [Resource-specific]
  ↓
Validation Middleware
  ↓
Handler (Business Logic)
  ↓
Global Error Handler
  ↓
Response
```

**Request Flow Example:**
1. Client sends `GET /api/v1/plants` with `Authorization: Bearer <jwt>`
2. Global middleware logs request and adds request_id
3. ClerkAuthMiddleware verifies JWT and stores claims in context
4. RequireAuth extracts user_id from claims
5. Handler validates request and queries database (filtered by user_id)
6. Response sent with structured JSON
7. Request logged with duration, status, and user_id

**Response Format:**
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

**Error Response Format:**
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

---

##  API Integration (Frontend)

* REST API with **OpenAPI/Swagger specification**
* Type-safe API client using **ts-rest**
* Automatic type generation from OpenAPI specs
* Request/response validation
* Error handling with retries

---

##  Error Handling

**Custom Error Types** (`internal/lib/errs/`):
* `HTTPError` - Structured error responses with status codes
* Error codes: `UNAUTHORIZED`, `INTERNAL_SERVER_ERROR`, `VALIDATION_ERROR`, etc.
* Consistent JSON error format across all endpoints

**Global Error Handler** (`internal/middleware/global.go`):
* Centralized error handling for all routes
* Logs errors with request context (request_id, user_id, method, path)
* Returns appropriate HTTP status codes
* Hides internal error details in production

**Error Flow:**
```
Handler Error → Custom HTTPError → Global Error Handler → JSON Response
```

**Logging:**
* Structured logging with Zerolog
* Request-scoped context (request_id, trace_id, span_id)
* Different log levels: DEBUG, INFO, WARN, ERR, FTL
* Integration with New Relic for distributed tracing

---

##  Testing Strategy

### Backend

**Unit Tests:**
* Configuration loading (`internal/config/config_test.go`)
* JWT verification helpers (`internal/lib/jwt/clerk_test.go`)
* Middleware behavior (`internal/middleware/auth_test.go`)
* Mock-based testing with `testify/assert` and `testify/require`

**Integration Tests** (`tests/integration/`):
* Full middleware chain testing
* Route-level authentication verification
* Error handling and HTTP status codes
* Uses `httptest` for request/response simulation
* Tests all CRUD operations (GET, POST, PUT, DELETE)

**Manual Testing** (`tests/manual/`):
* `.http` files for REST client testing (Bruno, HTTPie, Postman, VS Code REST Client)
* Comprehensive test cases for all auth scenarios:
  - Missing Authorization header
  - Invalid token formats
  - Malformed JWTs
  - Valid JWT authentication
* Includes verification checklists and troubleshooting guides

**Testing Best Practices:**
* Write tests alongside implementation (TDD approach)
* Test both positive and negative cases
* Use table-driven tests for multiple scenarios
* Verify error messages and status codes
* Test middleware chain ordering
* Run `task test` before committing

### Frontend

* Component tests with **React Testing Library**
* Integration tests
* E2E tests with **Cypress**

---

##  Deployment & Operations

1. **Environment Configuration**

   * Environment-specific settings
   * Secret management
   * Feature flags

2. **Monitoring**

   * APM with New Relic
   * Error tracking
   * Performance monitoring
   * Log aggregation

3. **Security**

   * Authentication with Clerk
   * Authorization middleware
   * Input validation
   * Rate limiting
   * CORS policies

---

##  Common Issues & Troubleshooting

### Port Already in Use
```bash
# Error: listen tcp :8080: bind: address already in use
lsof -ti:8080 | xargs kill -9
task run
```

### Database Permission Denied
```bash
# Error: permission denied for table plants
psql -U postgres -d gardenjournal -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO <db_user>;"
psql -U postgres -d gardenjournal -c "GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO <db_user>;"
```

### JWT Token Invalid/Expired
- Decode token at https://jwt.io to check expiration
- Verify `iss` claim matches `CLERK.JWT_ISSUER` in .env
- Generate fresh token from browser console (see "Getting Fresh JWT Tokens" above)
- Ensure Clerk Secret Key is correct

### 401 Unauthorized Errors
1. Check Authorization header is present: `Authorization: Bearer <token>`
2. Verify token format (should have 3 parts separated by dots)
3. Check server logs for specific error message
4. Confirm Clerk credentials in .env are correct
5. Test with a fresh token

### Database Connection Issues
```bash
# Test connection
psql -U <db_user> -d gardenjournal -c "SELECT 1"

# Check if database exists
psql -U postgres -l | grep gardenjournal
```

### Server Won't Start
1. Check .env file exists and has required variables
2. Verify database is running
3. Verify Redis is running (if using background jobs)
4. Check for port conflicts
5. Review server logs for configuration errors

---

##  Future Considerations

1. **Scalability**

   * Horizontal scaling of API
   * Caching strategies (Redis for session data)
   * Database optimization (indexes, query optimization)
   * Connection pooling tuning

2. **Feature Enhancements**

   * Real-time updates (WebSockets)
   * Mobile responsiveness
   * Offline support (PWA)
   * Data export/import (CSV, JSON)
   * Image upload for plant photos
   * Reminder notifications for watering

3. **Integrations**

   * Weather API integration
   * Plant database (species information)
   * Image recognition (plant identification)
   * Social sharing
   * Calendar integration

4. **Security Enhancements**

   * Rate limiting per user
   * API key management
   * Role-based access control (RBAC)
   * Audit logging
   * HTTPS enforcement in production

