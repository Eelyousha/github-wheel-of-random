# GitHub Wheel of Random 🎡

> Spin a wheel of popular open-source GitHub projects. Filter by language, topic, or stars — or just leave it all in and discover something new.

## Screenshots

*(coming soon)*

## Quick Start

```bash
# 1. Clone
git clone https://github.com/<your-org>/github-wheel-of-random.git
cd github-wheel-of-random

# 2. Configure
cp .env.example .env
# Edit .env with your Supabase URL, key, and GitHub token

# 3. Start everything
docker compose up --build
```

- **Frontend:** http://localhost:3000
- **Backend:** http://localhost:8000
- **API health:** http://localhost:8000/api/v1/status

## Architecture (3-component)

```
Frontend (React SPA)  ──HTTP──▶  Backend (Go/Fiber)  ──PG──▶  Supabase
   GitHub Pages                      Railway / VPS                Cloud
```

| Component | Stack | Hosting |
|-----------|-------|---------|
| Frontend  | React + Vite, TypeScript | GitHub Pages |
| Backend   | Go (Fiber), robfig/cron   | Railway / Fly.io / VPS |
| Database  | PostgreSQL (Supabase)      | Supabase Cloud |

All three are containerized (Docker). A Redis container is also included for future caching.

## Project Structure

```
github-wheel-of-random/
├── _docs/                     # Architecture & planning docs
├── frontend/                  # React SPA
│   ├── src/
│   │   ├── components/        # Wheel, FilterPanel, ProjectCard, AdminPanel
│   │   ├── pages/             # Home, Admin
│   │   └── api.ts             # Backend HTTP client
│   ├── Dockerfile
│   └── vite.config.ts
├── backend/                   # Go API
│   ├── cmd/server/main.go     # Entry point
│   ├── internal/
│   │   ├── api/               # HTTP handlers & router
│   │   ├── crawler/           # GitHub crawler + scheduler
│   │   ├── db/                # Supabase client + migrations
│   │   └── model/             # Data types
│   ├── config/config.go
│   └── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

## How It Works

1. **Crawler** (scheduled via cron) fetches popular repos from the GitHub Search API and stores them in Supabase
2. **User** opens the web app and optionally sets filters (language, topics, min stars, count X)
3. **Spin** — frontend requests `GET /api/v1/projects/random?limit=X&lang=...&topic=...`
4. **Backend** returns X random matching projects
5. **Wheel** — frontend draws a wheel with equal slices and animates the spin
6. **Winner** — when the wheel stops, the highlighted project is shown as a card

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/projects` | Paginated list with filters |
| GET | `/api/v1/projects/random?limit=X` | X random projects (max 100) |
| GET | `/api/v1/filters` | Available languages & topics |
| POST | `/api/v1/crawl` | Trigger manual crawl |
| GET | `/api/v1/status` | Crawl stats |

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SUPABASE_URL` | ✅ | Supabase project URL |
| `SUPABASE_KEY` | ✅ | Supabase service role key |
| `GITHUB_TOKEN` | ❌ | Personal access token (higher rate limit) |
| `CRAWL_INTERVAL` | ❌ | Cron or duration (default `6h`) |
| `BACKEND_PORT` | ❌ | API port (default `8000`) |

## Deployment

### Frontend → GitHub Pages

A GitHub Actions workflow builds the Vite app and deploys to the `gh-pages` branch.

1. Set `VITE_API_BASE_URL` as a repository secret (your backend URL)
2. Push to `main` → automatic deploy

### Backend → Railway

1. Connect your GitHub repo to Railway
2. Set env vars in Railway dashboard
3. Railway auto-deploys on push to `main`

### Local

```bash
docker compose up --build
```

## Related Docs

- [Architecture](_docs/1-architecture.md)
- [Project Plan](_docs/2-plan.md)
- [Components](_docs/3-components.md)
- [Workflow](_docs/4-workflow.md)

## License

MIT
