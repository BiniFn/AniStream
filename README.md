# AniWays 🍿

An anime tracking and streaming platform that helps users discover, track, and manage their anime library with metadata from multiple sources.

> ⚠️ **Under Active Development** - This project is currently in active development. Features and APIs may change frequently.

> ⚠️ **Disclaimer** - This project is for educational purposes only. Please respect the terms of service of all integrated platforms and ensure compliance with applicable laws in your jurisdiction. Deploying this application may violate terms of service of streaming platforms - use at your own risk.

## 🚀 Features

- **Anime Discovery & Metadata**: Comprehensive anime database with metadata from AniList, MyAnimeList, and Shikimori
- **Personal Library Management**: Track your watching progress with custom statuses (watching, completed, planned, etc.)
- **OAuth Authentication**: Secure login with support for multiple providers (AniList, MAL)
- **Episode Streaming**: Integration with HiAnime for episode streaming and availability
- **HLS Streaming Proxy**: Dedicated proxy server for handling HLS video streams
- **Multi-Source Sync**: Synchronize your anime lists across different platforms
- **RESTful API**: Well-documented API with OpenAPI specifications
- **Modern Web Interface**: Clean, responsive UI built with SvelteKit and TailwindCSS
- **Real-time Updates**: Background workers for metadata refreshing and library synchronization

## 🏗️ Architecture

AniWays follows a clean architecture pattern with separate concerns:

- **API Server** (`cmd/api`): REST API server handling HTTP requests
- **Proxy Server** (`cmd/proxy`): HLS streaming proxy for video content
- **Worker** (`cmd/worker`): Background job processing for data synchronization
- **Web Frontend** (`web/`): SvelteKit-based user interface

### Tech Stack

**Backend:**

- Go 1.24+ with Chi router
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

## 🛠️ Local Development Setup

### Prerequisites

- Go 1.24+
- Node.js 18+ (with Bun recommended)
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
bun install  # or npm install
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

## 🔧 Configuration

The application uses environment variables for configuration. Create a `.env.local` file with the following variables:

### Required Environment Variables

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/aniways

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

# Application
COOKIE_DOMAIN=localhost
```

### Optional Environment Variables

```bash
APP_ENV=development
APP_PORT=8080
ALLOWED_ORIGINS=*
FRONTEND_URL=http://localhost:3000
API_URL=http://localhost:8080
```

## 📋 Available Commands

### Development

```bash
make dev-api          # Run API server with hot reload
make dev-worker       # Run background worker
make dev-proxy        # Run HLS streaming proxy server
make tmux             # Start all services in tmux session
```

### Code Generation

```bash
make sqlc             # Generate type-safe database code
make genqlient        # Generate GraphQL client from schema
make openapi          # Generate OpenAPI documentation
```

### Database

```bash
make migrate <name>   # Create new migration files
```

### Docker

```bash
make dev-docker-up    # Start development containers
make dev-docker-down  # Stop development containers
make dev-docker-logs  # View container logs
```

### Building

```bash
make build            # Build API and proxy binaries
make docker-build-api # Build API Docker image
```

## 🌐 API Documentation

Once the API server is running, you can access:

- **Swagger UI**: `http://localhost:8080/swagger/`
- **API Specification**: `/docs/openapi.yaml`

The frontend runs on `http://localhost:3000` by default.

## 🚀 Deployment

### Docker Swarm Deployment

The project includes production-ready Docker Swarm configuration with Traefik reverse proxy and Let's Encrypt SSL certificates.

#### Prerequisites

- Docker Swarm initialized on your VPS
- Domain names pointing to your VPS

#### Environment Setup

Create the following secrets in your swarm:

```bash
# Create secrets
echo "your_postgres_password" | docker secret create postgres_password -
echo "your_redis_password" | docker secret create redis_password -

# Create environment file secret
cat > aniways.env << EOF
DATABASE_URL=postgres://postgres:your_postgres_password@db:5432/aniways
REDIS_ADDR=redis:6379
REDIS_PASSWORD=your_redis_password
MYANIMELIST_CLIENT_ID=your_mal_client_id
MYANIMELIST_CLIENT_SECRET=your_mal_client_secret
ANILIST_CLIENT_ID=your_anilist_client_id
ANILIST_CLIENT_SECRET=your_anilist_client_secret
CLOUDINARY_NAME=your_cloudinary_name
CLOUDINARY_API_KEY=your_cloudinary_api_key
CLOUDINARY_API_SECRET=your_cloudinary_api_secret
RESEND_API_KEY=your_resend_api_key
RESEND_FROM_EMAIL=noreply@yourdomain.com
COOKIE_DOMAIN=yourdomain.com
API_URL=https://api.yourdomain.com
FRONTEND_URL=https://yourdomain.com
EOF

docker secret create aniways_env aniways.env
rm aniways.env
```

#### Deploy to Swarm

```bash
# Set environment variables
export POSTGRES_USER=postgres
export POSTGRES_DB=aniways
export LETSENCRYPT_EMAIL=your-email@example.com
export API_DOMAIN=api.yourdomain.com
export PROXY_DOMAIN=yourdomain.com

# Deploy the stack
docker stack deploy -c docker/docker-stack.yaml aniways
```

### Alternative: Dokploy

For easier deployment, you can use [Dokploy](https://dokploy.com/) which provides a simple web interface for Docker deployments:

1. Install Dokploy on your VPS
2. Create a new project in Dokploy
3. Connect your repository
4. Configure environment variables through the Dokploy interface
5. Deploy with automatic SSL and domain management

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and run `go fmt`
- Write tests for new features
- Update documentation when needed
- Use conventional commit messages

---

**Made with ❤️ by [Coeeter](https://github.com/coeeter)**

