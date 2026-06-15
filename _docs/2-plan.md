# Project Plan

## Milestones

### M1 ✅ — Foundation (Done)
- Initialized project structure (backend, frontend, docker-compose)
- Go module with Fiber, pgx, robfig/cron, golang.org/x/exp dependencies
- Supabase migration (`001_create_projects.sql`) — `projects` table with indexes
- Backend entry point (`cmd/server/main.go`) with graceful shutdown
- React + Vite + TypeScript scaffold with routing (`/` and `/admin`)
- Docker Compose with backend, frontend, redis services
- `.env.example` with all configuration variables
- `.gitignore`, `vite.config.ts` with configurable base path

### M2 ✅ — Backend API & Crawler (Done)
- **Config** (`config/config.go`) — env loader, auto-builds PostgreSQL connection string from `SUPABASE_URL` + `SUPABASE_KEY`
- **DB layer** (`internal/db/supabase.go`) — pgx connection pool, embedded SQL migration runner
- **Project handler** (`internal/api/handler/projects.go`):
  - `List` — paginated with filters: language, topic, min_stars, page, limit
  - `Random` — random selection with same filters, `ORDER BY RANDOM() LIMIT X`
  - `GetFilters` — distinct languages and topics from DB
- **Crawl handler** (`internal/api/handler/crawl.go`) — triggers crawl via channel, returns 202 or 429
- **Status handler** (`internal/api/handler/status.go`) — last crawl timestamp, project count, is_crawling flag
- **CORS middleware** (`internal/api/middleware/cors.go`)
- **Router** (`internal/api/router.go`) — all routes under `/api/v1`
- **GitHub client** (`internal/crawler/github.go`) — Search API client, paginated (10 pages × 100), rate-limit aware
- **Scheduler** (`internal/crawler/scheduler.go`) — `robfig/cron`, configurable interval, first crawl on startup, manual trigger via channel, concurrency guard
- **Models** (`internal/model/project.go`) — Project, ProjectFilter, Status, Filters types
- **Dockerfile** — multi-stage Alpine build, ~15MB binary
- **entrypoint.sh** — DNS workaround for Windows/VPN environments

### M3 ✅ — Frontend Core (Done)
- **API client** (`api.ts`) — typed fetch wrapper for all endpoints
- **Wheel** (`components/Wheel.tsx`) — Canvas-based:
  - Equal colored slices with project names
  - Pointer triangle indicator
  - Random spin duration + ease-out cubic animation
  - Click-to-spin or Spin button
  - Winner detection and callback
- **FilterPanel** (`components/FilterPanel.tsx`) — limit slider, text inputs for language/topic, number input for min stars
- **ProjectCard** (`components/ProjectCard.tsx`) — avatar, GitHub link, stars, language badge, topic tags
- **Home page** (`pages/Home.tsx`) — filters + spin + wheel + winner card
- **Admin page** (`pages/Admin.tsx`) — crawl trigger button, status display, available filters list
- **App.tsx** — BrowserRouter with `basename={import.meta.env.BASE_URL}`
- **vite-env.d.ts** — typed env variables (`VITE_API_BASE_URL`, `VITE_BASE_URL`)
- **index.css** — minimal global styles
- **Dockerfile** + **nginx.conf** — build → nginx-alpine static serve with `/api/` reverse proxy

### M4 ✅ — Deploy & CI/CD (Done)
- GitHub Actions: frontend deploy to GitHub Pages
  - Builds with `npm ci && npm run build`, passes `VITE_API_BASE_URL` and `VITE_BASE_URL`
  - Deploys `dist/` to GitHub Pages via official actions
- GitHub Actions: backend Docker image build & push
  - Builds multi-stage Docker image via BuildKit with layer caching
  - Pushes to `ghcr.io/<owner>/github-wheel-of-random-backend` (tags: `latest` + sha)
  - Triggers Render.com deploy via Deploy Hook (optional, requires `RENDER_DEPLOY_HOOK_URL` secret)
- Backend hosted on Render.com via Docker (pulls image from GHCR)
- Frontend connected to production backend via `VITE_API_BASE_URL` secret

### M5 ⏳ — Future Improvements
- Redis caching layer
- Better wheel animation (highlight winner, confetti)
- Multi-select filters (checkboxes instead of text inputs)
- Pagination for Admin page project list
- Search within projects
- Dark mode
- Request validation and health check endpoint
- Rate limiting on crawl endpoint
- Structured logging (slog or zerolog)

## Files Created

### Backend — 10 files

```
backend/
├── cmd/server/main.go                          # Entry point
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   │   ├── projects.go                     # List, Random, GetFilters
│   │   │   ├── crawl.go                        # TriggerCrawl
│   │   │   └── status.go                       # Status
│   │   ├── middleware/
│   │   │   └── cors.go                         # CORS
│   │   └── router.go                           # Routes
│   ├── crawler/
│   │   ├── github.go                           # GitHub API client
│   │   └── scheduler.go                        # Cron + upsert
│   ├── db/
│   │   ├── supabase.go                         # pgx pool + migrations
│   │   └── migrations/
│   │       └── 001_create_projects.sql          # Schema
│   └── model/
│       └── project.go                           # Types
├── config/
│   └── config.go                                # Env loader
├── go.mod / go.sum
├── Dockerfile
└── entrypoint.sh
```

### Frontend — 12 files

```
frontend/
├── src/
│   ├── components/
│   │   ├── Wheel.tsx                            # Canvas wheel
│   │   ├── FilterPanel.tsx                      # Filter controls
│   │   └── ProjectCard.tsx                      # Result card
│   ├── pages/
│   │   ├── Home.tsx                             # Main page
│   │   └── Admin.tsx                            # Admin panel
│   ├── hooks/                                   # (empty, future use)
│   ├── api.ts                                   # HTTP client
│   ├── vite-env.d.ts                            # Env types
│   ├── App.tsx                                  # Router with basename
│   ├── main.tsx                                 # Entry
│   └── index.css                                # Styles
├── index.html
├── vite.config.ts
├── tsconfig.json
├── Dockerfile
└── nginx.conf
```

### Infrastructure — 5 files

```
├── docker-compose.yml                           # Local dev
├── .env.example                                 # Env template
├── .gitignore
└── .github/workflows/
    ├── deploy-frontend.yml                      # GitHub Pages deploy
    └── deploy-backend.yml                       # Docker build + Render deploy
```

## Deployment

### Frontend → GitHub Pages
- **Trigger:** Push to `main` (changes in `frontend/` or workflow file)
- **Build:** `npm ci && npm run build` with `VITE_API_BASE_URL` and `VITE_BASE_URL` env
- **Deploy:** GitHub Pages via `actions/deploy-pages`
- **URL:** `https://<user>.github.io/github-wheel-of-random/`
- **Secret needed:** `VITE_API_BASE_URL` (production backend URL)

### Backend → Render.com
- **Trigger:** Push to `main` (changes in `backend/` or workflow file)
- **Build:** Multi-stage Docker build, push to `ghcr.io`
- **Deploy:** Render Deploy Hook (optional, requires `RENDER_DEPLOY_HOOK_URL` secret)
- **Setup on Render:**
  1. Create Web Service → Deploy via Docker
  2. Image: `ghcr.io/<user>/github-wheel-of-random-backend:latest`
  3. Add env vars: `SUPABASE_URL`, `SUPABASE_KEY`, `GITHUB_TOKEN`, `CRAWL_INTERVAL`
  4. Copy Deploy Hook URL → GitHub secret `RENDER_DEPLOY_HOOK_URL`
