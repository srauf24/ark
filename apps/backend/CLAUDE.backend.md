# Backend Quick Reference

## File Structure

```
apps/backend/
├── cmd/ark/              # Application entry point
├── internal/
│   ├── handler/          # HTTP handlers
│   ├── service/          # Business logic
│   ├── repository/       # Data access
│   ├── model/            # Domain models and DTOs
│   ├── middleware/       # Request pipeline
│   ├── router/           # Route registration
│   ├── server/           # Server setup
│   ├── database/         # DB connection and migrations
│   ├── lib/              # Shared libraries (email, jobs)
│   └── errs/             # Custom error types
├── tests/                # Integration and manual tests
└── Taskfile.yml          # Task runner configuration
```

## Common Commands

```bash
# Run the application
task run

# Run database migrations
task migrations:up

# Manual migration commands (new!)
./ark migrate up        # Run pending migrations
./ark migrate status    # Show current migration version
./ark migrate validate  # Validate schema (check tables exist)

# Create a new migration
task migrations:new name=migration_name

# Run tests
go test ./...

# Format and tidy code
task tidy
```

## Key Patterns

- **Clean Architecture**: Handlers → Services → Repositories
- **Multi-tenancy**: All queries scoped to `user_id`
- **Two-phase Auth**: ClerkAuthMiddleware → RequireAuth
- **Error Handling**: Custom error types in `internal/errs/`

See [BACKEND_GUIDE.md](../../docs/claude/BACKEND_GUIDE.md) for full documentation.
