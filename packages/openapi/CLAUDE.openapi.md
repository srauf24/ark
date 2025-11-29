# OpenAPI Package Quick Reference

## Purpose

Generates OpenAPI specification from Zod schemas and ts-rest contracts.

## Workflow

1. **Update Zod schemas** in `packages/zod/src/`
2. **Update contracts** in `packages/openapi/src/contracts/`
3. **Generate spec**: `bun run gen`
4. **Output**: `apps/backend/static/openapi.json`

## Commands

```bash
# Generate OpenAPI specification
bun run gen

# Build TypeScript contracts
bun run build
```

## Important Notes

- Frontend dev server waits for OpenAPI contracts to build
- Always run `bun gen` after changing contracts
- Verify spec at `http://localhost:8080/openapi.json`
- Backend route registration must match contracts exactly

## Verification

```bash
# Test that spec matches backend routes
go test -v ./internal/handler -run TestOpenAPI
```

See [DEV_GUIDE.md](../../docs/claude/DEV_GUIDE.md#working-with-the-openapi-spec) for full workflow.
