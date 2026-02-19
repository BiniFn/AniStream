# AniWays

An anime tracking and streaming platform that helps users discover, track, and manage their anime library with metadata from multiple sources.

> ⚠️ **Under Active Development**  
> Features, APIs, and architecture may change frequently.

> ⚠️ **Disclaimer**  
> This project is for educational purposes only. Use responsibly and comply with platform terms of service.

---

## ✨ Features

- **Anime Discovery & Metadata** – Integrated data from AniList, MyAnimeList, and Shikimori
- **Personal Library Management** – Track anime with custom statuses (watching, completed, planned, etc.)
- **OAuth Authentication** – Secure login via AniList & MAL
- **Episode Streaming** – Integration with HiAnime for episode playback
- **HLS Proxy Server** – Handles HLS streaming and bypasses CORS restrictions
- **Multi-Source Sync** – Keep anime lists synced across platforms
- **RESTful API** – OpenAPI-documented backend for integrations
- **Modern Web App** – Built with SvelteKit + TailwindCSS
- **Background Jobs** – Workers handle scraping, syncing, and metadata updates

---

## 🧱 Architecture

AniWays is built around four core services:

| Service          | Path         | Description                                             |
| ---------------- | ------------ | ------------------------------------------------------- |
| **API Server**   | `cmd/api`    | REST API handling authentication, library, and metadata |
| **Proxy Server** | `cmd/proxy`  | Dedicated HLS streaming proxy                           |
| **Worker**       | `cmd/worker` | Background processor for scraping and sync jobs         |
| **Web Frontend** | `web/`       | SvelteKit UI for users                                  |

### Tech Stack

**Backend**

- Go 1.24+ with Chi router
- PostgreSQL + SQLC (type-safe queries)
- Redis (cache, sessions)
- GraphQL (AniList)
- JWT + OAuth2 authentication

**Frontend**

- SvelteKit + TypeScript
- TailwindCSS + Shadcn-Svelte
- Vite build system

**Infra**

- Docker-based development
- Database migrations with `golang-migrate`
- Background tasks with cron jobs

---

## 🧰 Worker CLI

The worker service includes CLI tools for scraping and maintenance:

```bash
# Run worker in daemon mode
./worker

# Scraping operations
./worker scrape recently-updated      # Fetch 40 most recent anime
./worker scrape all-recently-updated  # Fetch all recent anime
./worker scrape full-seed             # Complete database seed (A–Z)

# Library sync operations
./worker library retry-failed         # Retry failed syncs

# Auth management
./worker auth refresh-tokens          # Refresh OAuth tokens
```

---

## ⚙️ Local Development

### Prerequisites

```bash
make setup
```

This command installs dependencies, generates code, starts containers, and configures your local environment automatically.

### Running services

```bash
# Start API server
make dev-api
# Start Proxy server
make dev-proxy
# Start Worker
make dev-worker
# Start Frontend
cd web && bun run dev
```

---

## 📚 API Documentation

Once the API is running:

- **Swagger UI** → [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- **Frontend** → [http://localhost:3000](http://localhost:3000)

---


## 🚀 Deploying to Render

This repo now includes a `render.yaml` Blueprint for the SvelteKit website (`web/`) using the Node adapter.

1. Create a new **Blueprint** service in Render and point it to this repository.
2. Set `PUBLIC_API_URL` in Render to your deployed API URL.
3. Deploy — Render will run `npm install && npm run build` inside `web/` and start the app with `npm run start`.

## 🧑‍💻 Contributing

1. Fork the repo
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit & push (`git commit -m "feat: add feature"`)
4. Open a Pull Request

---

**Made with ❤️ by [Coeeter](https://github.com/coeeter)**
