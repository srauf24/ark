# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with the Ark project.

## How to use this repository documentation

**IMPORTANT**: To avoid exceeding model context limits, DO NOT load all project documentation at once.

When generating code or plans:

1. **Load only the specific guide you need** from the docs below
2. Follow that guide as the single source of truth for that domain
3. If a user request spans multiple domains, load those specific guides

## Documentation Index

### High-Level Architecture / Concepts
- [`docs/claude/ARCHITECTURE.md`](docs/claude/ARCHITECTURE.md) - Project overview, goals, clean architecture, domain models, security
- [`docs/claude/DEV_GUIDE.md`](docs/claude/DEV_GUIDE.md) - Development commands, database management, common tasks, troubleshooting

### Frontend (React, Vite, Clerk, TanStack Query)
- [`apps/frontend/CLAUDE.frontend.md`](apps/frontend/CLAUDE.frontend.md) - Quick reference
- [`docs/claude/FRONTEND_GUIDE.md`](docs/claude/FRONTEND_GUIDE.md) - Tech stack, environment variables, API client, features, testing

### Backend (Go, Echo, PostgreSQL, Clean Architecture)
- [`apps/backend/CLAUDE.backend.md`](apps/backend/CLAUDE.backend.md) - Quick reference
- [`docs/claude/BACKEND_GUIDE.md`](docs/claude/BACKEND_GUIDE.md) - API endpoints, configuration, infrastructure, testing

### Contracts / OpenAPI / Zod
- [`packages/openapi/CLAUDE.openapi.md`](packages/openapi/CLAUDE.openapi.md) - Quick reference
- [`packages/zod/CLAUDE.zod.md`](packages/zod/CLAUDE.zod.md) - Quick reference

### CI/CD, Docker, Deployment, AWS
- [`docs/claude/CICD_GUIDE.md`](docs/claude/CICD_GUIDE.md) - GitHub Actions, Docker builds, deployment workflow
- [`docs/claude/DEPLOYMENT_GUIDE.md`](docs/claude/DEPLOYMENT_GUIDE.md) - AWS EC2, Caddy, production setup, troubleshooting
- [`docs/claude/OBSERVABILITY_GUIDE.md`](docs/claude/OBSERVABILITY_GUIDE.md) - Logging, New Relic, metrics, future enhancements

## Rules

- Prefer small, incremental changes
- Follow the project's clean architecture patterns
- Maintain type safety across frontend ↔ backend
- Match OpenAPI contracts exactly
- Never assume hidden state; refer to the docs above

## Task Workflow

Whenever I provide a JIRA ticket:

1. Read only the docs relevant to that ticket
2. Produce a step-by-step implementation plan first
3. Wait for approval before generating code
4. Generate minimal code changes following the plan
5. Add/update tests as needed

## Quick Start

### Prerequisites

**Local Development**:
- Go 1.24+, Node.js 22+, Bun, PostgreSQL 16+, Redis 8+
- Task runner (`brew install go-task`)

**Production Deployment**:
- Docker 20.10+ and Docker Compose v2
- AWS EC2 (t3.micro or larger)
- Domain with Cloudflare DNS (optional)

### Development Commands

```bash
# Install dependencies
bun install

# Start all services (backend + frontend)
bun dev

# Backend only
cd apps/backend && task run

# Frontend only
cd apps/frontend && bun dev

# Run tests
go test ./...                    # Backend
bun test                         # Frontend unit tests
bun test:e2e                     # Frontend E2E tests
```

See [`docs/claude/DEV_GUIDE.md`](docs/claude/DEV_GUIDE.md) for complete command reference.

## Project Structure

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
├── docker-compose.yml    # Production deployment
├── Dockerfile.backend    # Backend container
├── Dockerfile.frontend   # Frontend container
├── Caddyfile             # Reverse proxy config
└── CLAUDE.md             # This file
```

## Live Production

- **Frontend**: https://arkcore.dev
- **Backend API**: https://api.arkcore.dev
- **Infrastructure**: AWS EC2 t3.micro with Docker Compose
- **Database**: PostgreSQL 16 with pgvector
- **CI/CD**: GitHub Actions → GHCR → EC2

See [`docs/claude/DEPLOYMENT_GUIDE.md`](docs/claude/DEPLOYMENT_GUIDE.md) for deployment details.
