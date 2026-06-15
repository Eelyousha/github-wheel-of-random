# Architecture

## System Overview

```
┌─────────────────────┐     ┌──────────────────────┐     ┌────────────────────┐
│   Frontend (React)  │────▶│   Backend (Go/Fiber) │────▶│  Supabase (PG)    │
│   GitHub Pages      │     │   Railway / VPS      │     │                   │
│   Docker container  │     │   Docker container   │     │  Managed service  │
│   Static SPA        │     │   REST API + Crawler │     │                   │
└─────────────────────┘     └──────────┬───────────┘     └────────────────────┘
                                        │
                                ┌───────▼────────┐
                                │  Redis (future) │
                                │  Docker         │
                                │  Cache layer    │
                                └────────────────┘
```

Three independently containerized components communicating over HTTP. Each can be run locally via Docker Compose or deployed separately.

## Stack

| Component   | Technology                   | Hosting                          |
|-------------|------------------------------|----------------------------------|
| **Frontend**| React + Vite (TypeScript)   | GitHub Pages                     |
| **Backend** | Go (Fiber)                   | Railway / Fly.io / VPS           |
| **Database**| Supabase (PostgreSQL)        | Managed (Supabase Cloud)         |
| **Cache**   | Redis (future)               | Upstash / Docker                 |
| **Crawler** | Go (built-in, `robfig/cron`) | Embedded in Backend process      |

## Project Structure

```
github-wheel-of-random/
├── _docs/                         # Architecture & planning docs
├── backend/
│   ├── cmd/server/main.go         # Entry point, graceful shutdown
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/
│   │   │   │   ├── projects.go    # List, Random, GetFilters handlers
│   │   │   │   ├── crawl.go       # Manual crawl trigger
│   │   │   │   └── status.go      # Crawl status & project count
│   │   │   ├── middleware/
│   │   │   │   └── cors.go        # CORS middleware
│   │   │   └── router.go          # Fiber router registration
│   │   ├── crawler/
│   │   │   ├── github.go          # GitHub Search API client (paginated)
│   │   │   └── scheduler.go       # Cron scheduler + upsert logic
│   │   ├── db/
│   │   │   ├── supabase.go        # pgx pool, migrations runner
│   │   │   └── migrations/
│   │   │       └── 001_create_projects.sql
│   │   └── model/
│   │       └── project.go         # Project, ProjectFilter, Status types
│   ├── config/
│   │   └── config.go              # Env loader
│   ├── go.mod / go.sum
│   ├── Dockerfile                 # Multi-stage Alpine
│   └── entrypoint.sh              # Startup script
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Wheel.tsx          # Canvas-based spinning wheel
│   │   │   ├── FilterPanel.tsx    # Language, topics, stars, limit controls
│   │   │   └── ProjectCard.tsx    # GitHub-style result card
│   │   ├── pages/
│   │   │   ├── Home.tsx           # Wheel page with filters + result
│   │   │   └── Admin.tsx          # Crawl trigger, status, filter list
│   │   ├── api.ts                 # Typed HTTP client
│   │   ├── vite-env.d.ts          # Env type declarations
│   │   ├── App.tsx                # Router setup
│   │   ├── main.tsx               # Entry point
│   │   └── index.css              # Minimal styles
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── Dockerfile                 # Build → nginx-alpine
│   └── nginx.conf                 # Proxy /api/ → backend
├── docker-compose.yml             # backend + frontend + redis
├── .env.example
├── .gitignore
└── README.md
```

## Data Flow

1. **Crawler** (scheduled via `robfig/cron`, interval from `CRAWL_INTERVAL`) fetches popular repos from GitHub Search API (`stars:>1000`, top 10 pages × 100 per page)
2. Crawler upserts projects into Supabase (`projects` table) — `INSERT ... ON CONFLICT (full_name) DO UPDATE`
3. **Frontend** (React SPA) starts immediately — shows filter controls and Spin button
4. User sets filters (language, topic, min stars) and limit X via FilterPanel
5. User clicks "Spin" → calls `GET /api/v1/projects/random?limit=X&lang=...&topic=...`
6. **Backend** builds a dynamic SQL query: filter → `ORDER BY RANDOM() LIMIT X`
7. Returns X random projects as JSON array
8. **Frontend** draws a Canvas wheel with equal slices, animates a spin (random duration, ease-out cubic)
9. When animation stops, the slice under the pointer is the winner
10. **ProjectCard** renders below the wheel with avatar, stars, description, topics, and GitHub link

## API Endpoints

| Method | Path                              | Description                                      |
|--------|-----------------------------------|--------------------------------------------------|
| GET    | `/api/v1/projects`                | Paginated list (query: page, limit, lang, topic, min_stars) |
| GET    | `/api/v1/projects/random?limit=X` | X random projects (max 100)                      |
| GET    | `/api/v1/filters`                 | Available languages and topics for filter UI     |
| POST   | `/api/v1/crawl`                   | Trigger a manual crawl (returns 202 Accepted)    |
| GET    | `/api/v1/status`                  | Crawl stats: last_crawl, project_count, is_crawling |

## Database Schema

```sql
CREATE TABLE projects (
    id          BIGSERIAL PRIMARY KEY,
    full_name   TEXT UNIQUE NOT NULL,        -- e.g. "golang/go"
    name        TEXT NOT NULL,
    owner       TEXT NOT NULL,
    description TEXT,
    stars       INTEGER NOT NULL DEFAULT 0,
    language    TEXT,
    topics      TEXT[],                       -- PostgreSQL array
    html_url    TEXT NOT NULL,
    avatar_url  TEXT,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_language ON projects (language);
CREATE INDEX idx_projects_stars ON projects (stars DESC);
CREATE INDEX idx_projects_topics ON projects USING GIN (topics);
```

## Environment Variables

| Variable         | Required | Description                           |
|------------------|----------|---------------------------------------|
| `SUPABASE_URL`   | ✅       | Supabase project URL (with or without scheme) |
| `SUPABASE_KEY`   | ✅       | Supabase service role key             |
| `GITHUB_TOKEN`   | ❌       | GitHub PAT (higher API rate limit)    |
| `CRAWL_INTERVAL` | ❌       | Duration string, default `6h`         |
| `BACKEND_PORT`   | ❌       | API port, default `8000`              |

## Frontend Components

### Wheel (`Wheel.tsx`)
- Canvas-based spinning wheel with colored equal slices
- Each slice displays the project name (truncated to 12 chars)
- Pointer triangle at the top indicates the winner
- Animation: random number of full spins (5–10) + random slice offset
- Easing: ease-out cubic over ~4 seconds
- Click-to-spin or use the Spin button
- Accepts `projects: Project[]` and `onFinish: (winner: Project) => void` props

### FilterPanel (`FilterPanel.tsx`)
- Range slider for limit X (2–100)
- Text inputs for language and topic filter
- Number input for minimum stars
- Passes values up to Home page via props

### ProjectCard (`ProjectCard.tsx`)
- Displays winner: owner avatar, full_name as GitHub link, description, star count, language badge, up to 5 topic tags

### Pages
- **Home** — combines FilterPanel, Spin button, Wheel, and winner ProjectCard
- **Admin** — shows crawler status, manual crawl button, list of available languages and topics

## Crawler

- Uses GitHub Search API: `GET /search/repositories?q=stars:>1000&sort=stars&per_page=100`
- Paginates up to 10 pages (configurable via code constant)
- Respects rate limits; uses `GITHUB_TOKEN` if provided (5000 req/h) else anonymous (60 req/h)
- Upserts via `INSERT ... ON CONFLICT (full_name) DO UPDATE SET ...`
- First crawl runs immediately on startup, then repeats every `CRAWL_INTERVAL`
- Manual trigger via `POST /api/v1/crawl` (non-blocking, returns 202)
- Concurrency guard: skips if already crawling
- Scheduler uses `robfig/cron/v3` with `@every <interval>` syntax
