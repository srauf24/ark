# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is the React frontend for Ark, a homelab asset tracking and configuration log management application. The frontend is part of a Turborepo monorepo and integrates with a Go backend (Echo framework) using type-safe contracts.

**Tech Stack**: React 19, TypeScript, Vite, TailwindCSS v4, Clerk (auth), TanStack Query, React Router v7, ts-rest, shadcn/ui

## Development Commands

### Frontend-specific (from apps/frontend directory)
```bash
# Start development server (waits for OpenAPI contracts to build first)
bun dev

# Build for production (requires .env.local file)
bun build

# Testing
bun test              # Run unit/component tests
bun test:e2e          # Run E2E tests (headless)
bun test:e2e:ui       # Run E2E tests (UI mode)
bun test:e2e:debug    # Run E2E tests (debug mode)
bun test:e2e:report   # View E2E test report

# Linting
bun lint              # Check for issues
bun lint:fix          # Auto-fix issues

# Formatting
bun format            # Check formatting
bun format:fix        # Auto-fix formatting

# Type checking
bun typecheck

# Clean build artifacts
bun clean
```
### Design Principles
design-principles.md


### Monorepo commands (from root directory)
```bash
# Start all services (backend + frontend)
bun dev

# Build all packages
bun build

# Lint all packages
bun lint

# Type check all packages
bun typecheck
```

## Architecture

### Project Structure

```
src/
├── api/              # API client setup and utilities
│   ├── index.ts      # API client initialization with ts-rest
│   ├── types.ts      # Type exports from contracts
│   └── utils.ts      # Query keys for TanStack Query
├── config/           # Configuration and environment variables
│   └── env.ts        # Zod-validated environment schema
├── lib/              # Shared utilities
│   └── utils.ts      # Tailwind class merging (cn function)
├── App.tsx           # Root component
├── main.tsx          # Application entry point
└── index.css         # Global styles (Tailwind directives)
```

### Key Architectural Patterns

**Type-Safe API Client**:
- Uses `@ts-rest/core` to create a type-safe API client from shared contracts
- Located in `src/api/index.ts` with the `useApiClient` hook
- Automatically injects Clerk JWT tokens via `Authorization: Bearer <token>`
- Retries failed 401 requests up to 2 times (token refresh handling)
- Supports blob responses via `isBlob` parameter for file downloads

**Authentication Flow**:
- Clerk React SDK (`@clerk/clerk-react`) for authentication
- Custom JWT template named "custom" used for API requests
- Token automatically refreshed and injected into all API calls
- Backend validates JWT and enforces user_id scoping

**Environment Configuration**:
- Zod schema validation in `src/config/env.ts`
- All environment variables prefixed with `VITE_`
- Validates on application startup with human-readable error messages
- Uses Zod v4's `treeifyError` for better error display

**Shared Contracts**:
- API contracts from `@ark/openapi` workspace package
- Zod schemas from `@ark/zod` workspace package
- Full type safety between frontend and backend
- Contracts generated from OpenAPI specification

**Styling System**:
- TailwindCSS v4 with Vite plugin
- shadcn/ui components (New York style)
- `cn()` utility for conditional class merging (clsx + tailwind-merge)
- CSS variables for theming (neutral base color)
- Lucide React for icons

## Configuration

### Environment Variables

Create a `.env.local` file with the following required variables:

```bash
# Clerk Authentication (required)
VITE_CLERK_PUBLISHABLE_KEY=pk_test_...

# Backend API URL (default: http://localhost:3000)
VITE_API_URL=http://localhost:8080

# Environment (default: local)
VITE_ENV=local  # or "development" | "production"
```

See `env.local.sample` for reference.

**Important**: The build script (`bun build`) uses `env-cmd` to load `.env.local`, so this file is required for production builds.

### Path Aliases

The project uses TypeScript path aliases for cleaner imports:

```typescript
import { Button } from "@/components/ui/button"
import { cn } from "@/lib/utils"
import { API_URL } from "@/config/env"
```

**Configured aliases**:
- `@/*` → `./src/*`
- `@gardenjournal/openapi` → `../../packages/openapi/src` (legacy, should migrate to `@ark/openapi`)
- `@gardenjournal/zod` → `../../packages/zod/src` (legacy, should migrate to `@ark/zod`)

**Note**: The vite.config.ts still references `@gardenjournal/*` aliases from the previous codebase. These should be updated to `@ark/*` to match package.json dependencies.

### shadcn/ui Configuration

The project uses shadcn/ui components with the following configuration (in `components.json`):

- **Style**: New York
- **Base color**: Neutral
- **CSS Variables**: Enabled
- **Icon Library**: Lucide React
- **TypeScript**: Enabled

Add new components with:
```bash
npx shadcn@latest add <component-name>
```

## Integration with Backend

### API Communication

The frontend communicates with the backend using a type-safe client based on ts-rest:

**Example usage**:
```typescript
const apiClient = useApiClient();

// Type-safe API call
const response = await apiClient.assets.getAll({
  query: { page: 1, limit: 10 }
});

// Response is fully typed from backend contracts
if (response.status === 200) {
  const assets = response.body.data;
}
```

**Key features**:
- Automatic JWT injection from Clerk
- Type inference from backend contracts
- Built-in retry logic for 401 errors
- Support for blob responses (file downloads)

### Contract Updates

When backend API changes:

1. Backend updates Zod schemas in `packages/zod/`
2. Backend updates contracts in `packages/openapi/src/contracts/`
3. Run `bun run gen` from `packages/openapi/` to regenerate spec
4. Frontend automatically gets updated types on next build
5. TypeScript will catch any breaking changes

**Important**: The frontend dev server waits for OpenAPI contracts to build (`wait-on ../../packages/openapi/dist/index.js`) before starting.

## Common Development Tasks

### Adding a new API integration

1. Ensure backend contract exists in `packages/openapi/src/contracts/`
2. Import the API client in your component:
```typescript
import { useApiClient } from "@/api";

function MyComponent() {
  const apiClient = useApiClient();
  // Use apiClient.yourEndpoint.method(...)
}
```
3. Use TanStack Query for data fetching (recommended):
```typescript
import { useQuery } from "@tanstack/react-query";
import { useApiClient } from "@/api";

function MyComponent() {
  const apiClient = useApiClient();

  const { data, isLoading } = useQuery({
    queryKey: ["assets", "list"],
    queryFn: () => apiClient.assets.getAll({ query: {} })
  });
}
```

### Adding a new shadcn/ui component

```bash
# Navigate to frontend directory
cd apps/frontend

# Add the component
npx shadcn@latest add button

# Component will be added to src/components/ui/
```

### Working with forms

The project is set up for React Hook Form with Zod validation:

```typescript
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const schema = z.object({
  name: z.string().min(1)
});

function MyForm() {
  const form = useForm({
    resolver: zodResolver(schema)
  });

  // Use form...
}
```

### Updating legacy package references

The vite.config.ts still uses `@gardenjournal/*` aliases. To update:

1. Replace `@gardenjournal/openapi` with `@ark/openapi`
2. Replace `@gardenjournal/zod` with `@ark/zod`
3. Update any import statements in source files
4. Run `bun typecheck` to verify

## Code Style and Conventions

### Linting and Formatting

- **ESLint**: TypeScript ESLint with React Hooks and React Refresh plugins
- **Prettier**: Configured with import sorting via `@trivago/prettier-plugin-sort-imports`
  - Double quotes (not single quotes)
  - Trailing commas
  - Auto end-of-line
- **TypeScript**: Strict mode enabled

### Import Organization

Prettier automatically sorts imports. The general order is:
1. React and React-related imports
2. Third-party libraries
3. Internal packages (@ark/*, @/*)
4. Relative imports
5. CSS/styles

### Component Patterns

- Use functional components with hooks
- Prefer named exports for components
- Use TypeScript for all new code
- Extract reusable logic into custom hooks
- Use `cn()` utility for conditional Tailwind classes

## Testing Strategy

### Unit/Component Tests

The project uses Vite's test runner (Vitest expected):

```bash
# Run unit/component tests
bun test

# Run tests in watch mode
bun test --watch

# Run tests with coverage
bun test --coverage
```

**Testing guidelines**:
- Test user interactions and component behavior
- Mock API calls using the TanStack Query testing utilities
- Test form validation and submission
- Test authentication flows with Clerk mocks

### End-to-End Tests (Playwright)

The project uses Playwright for E2E testing across multiple browsers:

```bash
# Run all E2E tests (headless mode)
bun test:e2e

# Run tests in UI mode (interactive)
bun test:e2e:ui

# Run tests in debug mode (step-through)
bun test:e2e:debug

# View HTML test report
bun test:e2e:report
```

**E2E Test Structure**:
- Tests are located in the `e2e/` directory
- Configuration in `playwright.config.ts`
- Automatically starts dev server before tests (configured in webServer)
- Tests run against `http://localhost:3000` by default

**Playwright Configuration**:
- Tests run on Chromium, Firefox, and WebKit by default
- Parallel execution enabled (can be disabled with `--workers=1`)
- Automatic retries on CI (2 retries)
- Screenshots captured on failure
- Traces collected on first retry

**Writing E2E Tests**:
```typescript
import { test, expect } from "@playwright/test";

test.describe("Feature Name", () => {
  test("should do something", async ({ page }) => {
    await page.goto("/");

    // Your test assertions
    await expect(page.getByRole("heading")).toBeVisible();
  });
});
```

**Best Practices**:
- Use semantic selectors (role, label, text) over CSS selectors
- Write tests from the user's perspective
- Test critical user journeys (auth, CRUD operations)
- Mock external API calls when needed
- Keep tests independent and isolated

## Troubleshooting

### Development server won't start

```bash
# Ensure OpenAPI contracts are built first
cd ../../packages/openapi
bun run build

# Then start the frontend dev server
cd ../../apps/frontend
bun dev
```

### TypeScript errors after backend changes

```bash
# Rebuild OpenAPI contracts
cd ../../packages/openapi
bun run build

# Rebuild Zod schemas
cd ../zod
bun run build

# Type check frontend
cd ../../apps/frontend
bun typecheck
```

### Environment variable errors on startup

The application validates environment variables with Zod on startup. If you see errors:

1. Check that `.env.local` exists and has all required variables
2. Ensure `VITE_CLERK_PUBLISHABLE_KEY` is valid
3. Verify `VITE_API_URL` is a valid URL
4. Check for typos in variable names

### Clerk authentication issues

```bash
# Verify Clerk publishable key is correct
echo $VITE_CLERK_PUBLISHABLE_KEY

# Check that backend is configured with matching Clerk secret
# Backend needs ARK_AUTH.CLERK.SECRET_KEY and ARK_AUTH.CLERK.JWT_ISSUER

# Test token generation in browser console:
await window.Clerk.session.getToken({ template: "custom" })
```

### Build failures

```bash
# Clean and reinstall dependencies
bun clean
rm -rf node_modules
bun install

# Rebuild workspace dependencies
cd ../../packages/zod && bun run build
cd ../openapi && bun run build

# Try building again
cd ../../apps/frontend
bun build
```

## Prerequisites

- **Node.js**: 22+ (for frontend development)
- **Bun**: Latest version (package manager)
- **Backend**: Running on configured `VITE_API_URL` (default: http://localhost:8080)
- **Clerk Account**: For authentication (get publishable key from Clerk dashboard)

## Current State

The frontend is in early development stages:
- Basic React + Vite + TypeScript setup complete
- API client and authentication infrastructure ready
- shadcn/ui and TailwindCSS configured
- Environment validation implemented
- No UI components or routes implemented yet

**Next steps for development**:
1. Update legacy `@gardenjournal/*` references to `@ark/*`
2. Implement routing with React Router v7
3. Create layouts and page components
4. Add asset and log management UI
5. Implement TanStack Query hooks for data fetching
