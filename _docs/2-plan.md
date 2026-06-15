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
- **Crawl handler** (`internal/api/handler/crawl.go`) — triggers crawl via channel, returns 202
- **Status handler** (`internal/api/handler/status.go`) — last crawl timestamp, project count, is_crawling flag
- **CORS middleware** (`internal/api/middleware/cors.go`)
- **Router** (`internal/api/router.go`) — all routes under `/api/v1`
- **GitHub client** (`internal/crawler/github.go`) — Search API client, paginated (10 pages × 100), rate-limit aware
- **Scheduler** (`internal/crawler/scheduler.go`) — `robfig/cron`, configurable interval, first crawl on startup, manual trigger via channel, concurrency guard
- **Models** (`internal/model/project.go`) — Project, ProjectFilter, Status, Filters types
- **Dockerfile** — multi-stage Alpine build, ~15MB binary

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
- **App.tsx** — BrowserRouter with `/` and `/admin` routes
- **index.css** — minimal global styles
- **Dockerfile** + **nginx.conf** — build → nginx-alpine static serve with `/api/` reverse proxy

### M4 🔄 — Polish & Deploy (In Progress)
- GitHub Pages deployment workflow (GitHub Actions)
- Backend deployment guide (Railway / Fly.io)
- Production checklist review
- Load testing

### M5 ⏳ — Future Improvements
- Redis caching layer
- Better wheel animation (highlight winner, confetti)
- Multi-select filters (checkboxes instead of text inputs)
- Pagination for Admin page project list
- Search within projects
- Dark mode

## Files Created

### Backend — 9 files

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

### Frontend — 11 files

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
│   ├── api.ts                                   # HTTP client
│   ├── vite-env.d.ts                            # Env types
│   ├── App.tsx                                  # Router
│   ├── main.tsx                                 # Entry
│   └── index.css                                # Styles
├── index.html
├── vite.config.ts
├── tsconfig.json
├── Dockerfile
└── nginx.conf
```

### Infrastructure — 3 files

```
├── docker-compose.yml                           # Local dev
├── .env.example                                 # Env template
└── .gitignore
```

## Next Steps

### 1. GitHub Actions — Deploy Frontend to GitHub Pages
Create `.github/workflows/deploy-frontend.yml`:
- Trigger: push to `main`
- Build frontend with `npm ci && npm run build`
- Deploy `dist/` to `gh-pages` branch
- Set `VITE_API_BASE_URL` as repository secret (production backend URL)

### 2. GitHub Actions — Deploy Backend
Create `.github/workflows/deploy-backend.yml`:
- Trigger: push to `main`
- Build Docker image
- Push to Docker Hub / GitHub Container Registry
- Deploy to Railway / Fly.io (optional, depends on target)

### 3. Production Readiness
- Add request validation (limit bounds, sanitize inputs)
- Add structured logging (slog or zerolog)
- Add health check endpoint (`GET /api/v1/health`)
- Configure CORS for production domain
- Add rate limiting to crawl endpoint

### 4. Frontend Improvements
- Multi-select dropdowns for language and topic (instead of text inputs)
- Visual winner highlight after spin (glow, scale animation)
- Share result button (copy project URL)
- Responsive layout for mobile

### 5. Redis Caching (Future)
- Cache `GET /api/v1/filters` response (invalidated on crawl)
- Cache random selections (short TTL)
- Rate limiting counter
