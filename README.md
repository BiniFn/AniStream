# AniWays

An anime tracking and streaming platform that helps users discover, track, and manage their anime library with metadata from multiple sources.

> ⚠️ **Under Active Development** - This project is currently in active development. Features and APIs may change frequently.

> ⚠️ **Disclaimer** - This project is for educational purposes only. Please respect the terms of service of all integrated platforms and ensure compliance with applicable laws in your jurisdiction. Deploying this application may violate terms of service of streaming platforms - use at your own risk.

## Features

- **Anime Discovery & Metadata**: Comprehensive anime database with metadata from AniList, MyAnimeList, and Shikimori
- **Personal Library Management**: Track your watching progress with custom statuses (watching, completed, planned, etc.)
- **OAuth Authentication**: Secure login with support for multiple providers (AniList, MAL)
- **Episode Streaming**: Integration with HiAnime for episode streaming and availability
- **HLS Streaming Proxy**: Dedicated proxy server for handling HLS video streams
- **Multi-Source Sync**: Synchronize your anime lists across different platforms
- **RESTful API**: Well-documented API with OpenAPI specifications
- **Modern Web Interface**: Clean, responsive UI built with SvelteKit and TailwindCSS
- **Real-time Updates**: Background workers for metadata refreshing and library synchronization

## Architecture

AniWays follows a clean architecture pattern with separate concerns:

- **API Server** (`cmd/api`): REST API server handling HTTP requests
- **Proxy Server** (`cmd/proxy`): HLS streaming proxy for video content
- **Worker** (`cmd/worker`): Background job processing for data synchronization and scraping
- **Web Frontend** (`web/`): SvelteKit-based user interface

### Tech Stack

**Backend:**

- Go 1.24.1+ with Chi router
- PostgreSQL database with SQLC for type-safe queries
- Redis for caching and session storage
- GraphQL integration with AniList API
- JWT authentication with OAuth2 flows

**Frontend:**

- SvelteKit with TypeScript
- TailwindCSS for styling
- Shadcn-Svelte component library
- Vite for development and building

**Infrastructure:**

- Docker & Docker Compose for containerization
- Database migrations with golang-migrate
- Background job processing with cron scheduler

## Worker CLI

The worker service includes a powerful CLI for manual operations and debugging:

### Available Commands

```bash
# Daemon mode (default)
./worker                   # Run worker in daemon mode
./worker daemon            # Same as above

# Scraping operations
./worker scrape recently-updated       # Scrape top most 40 recently updated anime (hourly task)
./worker scrape all-recently-updated   # Scrape all recently updated pages and upsert into DB
./worker scrape full-seed              # Complete database seed (A-Z + all recently updated) # (one-time use)

# Library management
./worker library retry-failed          # Retry failed library syncs

# Authentication
./worker auth refresh-tokens           # Refresh OAuth tokens
```

## Local Development Setup

### Prerequisites

- Go 1.24.1+
- Node.js 18+ (Bun recommended)
- Docker & Docker Compose

### Step 0: Start Database Services

```bash
# Start PostgreSQL and Redis containers
make dev-docker-up
```

### Step 1: Install Required Binaries

```bash
# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/Khan/genqlient@latest
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/swaggo/swag/cmd/swag@latest
npm install -g swagger2openapi
```

### Step 2: Generate Required Files

```bash
# Generate database code, GraphQL client, and API docs
make sqlc
make genqlient
make openapi
```

### Step 3: Install Frontend Dependencies

```bash
cd web
bun install
cd ..
```

### Step 4: Start Development Servers

**Option A: Individual Terminals**

```bash
# Terminal 1: API Server
make dev-api

# Terminal 2: Background Worker
make dev-worker

# Terminal 3: HLS Streaming Proxy
make dev-proxy

# Terminal 4: Frontend
cd web && bun run dev
```

**Option B: Tmux Session (Easier Setup)**

```bash
make tmux
```

## Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env.local` and update the values:

```bash
cp .env.example .env.local
```

### Required Environment Variables

```bash
# Application
APP_ENV=development
APP_PORT=8080

# URLs
ALLOWED_ORIGINS=http://localhost:3000
FRONTEND_URL=http://localhost:3000
API_URL=http://localhost:8080

# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/aniways?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=password

# OAuth Providers
MYANIMELIST_CLIENT_ID=your_mal_client_id
MYANIMELIST_CLIENT_SECRET=your_mal_client_secret
ANILIST_CLIENT_ID=your_anilist_client_id
ANILIST_CLIENT_SECRET=your_anilist_client_secret

# Cloudinary (for image uploads)
CLOUDINARY_NAME=your_cloudinary_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret

# Email Service
RESEND_API_KEY=your_resend_api_key
RESEND_FROM_EMAIL=noreply@yourdomain.com

# Cookie Domain
COOKIE_DOMAIN=localhost
```

## Available Commands

### Local Development

```bash
make dev-api          # Run API server with hot reload
make dev-proxy        # Run HLS streaming proxy server
make dev-worker       # Run background worker
make tmux             # Start all services in tmux session
```

### Code Generation

```bash
make sqlc             # Generate type-safe database code
make genqlient        # Generate GraphQL client from schema
make openapi          # Generate OpenAPI documentation
make migrate          # Create new migration files (interactive)
```

### Docker Compose (Dev)

```bash
make dev-docker-up    # Start development containers
make dev-docker-down  # Stop development containers
make dev-docker-logs  # View container logs
```

### Help

```bash
make help             # Show all available commands
```

## API Documentation

Once the API server is running, you can access:

- **Swagger UI**: `http://localhost:8080/swagger/`
- **API Specification**: `/docs/openapi.yaml`

The frontend runs on `http://localhost:3000` by default.

## Deployment

### Dokploy Deployment (Recommended)

[Dokploy](https://dokploy.com/) provides a simple web interface for Docker deployments with automatic SSL and domain management.

#### Prerequisites

- VPS with Docker installed
- Domain names pointing to your VPS
- Dokploy installed on your VPS
- Docker Hub account for CI/CD

#### CI/CD Workflow

The project uses GitHub Actions for automated deployment:

1. **Push to main branch** → Triggers GitHub Actions
2. **GitHub Actions** → Builds Docker images and pushes to Docker Hub
3. **Docker Hub webhook** → Notifies Dokploy of new images
4. **Dokploy** → Pulls latest images and redeploys services

#### Setup Steps

1. **Install Dokploy** on your VPS following their [installation guide](https://dokploy.com/docs/installation/)

2. **Configure GitHub Secrets**:
   - `DOCKER_USERNAME`: Your Docker Hub username
   - `DOCKER_PASSWORD`: Your Docker Hub password/token

3. **Create databases in Dokploy**:
   - **PostgreSQL**: Create a PostgreSQL database service
   - **Redis**: Create a Redis service

4. **Create application services in Dokploy**:
   - **API Service**: Docker image `yourusername/aniways-api:latest`, port 8080
   - **Proxy Service**: Docker image `yourusername/aniways-proxy:latest`, port 1234
   - **Worker Service**: Docker image `yourusername/aniways-worker:latest` (background service)

5. **Set environment variables** in Dokploy interface:

```bash
# Application
APP_ENV=production
APP_PORT=8080

# URLs
ALLOWED_ORIGINS=https://yourdomain.com
FRONTEND_URL=https://yourdomain.com
API_URL=https://api.yourdomain.com

# Database
DATABASE_URL=postgres://postgres:your_password@postgres:5432/aniways

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=your_redis_password

# OAuth Providers
MYANIMELIST_CLIENT_ID=your_mal_client_id
MYANIMELIST_CLIENT_SECRET=your_mal_client_secret
ANILIST_CLIENT_ID=your_anilist_client_id
ANILIST_CLIENT_SECRET=your_anilist_client_secret

# Cloudinary
CLOUDINARY_NAME=your_cloudinary_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret

# Email
RESEND_API_KEY=your_resend_api_key
RESEND_FROM_EMAIL=noreply@yourdomain.com

# Cookie Domain
COOKIE_DOMAIN=yourdomain.com
```

6. **Configure Docker Hub webhooks** to trigger Dokploy redeployment on new image pushes

#### Benefits of This Setup

- **Automated CI/CD**: Push to main → automatic deployment
- **Docker Hub Integration**: Centralized image registry
- **Webhook Triggers**: Dokploy automatically pulls latest images
- **Automatic SSL**: Let's Encrypt certificates managed by Dokploy
- **Service Monitoring**: Built-in health checks and monitoring
- **Easy Rollbacks**: Previous image tags available for quick rollbacks

### Docker Compose Deployment (Alternative)

For simpler deployments without CI/CD, you can use Docker Compose directly.

#### Prerequisites

- VPS with Docker and Docker Compose installed
- Domain names pointing to your VPS
- Reverse proxy (Nginx/Traefik) for SSL termination

#### Setup Steps

1. **Create environment file**:

```bash
# Create production environment file
cat > .env.prod << EOF
# Docker Hub
DOCKER_USERNAME=yourusername

# Database passwords
POSTGRES_PASSWORD=your_secure_postgres_password
REDIS_PASSWORD=your_secure_redis_password

# OAuth Providers
MYANIMELIST_CLIENT_ID=your_mal_client_id
MYANIMELIST_CLIENT_SECRET=your_mal_client_secret
ANILIST_CLIENT_ID=your_anilist_client_id
ANILIST_CLIENT_SECRET=your_anilist_client_secret

# Cloudinary
CLOUDINARY_NAME=your_cloudinary_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret

# Email
RESEND_API_KEY=your_resend_api_key
RESEND_FROM_EMAIL=noreply@yourdomain.com

# Application URLs
COOKIE_DOMAIN=yourdomain.com
ALLOWED_ORIGINS=https://yourdomain.com
FRONTEND_URL=https://yourdomain.com
API_URL=https://api.yourdomain.com
EOF
```

2. **Deploy with Docker Compose**:

```bash
# Pull latest images
docker compose -f docker/docker-compose.prod.yaml --env-file .env.prod pull

# Start services
docker compose -f docker/docker-compose.prod.yaml --env-file .env.prod up -d
```

3. **Set up reverse proxy** (Nginx example):

```nginx
# /etc/nginx/sites-available/aniways
server {
    listen 80;
    server_name yourdomain.com api.yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name yourdomain.com;

    ssl_certificate /path/to/ssl/cert.pem;
    ssl_certificate_key /path/to/ssl/key.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

server {
    listen 443 ssl;
    server_name api.yourdomain.com;

    ssl_certificate /path/to/ssl/cert.pem;
    ssl_certificate_key /path/to/ssl/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

#### Benefits of Docker Compose

- **Simple setup**: Single command deployment
- **Full control**: Direct access to all services
- **Easy debugging**: Direct access to logs and containers
- **Cost effective**: No additional platform fees
- **Customizable**: Easy to modify configuration

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature-name`)
3. Commit your changes (`git commit -m 'Add your feature'`)
4. Push to the branch (`git push origin feature/your-feature-name`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and run `go fmt`
- Write tests for new features
- Update documentation when needed
- Use conventional commit messages

---

**Made with ❤️ by [Coeeter](https://github.com/coeeter)**
