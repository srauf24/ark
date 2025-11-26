# Build stage
FROM oven/bun:1 AS builder

WORKDIR /app

# Copy root package files for workspace setup
COPY ../../package.json ../../bun.lockb ./

# Copy workspace packages
COPY ../../packages ./packages

# Copy frontend package files
COPY package.json ./apps/frontend/

# Install all workspace dependencies
RUN bun install --frozen-lockfile

# Copy frontend source
COPY . ./apps/frontend/

# Build the frontend application
WORKDIR /app/apps/frontend
RUN bun run build

# Production stage
FROM nginx:alpine

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder
COPY --from=builder /app/apps/frontend/dist /usr/share/nginx/html

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
