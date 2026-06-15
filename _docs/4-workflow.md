# Git Workflow

## Branching

- `main` — stable, deployable
- `feat/<name>` — feature branches (merged to `main`)

## Commit Convention

```
<type>(<scope>): <message>
```

Types: `feat` / `fix` / `chore` / `docs` / `refactor` / `test`

Examples:
```
feat(crawler): implement GitHub search client
fix(api): clamp random limit to max 100
docs(plan): update milestones with completed tasks
```

## CI/CD Pipeline

### Frontend → GitHub Pages
- **Trigger:** Push to `main` with changes in `frontend/**` or `.github/workflows/deploy-frontend.yml`
- **Secret needed:** `VITE_API_BASE_URL` (production backend URL, e.g. `https://my-app.onrender.com`)
- **Env passed at build time:** `VITE_BASE_URL=/github-wheel-of-random/`

### Backend → GHCR + Render.com
- **Trigger:** Push to `main` with changes in `backend/**` or `.github/workflows/deploy-backend.yml`
- Builds Docker image → pushes to `ghcr.io/<owner>/github-wheel-of-random-backend:latest`
- (Optional) Triggers Render Deploy Hook if `RENDER_DEPLOY_HOOK_URL` secret is set

### Manual Runs
Both workflows support `workflow_dispatch` — can be triggered manually from GitHub Actions UI.

## Local Development

```bash
# Start all services
docker compose up --build

# Or run backend directly (for faster iteration)
cd backend && go run ./cmd/server/

# Frontend dev server with hot reload
cd frontend && npm run dev
```

## Environment Setup

1. Copy `.env.example` → `.env`
2. Fill in `SUPABASE_URL`, `SUPABASE_KEY`, optionally `GITHUB_TOKEN`
3. For frontend dev, `VITE_API_BASE_URL=http://localhost:8000`
4. Run with Docker Compose or individual dev servers
