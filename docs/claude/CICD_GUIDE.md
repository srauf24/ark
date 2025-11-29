# ARK CI/CD Guide

## GitHub Actions Workflow

ARK uses GitHub Actions for automated Docker image building and publishing to GitHub Container Registry (GHCR).

**Workflow** (`.github/workflows/build-push.yml`):
- Triggers on push to `main` branch or manual dispatch
- Builds backend and frontend images in parallel
- Tags images as `latest`
- Publishes to `ghcr.io/srauf24/ark-backend` and `ghcr.io/srauf24/ark-frontend`

**Why GitHub Actions for Building**:
- Avoids memory constraints on small EC2 instances (t3.micro has only 1GB RAM)
- Go compilation and Bun builds are memory-intensive
- Faster builds on GitHub's infrastructure
- Professional CI/CD practice for resume demonstration

**Build Arguments**:
- Frontend receives `VITE_API_URL` and `VITE_CLERK_PUBLISHABLE_KEY` at build time
- Production values: `https://api.arkcore.dev` and Clerk live publishable key

## Deployment Workflow

1. **Developer pushes to `main`**
2. **GitHub Actions builds and pushes images to GHCR**
3. **SSH into EC2 instance**
4. **Pull latest images**: `docker compose pull`
5. **Restart services**: `docker compose up -d`
6. **Verify deployment**: Check logs and health endpoints

**Future Enhancement**: Automate step 3-6 with GitHub Actions SSH deployment or AWS CodeDeploy.

## Docker Containerization

ARK uses Docker for production deployment with multi-stage builds optimized for small image sizes and security.

### Docker Images

**Backend** (`Dockerfile.backend`):
- Multi-stage build using Go 1.24 and Alpine Linux
- Binary built with CGO disabled for static linking
- Runs as non-root user (uid 1000)
- Includes database migrations
- Health check on `/health` endpoint
- Final image size: ~20MB

**Frontend** (`Dockerfile.frontend`):
- Build stage uses Bun for TypeScript compilation
- Builds shared packages (`@ark/zod`, `@ark/openapi`) first
- Vite production build with environment variable injection
- Runtime stage uses Caddy for static file serving
- SPA routing with fallback to `index.html`
- Gzip compression enabled
- Final image size: ~15MB

### Docker Compose Setup

The `docker-compose.yml` orchestrates all services for production deployment:

**Services**:
- `postgres`: PostgreSQL 16 with persistent volume and health checks
- `redis`: Redis 7 with AOF persistence and health checks
- `backend`: Go API server (pulls from GHCR)
- `frontend`: React SPA (pulls from GHCR)
- `caddy`: Reverse proxy with automatic HTTPS

**Networking**:
- All services on `ark-network` bridge network
- Only Caddy exposes ports 80/443 to host
- Internal service-to-service communication via Docker DNS

**Volumes**:
- `postgres_data`: Database persistence
- `redis_data`: Redis AOF persistence
- `caddy_data`: SSL certificates
- `caddy_config`: Caddy configuration cache

### Running with Docker Compose

```bash
# Pull latest images and start all services
docker compose pull
docker compose up -d

# View logs
docker compose logs -f

# View specific service logs
docker compose logs -f backend

# Stop all services
docker compose down

# Stop and remove volumes (WARNING: destroys data)
docker compose down -v

# Restart a specific service
docker compose restart backend

# Check service health
docker compose ps
```
