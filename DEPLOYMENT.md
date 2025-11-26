# Ark Deployment Guide

This guide covers deploying Ark to various platforms. Choose the option that best fits your needs.

## Prerequisites

All deployment options require:
- **Clerk Account**: Sign up at https://clerk.com/ for authentication
- **Resend Account**: Sign up at https://resend.com/ for transactional emails
- **PostgreSQL 16+**: Database (provided by platform or self-hosted)
- **Redis 8+**: Cache and job queue (provided by platform or self-hosted)

## Quick Start: Docker Compose (Self-Hosted)

Perfect for running Ark in your homelab on any server, NAS, or Raspberry Pi.

### 1. Clone and Configure

```bash
# Clone the repository
git clone <your-repo-url>
cd ark

# Create environment file
cp .env.docker.example .env

# Edit .env with your actual values
nano .env
```

### 2. Required Environment Variables

Edit `.env` and set these **required** values:

```bash
# Clerk Authentication (from https://dashboard.clerk.com/)
CLERK_SECRET_KEY=sk_live_...
CLERK_PUBLISHABLE_KEY=pk_live_...
CLERK_JWT_ISSUER=https://your-app.clerk.accounts.dev

# Auth Secret (generate random 32 chars)
ARK_AUTH_SECRET_KEY=your-random-secret-here

# Resend API (from https://resend.com/api-keys)
RESEND_API_KEY=re_...

# Database Password (choose a strong password)
POSTGRES_PASSWORD=your-secure-password
```

### 3. Deploy

```bash
# Start all services
docker compose up -d

# View logs
docker compose logs -f

# Check health
curl http://localhost:8080/health
```

### 4. Run Database Migrations

```bash
# Access backend container
docker compose exec backend sh

# Run migrations (install tern first)
go install github.com/jackc/tern/v2@latest
tern migrate -m ./internal/database/migrations

# Exit container
exit
```

### 5. Access Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Docs**: http://localhost:8080/openapi.json

### Updating

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker compose down
docker compose up -d --build
```

---

## Platform Deployment Options

### Option 1: Railway (Recommended)

**Best for**: Quick deployment with managed databases

#### Steps:

1. **Create Railway Account**: https://railway.app/

2. **Create New Project** → Deploy from GitHub

3. **Add Services**:
   - Add PostgreSQL database
   - Add Redis database
   - Railway will auto-detect backend (Go) and frontend (Node.js)

4. **Configure Backend Environment Variables**:
   ```
   ARK_PRIMARY.ENV=production
   ARK_SERVER.CORS_ALLOWED_ORIGINS=https://your-frontend-url.railway.app
   ARK_DATABASE.HOST=${{Postgres.PGHOST}}
   ARK_DATABASE.PORT=${{Postgres.PGPORT}}
   ARK_DATABASE.USER=${{Postgres.PGUSER}}
   ARK_DATABASE.PASSWORD=${{Postgres.PGPASSWORD}}
   ARK_DATABASE.NAME=${{Postgres.PGDATABASE}}
   ARK_DATABASE.SSL_MODE=require
   ARK_AUTH.SECRET_KEY=<random-secret>
   ARK_AUTH.CLERK.SECRET_KEY=<clerk-secret>
   ARK_AUTH.CLERK.JWT_ISSUER=<clerk-issuer>
   ARK_INTEGRATION.RESEND_API_KEY=<resend-key>
   ARK_REDIS.ADDRESS=${{Redis.REDIS_URL}}
   ARK_OBSERVABILITY.LOGGING.LEVEL=info
   ARK_OBSERVABILITY.LOGGING.FORMAT=json
   ARK_OBSERVABILITY.NEW_RELIC.LICENSE_KEY=<optional>
   ARK_OBSERVABILITY.HEALTH_CHECKS.INTERVAL=30s
   ARK_OBSERVABILITY.HEALTH_CHECKS.TIMEOUT=5s
   ```

5. **Configure Frontend Environment Variables**:
   ```
   VITE_API_URL=https://your-backend-url.railway.app
   VITE_CLERK_PUBLISHABLE_KEY=<clerk-publishable>
   VITE_ENV=production
   ```

6. **Deploy**: Railway auto-deploys on git push

#### Cost:
- Free tier: $5 credit/month
- Hobby: ~$5-15/month
- Pro: $20+/month

---

### Option 2: Render

**Best for**: Free tier testing, simple deployments

#### Steps:

1. **Create Render Account**: https://render.com/

2. **Create Services**:
   - New PostgreSQL database (free tier available)
   - New Redis instance (free 25MB)
   - New Web Service (backend) → Dockerfile: `apps/backend/Dockerfile`
   - New Static Site (frontend) → Build: `cd apps/frontend && bun install && bun build`

3. **Set Environment Variables** (similar to Railway, adjust database URLs for Render's format)

4. **Configure Build Commands**:
   - Backend: Auto-detected from Dockerfile
   - Frontend: `cd apps/frontend && bun run build`
   - Frontend Publish: `apps/frontend/dist`

#### Cost:
- Free tier: Available (databases expire after 90 days)
- Starter: $7/month per service

---

### Option 3: Fly.io

**Best for**: Global edge deployment, Docker-based apps

#### Prerequisites:
- Install flyctl: `curl -L https://fly.io/install.sh | sh`

#### Steps:

1. **Login**: `fly auth login`

2. **Create App**:
   ```bash
   cd apps/backend
   fly launch --no-deploy
   ```

3. **Add PostgreSQL**:
   ```bash
   fly postgres create --name ark-db
   fly postgres attach --app <your-app-name> ark-db
   ```

4. **Add Redis** (via Upstash):
   ```bash
   fly redis create --name ark-redis
   ```

5. **Set Secrets**:
   ```bash
   fly secrets set \
     ARK_AUTH.CLERK.SECRET_KEY=sk_... \
     ARK_AUTH.CLERK.JWT_ISSUER=https://... \
     ARK_INTEGRATION.RESEND_API_KEY=re_... \
     ARK_AUTH.SECRET_KEY=random-secret
   ```

6. **Deploy**:
   ```bash
   fly deploy
   ```

7. **Deploy Frontend**:
   ```bash
   cd apps/frontend
   fly launch --no-deploy
   fly deploy
   ```

#### Cost:
- Free tier: 3 VMs with 256MB RAM
- Paid: ~$5-20/month

---

### Option 4: DigitalOcean App Platform

**Best for**: Managed infrastructure with droplet flexibility

#### Steps:

1. **Create DigitalOcean Account**: https://digitalocean.com/

2. **Create App** → Deploy from GitHub

3. **Add Components**:
   - Database: Managed PostgreSQL
   - Database: Managed Redis
   - Service: Backend (Dockerfile)
   - Static Site: Frontend

4. **Configure Environment Variables** (similar to Railway)

5. **Deploy**: Auto-deploys on push

#### Cost:
- Basic: $5/month (512MB)
- Professional: $12+/month (1GB+)
- Managed DB: $15+/month

---

### Option 5: Self-Hosted VPS

**Best for**: Full control, cost-effective at scale

#### Providers:
- **DigitalOcean**: $6/month (1GB RAM)
- **Linode/Akamai**: $5/month (1GB RAM)
- **Hetzner**: €4/month (2GB RAM, best value)
- **Vultr**: $6/month (1GB RAM)

#### Setup:

1. **Create VPS** (Ubuntu 24.04 LTS recommended)

2. **SSH into server**:
   ```bash
   ssh root@your-server-ip
   ```

3. **Install Docker**:
   ```bash
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh
   ```

4. **Clone repository**:
   ```bash
   git clone <your-repo> /opt/ark
   cd /opt/ark
   ```

5. **Configure environment**:
   ```bash
   cp .env.docker.example .env
   nano .env
   ```

6. **Deploy**:
   ```bash
   docker compose up -d
   ```

7. **Setup Nginx reverse proxy** (optional, for HTTPS):
   ```bash
   apt install nginx certbot python3-certbot-nginx
   ```

8. **Configure domain and SSL**:
   ```bash
   certbot --nginx -d yourdomain.com
   ```

---

## Environment Variables Reference

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `ARK_AUTH.CLERK.SECRET_KEY` | Clerk backend secret | `sk_live_...` |
| `ARK_AUTH.CLERK.JWT_ISSUER` | Clerk JWT issuer URL | `https://app.clerk.accounts.dev` |
| `CLERK_PUBLISHABLE_KEY` | Clerk frontend key | `pk_live_...` |
| `ARK_INTEGRATION.RESEND_API_KEY` | Resend email API key | `re_...` |
| `ARK_AUTH.SECRET_KEY` | Session secret (32+ chars) | Random string |
| `ARK_OBSERVABILITY.LOGGING.LEVEL` | Log level | `info` or `debug` |
| `ARK_OBSERVABILITY.LOGGING.FORMAT` | Log format | `json` or `console` |
| `ARK_OBSERVABILITY.NEW_RELIC.LICENSE_KEY` | New Relic key (can be empty) | `...` or empty |
| `ARK_OBSERVABILITY.HEALTH_CHECKS.INTERVAL` | Health check interval | `30s` |
| `ARK_OBSERVABILITY.HEALTH_CHECKS.TIMEOUT` | Health check timeout | `5s` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEW_RELIC_LICENSE_KEY` | New Relic APM key | Empty (disabled) |
| `LOG_LEVEL` | Logging verbosity | `info` |
| `REDIS_PASSWORD` | Redis password | Empty |
| `CORS_ORIGINS` | Allowed CORS origins | `http://localhost:3000` |

### Database Variables (Docker Compose)

These are auto-configured in docker-compose.yml but can be overridden:

```bash
ARK_DATABASE.HOST=postgres
ARK_DATABASE.PORT=5432
ARK_DATABASE.USER=ark_user
ARK_DATABASE.PASSWORD=${POSTGRES_PASSWORD}
ARK_DATABASE.NAME=ark
ARK_DATABASE.SSL_MODE=disable
```

---

## Troubleshooting

### Config Validation Errors

**Error**: `Key: 'Config.Auth.Clerk.SecretKey' Error:Field validation for 'SecretKey' failed on the 'required' tag`

**Solution**: Ensure ALL required environment variables are set:
```bash
# Check which vars are missing
docker compose config

# Verify .env file
cat .env | grep CLERK
```

### Database Connection Issues

**Error**: `failed to connect to database`

**Solution**:
```bash
# Check if database is running
docker compose ps

# Check database logs
docker compose logs postgres

# Verify credentials
docker compose exec postgres psql -U ark_user -d ark -c "SELECT 1;"
```

### Redis Connection Issues

**Error**: `failed to connect to redis`

**Solution**:
```bash
# Check Redis status
docker compose exec redis redis-cli ping

# If using password
docker compose exec redis redis-cli -a ${REDIS_PASSWORD} ping
```

### Migration Errors

**Error**: `relation "assets" does not exist`

**Solution**: Run database migrations:
```bash
# Install tern in backend container
docker compose exec backend sh
go install github.com/jackc/tern/v2@latest
tern migrate -m ./internal/database/migrations
```

### CORS Errors

**Error**: `blocked by CORS policy`

**Solution**: Update backend CORS configuration:
```bash
# In .env or platform config
ARK_SERVER.CORS_ALLOWED_ORIGINS=https://your-frontend-domain.com,http://localhost:3000
```

### Authentication 401 Errors

**Error**: `unauthorized: user not authenticated`

**Solutions**:
1. Verify Clerk keys match same project:
   - Backend uses `CLERK_SECRET_KEY` (sk_...)
   - Frontend uses `CLERK_PUBLISHABLE_KEY` (pk_...)

2. Check JWT issuer matches:
   ```bash
   # Decode JWT at https://jwt.io
   # Verify 'iss' claim matches ARK_AUTH.CLERK.JWT_ISSUER
   ```

3. Ensure JWT template exists in Clerk dashboard:
   - Go to Clerk Dashboard → JWT Templates
   - Create template named "api-test"

---

## Platform Comparison

| Platform | Ease of Use | Cost | Control | Best For |
|----------|-------------|------|---------|----------|
| **Railway** | ⭐⭐⭐⭐⭐ | $$ | Medium | Quick deploys |
| **Render** | ⭐⭐⭐⭐⭐ | $ | Medium | Free testing |
| **Fly.io** | ⭐⭐⭐⭐ | $$ | High | Global edge |
| **DigitalOcean** | ⭐⭐⭐⭐ | $$$ | Medium | Managed services |
| **Docker Compose** | ⭐⭐⭐ | $ | Highest | Homelab/self-host |
| **VPS** | ⭐⭐⭐ | $ | Highest | Cost-effective scale |

## Recommended Path

**For Learning/Testing**: Render (free tier) or Railway (generous free tier)

**For Homelab Use**: Docker Compose on your existing server/NAS

**For Production**: Railway or Fly.io (managed databases + easy deploys)

**For Cost Optimization**: Hetzner VPS with Docker Compose (~€4/month)

---

## Next Steps

1. Choose your deployment platform
2. Set up Clerk and Resend accounts
3. Configure environment variables
4. Deploy using platform-specific instructions
5. Run database migrations
6. Access your Ark instance!

Need help? Check the main README.md or open an issue.
