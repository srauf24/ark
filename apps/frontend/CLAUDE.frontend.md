# Frontend Quick Reference

## File Structure

```
apps/frontend/
├── src/
│   ├── api/              # Type-safe API client (ts-rest)
│   ├── components/       # React components
│   ├── pages/            # Page components
│   ├── hooks/            # Custom React hooks
│   ├── lib/              # Utilities
│   └── config/           # Configuration
├── e2e/                  # Playwright E2E tests
└── public/               # Static assets
```

## Common Commands

```bash
# Start development server
bun dev

# Run tests
bun test                  # Unit/component tests (Vitest)
bun test:e2e              # E2E tests (Playwright)

# Build for production
bun build

# Type check
bun typecheck
```

## Key Patterns

- **Environment Variables**: Use `import.meta.env.VITE_*`, NOT `process.env`
- **API Client**: `useApiClient()` hook with full type safety
- **Authentication**: Clerk provider with JWT template "api-test"
- **State Management**: TanStack Query for server state

See [FRONTEND_GUIDE.md](../../docs/claude/FRONTEND_GUIDE.md) for full documentation.
