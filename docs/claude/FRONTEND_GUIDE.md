# ARK Frontend Guide

## Tech Stack
- **React 19**: Latest React with concurrent features
- **TypeScript**: Strict mode enabled for type safety
- **Vite 7**: Ultra-fast build tool and dev server
- **TailwindCSS v4**: Utility-first CSS with Vite plugin
- **Clerk**: Authentication provider (`@clerk/clerk-react`)
- **TanStack Query**: Server state management and data fetching
- **React Router v7**: Client-side routing
- **ts-rest**: Type-safe API client from OpenAPI contracts
- **shadcn/ui**: Accessible component library (New York style, Lucide icons)
- **react-hook-form**: Form state management with Zod validation
- **date-fns**: Date formatting and manipulation
- **Playwright**: E2E testing framework
- **Vitest**: Unit and component testing

## Environment Variables (Frontend)

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

## Type-Safe API Client

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

## Implemented Features

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

## Frontend Testing

### Unit/Component Tests (Vitest)
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

### E2E Tests (Playwright)
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
