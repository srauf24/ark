# Zod Package Quick Reference

## Purpose

Shared Zod schemas for request/response validation across frontend and backend.

## Schema Patterns

- **Request DTOs**: Validation for API inputs
- **Response DTOs**: Type definitions for API outputs
- **Domain Models**: Asset, AssetLog types
- **Validation Rules**: Max lengths, required fields, formats

## Commands

```bash
# Build shared Zod schemas and TypeScript types
bun run build
```

## Usage

### In Backend (Go)
- Schemas define OpenAPI spec
- Validation happens in handlers

### In Frontend (TypeScript)
- Import types from `@ark/zod`
- Use with react-hook-form for form validation
- Type-safe API client uses these types

## Important Notes

- Changes to schemas require rebuilding both `zod` and `openapi` packages
- Frontend depends on built output, not source files
- Always run `bun run build` after schema changes

## Related Documentation

- [`docs/claude/BACKEND_GUIDE.md`](../../docs/claude/BACKEND_GUIDE.md) - **Validation**: How schemas are used in the backend
- [`docs/claude/FRONTEND_GUIDE.md`](../../docs/claude/FRONTEND_GUIDE.md) - **Form Types**: Using schemas with react-hook-form

See [DEV_GUIDE.md](../../docs/claude/DEV_GUIDE.md#shared-packages) for full workflow.
