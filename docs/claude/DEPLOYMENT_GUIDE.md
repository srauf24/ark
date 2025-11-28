# ARK Deployment Guide

## AWS EC2 Deployment

### Infrastructure Setup

**Instance Type**: AWS EC2 t3.micro (1 vCPU, 1GB RAM)
- Free tier eligible for 12 months
- Sufficient for Docker Compose stack with resource limits
- Cost: ~$7-10/month after free tier

**Operating System**: Ubuntu 22.04 LTS
- Long-term support and stability
- Docker and Docker Compose compatibility
- Familiar package management (apt)

**Security Group Configuration**:
```
Inbound Rules:
- Port 22 (SSH): Your IP only (for security)
- Port 80 (HTTP): 0.0.0.0/0 (Caddy redirects to HTTPS)
- Port 443 (HTTPS): 0.0.0.0/0 (public access)

Outbound Rules:
- All traffic: 0.0.0.0/0 (for package updates, Docker pulls)
```

**IAM Best Practices**:
- Create dedicated IAM user for EC2 management
- Use SSH key pairs (never password authentication)
- Rotate keys periodically
- Enable CloudWatch monitoring (optional)

### Server Setup

```bash
# 1. SSH into EC2 instance
ssh -i ~/.ssh/ark-key.pem ubuntu@your-ec2-ip

# 2. Update system packages
sudo apt update && sudo apt upgrade -y

# 3. Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker ubuntu

# 4. Install Docker Compose
sudo apt install docker-compose-plugin -y

# 5. Verify installations
docker --version
docker compose version

# 6. Clone repository or copy docker-compose.yml
mkdir -p ~/ark
cd ~/ark
# Copy docker-compose.yml, Caddyfile, and .env files

# 7. Create .env file with production secrets
nano .env
# Add ARK_DATABASE.PASSWORD, ARK_AUTH.CLERK.SECRET_KEY, etc.

# 8. Start services
docker compose up -d

# 9. Check logs
docker compose logs -f
```

### Domain and DNS Configuration

**Domain**: `arkcore.dev` (owned by Samee)

**DNS Records** (Cloudflare):
```
Type    Name    Content             Proxy Status
A       @       <EC2-IP>            Proxied (orange cloud)
A       api     <EC2-IP>            Proxied (orange cloud)
CNAME   www     arkcore.dev         Proxied (orange cloud)
```

**Cloudflare Settings**:
- SSL/TLS Mode: Full (strict) - Caddy handles SSL termination
- Always Use HTTPS: Enabled
- Automatic HTTPS Rewrites: Enabled
- Proxy Status: Enabled (DDoS protection, caching, WAF)

**Important**: If experiencing 502 errors, check Cloudflare proxy settings and ensure Caddy is properly serving HTTPS.

### Caddy Reverse Proxy

Caddy handles automatic HTTPS, reverse proxying, and security headers.

**Configuration** (`Caddyfile`):
```
# API backend
api.arkcore.dev {
    reverse_proxy backend:8080
    
    header {
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        X-XSS-Protection "1; mode=block"
    }
    
    encode gzip
    
    log {
        output stdout
        format json
    }
}

# Frontend
arkcore.dev {
    reverse_proxy frontend:3000
    
    header {
        X-Content-Type-Options nosniff
        X-Frame-Options SAMEORIGIN
        X-XSS-Protection "1; mode=block"
    }
    
    encode gzip
    
    log {
        output stdout
        format json
    }
}

# Redirect www to non-www
www.arkcore.dev {
    redir https://arkcore.dev{uri} permanent
}
```

**Features**:
- Automatic SSL certificate provisioning via Let's Encrypt
- HTTP to HTTPS redirection
- Security headers for XSS and clickjacking protection
- Gzip compression for performance
- JSON logging for observability
- www to non-www redirect

### Production Environment Variables

**Critical Variables** (in `.env` on EC2):
```bash
# Database
POSTGRES_PASSWORD=<strong-random-password>
ARK_DATABASE.HOST=postgres
ARK_DATABASE.NAME=ark
ARK_DATABASE.USER=ark
ARK_DATABASE.PASSWORD=<strong-random-password>
ARK_DATABASE.MAX_OPEN_CONNS=25

# Clerk Authentication (LIVE keys, not test)
ARK_AUTH.CLERK.SECRET_KEY=sk_live_...
ARK_AUTH.CLERK.JWT_ISSUER=https://clerk.arkcore.dev

# Redis
ARK_REDIS.ADDRESS=redis:6379

# Server
ARK_SERVER.PORT=8080
ARK_SERVER.CORS_ALLOWED_ORIGINS=https://arkcore.dev

# Resend (transactional emails)
ARK_INTEGRATION.RESEND_API_KEY=re_...

# OpenAI (AI features - planned)
ARK_OPENAI.API_KEY=sk-...
ARK_OPENAI.MODEL=gpt-4o-mini

# New Relic (observability)
ARK_OBSERVABILITY.NEW_RELIC.LICENSE_KEY=...
ARK_OBSERVABILITY.LOGGING.LEVEL=info
```

**Security Notes**:
- Never commit `.env` to version control
- Use strong random passwords (e.g., `openssl rand -base64 32`)
- Rotate secrets periodically
- Use Clerk LIVE keys (`sk_live_...`, `pk_live_...`) for production

### Cost Optimization

**Monthly Budget**: $15-25 maximum

**Cost Breakdown**:
- EC2 t3.micro: ~$7-10/month (after free tier)
- Cloudflare: Free tier (sufficient for current traffic)
- Clerk: Free tier (up to 10,000 MAU)
- New Relic: Free tier (100GB/month data ingest)
- Resend: Free tier (3,000 emails/month)
- Domain (arkcore.dev): ~$12/year

**Optimization Strategies**:
- Single VM deployment (no separate database server)
- PostgreSQL with pgvector (no dedicated vector DB like Pinecone)
- GitHub Actions for builds (avoid EC2 memory constraints)
- Aggressive caching with Redis
- Cloudflare CDN for static assets
- Free tier services where possible

### Monitoring and Maintenance

**Health Monitoring**:
```bash
# Check service health
curl https://api.arkcore.dev/health

# View Docker container status
docker compose ps

# View resource usage
docker stats

# View logs
docker compose logs -f backend
docker compose logs -f frontend
```

**Maintenance Tasks**:
```bash
# Update Docker images (after CI/CD push)
docker compose pull
docker compose up -d

# Backup database
docker compose exec postgres pg_dump -U ark ark > backup.sql

# Restore database
docker compose exec -T postgres psql -U ark ark < backup.sql

# Clean up old Docker images
docker image prune -a

# View disk usage
df -h
docker system df
```

**Automated Backups** (Future):
- Daily PostgreSQL dumps to S3
- Retention policy: 7 daily, 4 weekly, 12 monthly
- Automated via cron job or AWS Backup

## Deployment Troubleshooting

### 502 Bad Gateway (Cloudflare)

**Symptoms**: Cloudflare shows 502 error, backend is running

**Causes and Fixes**:
1. **Caddy not serving HTTPS**: Check Caddy logs for SSL certificate errors
   ```bash
   docker compose logs caddy
   ```
2. **Cloudflare SSL mode mismatch**: Ensure SSL/TLS mode is "Full (strict)" in Cloudflare dashboard
3. **Backend not responding**: Check backend health
   ```bash
   docker compose exec backend wget -O- http://localhost:8080/health
   ```

### Frontend Container Crashes

**Symptoms**: Frontend container exits immediately after start

**Causes and Fixes**:
1. **Missing build arguments**: Ensure `VITE_API_URL` and `VITE_CLERK_PUBLISHABLE_KEY` are set in GitHub Actions
2. **Build failures**: Check GitHub Actions logs for build errors
3. **Caddy configuration error**: Verify Caddyfile syntax in frontend Dockerfile

### CORS Errors

**Symptoms**: Browser console shows CORS errors when calling API

**Causes and Fixes**:
1. **Backend CORS configuration**: Ensure `ARK_SERVER.CORS_ALLOWED_ORIGINS` includes production domain
   ```bash
   ARK_SERVER.CORS_ALLOWED_ORIGINS=https://arkcore.dev
   ```
2. **Cloudflare interference**: Disable Cloudflare proxy temporarily to test
3. **Preflight request failures**: Check backend logs for OPTIONS requests

### Database Connection Issues

**Symptoms**: Backend logs show "failed to connect to database"

**Causes and Fixes**:
1. **PostgreSQL not ready**: Wait for health check to pass
   ```bash
   docker compose ps postgres
   ```
2. **Wrong credentials**: Verify `ARK_DATABASE.PASSWORD` matches `POSTGRES_PASSWORD`
3. **Network issues**: Ensure backend and postgres are on same Docker network
   ```bash
   docker network inspect ark_ark-network
   ```

### Memory Issues on t3.micro

**Symptoms**: Services crash with OOM (Out of Memory) errors

**Causes and Fixes**:
1. **Too many services**: Consider disabling New Relic in production if not needed
2. **Resource limits**: Add memory limits to docker-compose.yml
   ```yaml
   backend:
     mem_limit: 512m
   ```
3. **Swap space**: Enable swap on EC2 instance
   ```bash
   sudo fallocate -l 2G /swapfile
   sudo chmod 600 /swapfile
   sudo mkswap /swapfile
   sudo swapon /swapfile
   ```
