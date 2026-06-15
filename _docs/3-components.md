# Components

## System Components

### 1. Frontend — React SPA

**Location:** `frontend/`

**Purpose:** Single-page application hosted on GitHub Pages. Users spin the wheel, configure filters, and view results.

**Runtime:** Static files served via GitHub Pages (production) or nginx (locally via Docker Compose).

**Key dependencies:**
- `react` + `react-dom` — UI library
- `react-router-dom` — client-side routing (`/`, `/admin`)

**Pages:**
| Route   | Component | Description                        |
|---------|-----------|------------------------------------|
| `/`     | `Home`    | Wheel + filters + result           |
| `/admin`| `Admin`   | Crawl trigger, status, filter list |

**Components:**
| Component       | Responsibility                                          |
|-----------------|---------------------------------------------------------|
| `Wheel`         | Canvas element that draws equal slices, animates spin (ease-out cubic, ~4s), highlights winner |
| `FilterPanel`   | Controls: limit slider (2–100), language input, topic input, min-stars input |
| `ProjectCard`   | Displays winner: avatar, full_name, stars, description, language, topics, GitHub link |
| (Admin page)    | Shows last crawl timestamp, project count, is_crawling flag, manual crawl button, available languages/topics |

**Flow:**
1. Mount → user sees filter controls and Spin button
2. User optionally configures filters (persisted in component state)
3. User clicks "Spin"
4. `GET /api/v1/projects/random?limit=X&lang=...&topic=...` called
5. Response array → Wheel draws N equal colored slices with project names
6. Animation plays (random spin count 5–10, ease-out cubic)
7. Winner determined by pointer position → `onFinish` callback
8. ProjectCard rendered below the wheel

**API client** (`api.ts`):
- Typed functions: `getProjects`, `getRandomProjects`, `getFilters`, `getStatus`, `triggerCrawl`
- Base URL from `VITE_API_BASE_URL` env var (default `http://localhost:8000`)

**Router config** (`App.tsx`):
- Uses `BrowserRouter` with `basename={import.meta.env.BASE_URL}`
- `import.meta.env.BASE_URL` comes from Vite's `base` config option
- On GitHub Pages: `/github-wheel-of-random/`, locally: `/`

**Vite config** (`vite.config.ts`):
- `base` defaults to `/github-wheel-of-random/` (can be overridden via `VITE_BASE_URL` env)

---

### 2. Backend — Go / Fiber API

**Location:** `backend/`

**Purpose:** REST API that serves project data, runs the GitHub crawler on a schedule, and provides filter metadata.

**Runtime:** Compiled Go binary, listens on `:8000`.

**Key dependencies:**
- `github.com/gofiber/fiber/v2` — HTTP framework
- `github.com/jackc/pgx/v5` — PostgreSQL driver (Supabase)
- `github.com/robfig/cron/v3` — Scheduler
- `golang.org/x/exp` — rand.Shuffle

**API Handlers:**

| Handler | Endpoint | Logic |
|---|---|---|
| `List` | `GET /api/v1/projects` | Dynamic SQL with filters: language, topic (ANY), min_stars. Paginated with page/limit. Ordered by stars DESC. |
| `Random` | `GET /api/v1/projects/random` | Same filters as List but `ORDER BY RANDOM() LIMIT X`. Max limit = 100. |
| `GetFilters` | `GET /api/v1/filters` | `SELECT DISTINCT language FROM projects` + `SELECT DISTINCT unnest(topics)` |
| `Trigger` | `POST /api/v1/crawl` | Sends to scheduler channel, returns 202 Accepted or 429 if already running |
| `Status` | `GET /api/v1/status` | Returns `{ last_crawl, project_count, is_crawling }` |

**Config** (`config/config.go`):
- Loads env vars: `SUPABASE_URL`, `SUPABASE_KEY`, `GITHUB_TOKEN`, `CRAWL_INTERVAL`, `BACKEND_PORT`
- Auto-builds PostgreSQL connection string if `SUPABASE_URL` doesn't have a `postgres://` scheme
- Formats: `db.xxxxx.supabase.co` or `https://db.xxxxx.supabase.co` → `postgresql://postgres:KEY@host:5432/postgres?sslmode=require`

**Crawler** (`internal/crawler/`):

- **Scheduler:** `robfig/cron/v3` job every `CRAWL_INTERVAL` (default `6h`)
- **First run:** Immediately on startup
- **Trigger:** Channel-based manual trigger (`POST /api/v1/crawl`)
- **Concurrency guard:** Atomic boolean flag, skips if already running; returns 429 if crawl in progress
- **GitHub client:** `GET /search/repositories?q=stars:>1000&sort=stars&per_page=100`
  - Paginates through up to 10 pages
  - Uses `GITHUB_TOKEN` if available (5000 req/h), else anonymous (60 req/h)
- **Upsert:** `INSERT INTO projects ... ON CONFLICT (full_name) DO UPDATE SET ...`
- **Fields saved:** full_name, name, owner, description, stars, language, topics, html_url, avatar_url, last_updated

**DB layer** (`internal/db/supabase.go`):
- pgx connection pool with SSL
- Override DNS resolver to use Google DNS (8.8.8.8) for container compatibility
- Embedded SQL migration runner (`//go:embed migrations/*.sql`)
- Single migration: creates `projects` table + 3 indexes (language, stars DESC, GIN on topics)

**Startup sequence:**
1. Load config from env
2. Connect to Supabase (pgx pool)
3. Run migrations
4. Register routes
5. Start scheduler (triggers first crawl)
6. Listen on `:8000`
7. Graceful shutdown on SIGINT/SIGTERM

---

### 3. Database — Supabase (PostgreSQL)

**Location:** Managed cloud service

**Table: `projects`**

| Column | Type | Constraints | Notes |
|---|---|---|---|
| `id` | `BIGSERIAL` | `PRIMARY KEY` | Auto-increment |
| `full_name` | `TEXT` | `UNIQUE NOT NULL` | e.g. "golang/go" |
| `name` | `TEXT` | `NOT NULL` | e.g. "go" |
| `owner` | `TEXT` | `NOT NULL` | e.g. "golang" |
| `description` | `TEXT` | | May be null |
| `stars` | `INTEGER` | `DEFAULT 0` | Updated on each crawl |
| `language` | `TEXT` | | Main language |
| `topics` | `TEXT[]` | | PostgreSQL array of topic strings |
| `html_url` | `TEXT` | `NOT NULL` | Full GitHub URL |
| `avatar_url` | `TEXT` | | Owner avatar URL |
| `last_updated` | `TIMESTAMPTZ` | `DEFAULT NOW()` | Crawl timestamp |

**Indexes:**
- `idx_projects_language` on `language`
- `idx_projects_stars` on `stars DESC`
- `idx_projects_topics` GIN index on `topics` (supports `ANY` and `@>` operators)

**Key Query — Random Selection:**
```sql
SELECT * FROM projects
WHERE ($1 = '' OR LOWER(language) = LOWER($1))
  AND ($2 = '' OR $2 = ANY(topics))
  AND ($3 = 0 OR stars >= $3)
ORDER BY RANDOM()
LIMIT $4;
```

---

### 4. Redis (Future)

**Location:** Docker container / Upstash

**Purpose (planned):**
- Cache `GET /api/v1/filters` response (invalidated on crawl)
- Rate limiting for `/crawl` endpoint
- Optional: cache random selections (TTL of a few seconds)

---

### 5. Docker Compose (Local Dev)

```yaml
services:
  backend:
    build: ./backend
    ports: ["8000:8000"]
    env_file: .env
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports: ["3000:80"]
    depends_on: [backend]
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
    restart: unless-stopped
```

**nginx config** (frontend → backend proxy):
- `/api/` requests proxied to `http://backend:8000`
- All other routes → `index.html` (SPA fallback)

## CI/CD — GitHub Actions

### Frontend Deploy (`deploy-frontend.yml`)
- **Trigger:** push to `main` (changes in `frontend/` or workflow)
- **Steps:**
  1. `actions/checkout`, `actions/setup-node`, `npm ci`
  2. `npm run build` with `VITE_API_BASE_URL` (from secrets) and `VITE_BASE_URL=/github-wheel-of-random/`
  3. `actions/configure-pages`, `actions/upload-pages-artifact`, `actions/deploy-pages`

### Backend Deploy (`deploy-backend.yml`)
- **Trigger:** push to `main` (changes in `backend/` or workflow)
- **Steps:**
  1. `actions/checkout`, `docker/setup-buildx-action`
  2. Login to GHCR via `docker/login-action`
  3. Build and push image with tags `latest` + short sha
  4. (Optional) Trigger Render Deploy Hook via curl

## Communication

```
Frontend ──HTTP──▶ Backend ──PG wire──▶ Supabase
                      │
                      ├── GitHub API (crawler)
                      │
                      └── Redis (future, cache)
```

- Frontend never talks directly to Supabase or GitHub
- Backend is the single source of truth for data
- Frontend is pure static; no server-side rendering needed (GitHub Pages compatible)
