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
- `@ark/openapi` → `../../packages/openapi/src`
- `@ark/zod` → `../../packages/zod/src`

**Note**: The vite.config.ts has been updated to use `@ark/*` aliases.

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
await window.Clerk.session.getToken({ template: "api-test" })
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

## Current State

The frontend is in active development:
- **Asset Management**: Full CRUD implemented (List, Create, Edit, Delete)
- **Authentication**: Clerk integration with JWT injection and protected routes
- **Architecture**: Type-safe API client, TanStack Query, and centralized hooks
- **UI**: shadcn/ui components, responsive layout, and toast notifications
- **Pending**: Log management UI, AI integration

# Ark Frontend Guide

## Overview
The frontend for **Ark** is built with React, Vite, and Tailwind CSS. It provides a modern, responsive interface for managing assets and logs.

**Next steps for development**:
**Next steps for development**:
1. Implement Log management UI (List, Create, Edit, Delete)
2. Implement AI Query interface
3. Add E2E tests for new features

## Design Principles

# S-Tier SaaS Dashboard Design Checklist (Inspired by Stripe, Airbnb, Linear)

## I. Core Design Philosophy & Strategy

*   [ ] **Users First:** Prioritize user needs, workflows, and ease of use in every design decision.
*   [ ] **Meticulous Craft:** Aim for precision, polish, and high quality in every UI element and interaction.
*   [ ] **Speed & Performance:** Design for fast load times and snappy, responsive interactions.
*   [ ] **Simplicity & Clarity:** Strive for a clean, uncluttered interface. Ensure labels, instructions, and information are unambiguous.
*   [ ] **Focus & Efficiency:** Help users achieve their goals quickly and with minimal friction. Minimize unnecessary steps or distractions.
*   [ ] **Consistency:** Maintain a uniform design language (colors, typography, components, patterns) across the entire dashboard.
*   [ ] **Accessibility (WCAG AA+):** Design for inclusivity. Ensure sufficient color contrast, keyboard navigability, and screen reader compatibility.
*   [ ] **Opinionated Design (Thoughtful Defaults):** Establish clear, efficient default workflows and settings, reducing decision fatigue for users.

## II. Design System Foundation (Tokens & Core Components)

*   [ ] **Define a Color Palette:**
    *   [ ] **Primary Brand Color:** User-specified, used strategically.
    *   [ ] **Neutrals:** A scale of grays (5-7 steps) for text, backgrounds, borders.
    *   [ ] **Semantic Colors:** Define specific colors for Success (green), Error/Destructive (red), Warning (yellow/amber), Informational (blue).
    *   [ ] **Dark Mode Palette:** Create a corresponding accessible dark mode palette.
    *   [ ] **Accessibility Check:** Ensure all color combinations meet WCAG AA contrast ratios.
*   [ ] **Establish a Typographic Scale:**
    *   [ ] **Primary Font Family:** Choose a clean, legible sans-serif font (e.g., Inter, Manrope, system-ui).
    *   [ ] **Modular Scale:** Define distinct sizes for H1, H2, H3, H4, Body Large, Body Medium (Default), Body Small/Caption. (e.g., H1: 32px, Body: 14px/16px).
    *   [ ] **Font Weights:** Utilize a limited set of weights (e.g., Regular, Medium, SemiBold, Bold).
    *   [ ] **Line Height:** Ensure generous line height for readability (e.g., 1.5-1.7 for body text).
*   [ ] **Define Spacing Units:**
    *   [ ] **Base Unit:** Establish a base unit (e.g., 8px).
    *   [ ] **Spacing Scale:** Use multiples of the base unit for all padding, margins, and layout spacing (e.g., 4px, 8px, 12px, 16px, 24px, 32px).
*   [ ] **Define Border Radii:**
    *   [ ] **Consistent Values:** Use a small set of consistent border radii (e.g., Small: 4-6px for inputs/buttons; Medium: 8-12px for cards/modals).
*   [ ] **Develop Core UI Components (with consistent states: default, hover, active, focus, disabled):**
    *   [ ] Buttons (primary, secondary, tertiary/ghost, destructive, link-style; with icon options)
    *   [ ] Input Fields (text, textarea, select, date picker; with clear labels, placeholders, helper text, error messages)
    *   [ ] Checkboxes & Radio Buttons
    *   [ ] Toggles/Switches
    *   [ ] Cards (for content blocks, multimedia items, dashboard widgets)
    *   [ ] Tables (for data display; with clear headers, rows, cells; support for sorting, filtering)
    *   [ ] Modals/Dialogs (for confirmations, forms, detailed views)
    *   [ ] Navigation Elements (Sidebar, Tabs)
    *   [ ] Badges/Tags (for status indicators, categorization)
    *   [ ] Tooltips (for contextual help)
    *   [ ] Progress Indicators (Spinners, Progress Bars)
    *   [ ] Icons (use a single, modern, clean icon set; SVG preferred)
    *   [ ] Avatars

## III. Layout, Visual Hierarchy & Structure

*   [ ] **Responsive Grid System:** Design based on a responsive grid (e.g., 12-column) for consistent layout across devices.
*   [ ] **Strategic White Space:** Use ample negative space to improve clarity, reduce cognitive load, and create visual balance.
*   [ ] **Clear Visual Hierarchy:** Guide the user's eye using typography (size, weight, color), spacing, and element positioning.
*   [ ] **Consistent Alignment:** Maintain consistent alignment of elements.
*   [ ] **Main Dashboard Layout:**
    *   [ ] Persistent Left Sidebar: For primary navigation between modules.
    *   [ ] Content Area: Main space for module-specific interfaces.
    *   [ ] (Optional) Top Bar: For global search, user profile, notifications.
*   [ ] **Mobile-First Considerations:** Ensure the design adapts gracefully to smaller screens.

## IV. Interaction Design & Animations

*   [ ] **Purposeful Micro-interactions:** Use subtle animations and visual feedback for user actions (hovers, clicks, form submissions, status changes).
    *   [ ] Feedback should be immediate and clear.
    *   [ ] Animations should be quick (150-300ms) and use appropriate easing (e.g., ease-in-out).
*   [ ] **Loading States:** Implement clear loading indicators (skeleton screens for page loads, spinners for in-component actions).
*   [ ] **Transitions:** Use smooth transitions for state changes, modal appearances, and section expansions.
*   [ ] **Avoid Distraction:** Animations should enhance usability, not overwhelm or slow down the user.
*   [ ] **Keyboard Navigation:** Ensure all interactive elements are keyboard accessible and focus states are clear.

## V. Specific Module Design Tactics

### A. Multimedia Moderation Module

*   [ ] **Clear Media Display:** Prominent image/video previews (grid or list view).
*   [ ] **Obvious Moderation Actions:** Clearly labeled buttons (Approve, Reject, Flag, etc.) with distinct styling (e.g., primary/secondary, color-coding). Use icons for quick recognition.
*   [ ] **Visible Status Indicators:** Use color-coded Badges for content status (Pending, Approved, Rejected).
*   [ ] **Contextual Information:** Display relevant metadata (uploader, timestamp, flags) alongside media.
*   [ ] **Workflow Efficiency:**
    *   [ ] Bulk Actions: Allow selection and moderation of multiple items.
    *   [ ] Keyboard Shortcuts: For common moderation actions.
*   [ ] **Minimize Fatigue:** Clean, uncluttered interface; consider dark mode option.

### B. Data Tables Module (Contacts, Admin Settings)

*   [ ] **Readability & Scannability:**
    *   [ ] Smart Alignment: Left-align text, right-align numbers.
    *   [ ] Clear Headers: Bold column headers.
    *   [ ] Zebra Striping (Optional): For dense tables.
    *   [ ] Legible Typography: Simple, clean sans-serif fonts.
    *   [ ] Adequate Row Height & Spacing.
*   [ ] **Interactive Controls:**
    *   [ ] Column Sorting: Clickable headers with sort indicators.
    *   [ ] Intuitive Filtering: Accessible filter controls (dropdowns, text inputs) above the table.
    *   [ ] Global Table Search.
*   [ ] **Large Datasets:**
    *   [ ] Pagination (preferred for admin tables) or virtual/infinite scroll.
    *   [ ] Sticky Headers / Frozen Columns: If applicable.
*   [ ] **Row Interactions:**
    *   [ ] Expandable Rows: For detailed information.
    *   [ ] Inline Editing: For quick modifications.
    *   [ ] Bulk Actions: Checkboxes and contextual toolbar.
    *   [ ] Action Icons/Buttons per Row: (Edit, Delete, View Details) clearly distinguishable.

### C. Configuration Panels Module (Microsite, Admin Settings)

*   [ ] **Clarity & Simplicity:** Clear, unambiguous labels for all settings. Concise helper text or tooltips for descriptions. Avoid jargon.
*   [ ] **Logical Grouping:** Group related settings into sections or tabs.
*   [ ] **Progressive Disclosure:** Hide advanced or less-used settings by default (e.g., behind "Advanced Settings" toggle, accordions).
*   [ ] **Appropriate Input Types:** Use correct form controls (text fields, checkboxes, toggles, selects, sliders) for each setting.
*   [ ] **Visual Feedback:** Immediate confirmation of changes saved (e.g., toast notifications, inline messages). Clear error messages for invalid inputs.
*   [ ] **Sensible Defaults:** Provide default values for all settings.
*   [ ] **Reset Option:** Easy way to "Reset to Defaults" for sections or entire configuration.
*   [ ] **Microsite Preview (If Applicable):** Show a live or near-live preview of microsite changes.

## VI. CSS & Styling Architecture

*   [ ] **Choose a Scalable CSS Methodology:**
    *   [ ] **Utility-First (Recommended for LLM):** e.g., Tailwind CSS. Define design tokens in config, apply via utility classes.
    *   [ ] **BEM with Sass:** If not utility-first, use structured BEM naming with Sass variables for tokens.
    *   [ ] **CSS-in-JS (Scoped Styles):** e.g., Stripe's approach for Elements.
*   [ ] **Integrate Design Tokens:** Ensure colors, fonts, spacing, radii tokens are directly usable in the chosen CSS architecture.
*   [ ] **Maintainability & Readability:** Code should be well-organized and easy to understand.
*   [ ] **Performance:** Optimize CSS delivery; avoid unnecessary bloat.

## VII. General Best Practices

*   [ ] **Iterative Design & Testing:** Continuously test with users and iterate on designs.
*   [ ] **Clear Information Architecture:** Organize content and navigation logically.
*   [ ] **Responsive Design:** Ensure the dashboard is fully functional and looks great on all device sizes (desktop, tablet, mobile).
*   [ ] **Documentation:** Maintain clear documentation for the design system and components.