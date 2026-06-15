# Git Workflow

## Branching

- `main` — stable, deployable
- `develop` — integration branch
- `feat/<name>` — feature branches (merged to `develop`)

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

## PR Workflow

1. Branch from `develop`
2. Implement changes
3. Open PR → `develop`
4. Squash-merge when approved
5. Merge `develop` → `main` for release

## Deploy Automation (Planned)

### Frontend → GitHub Pages
- Push to `main` → GitHub Action builds frontend (`npm ci && npm run build`)
- Deploys `dist/` to `gh-pages` branch
- Requires: `VITE_API_BASE_URL` as repository secret

### Backend → Railway / Fly.io
- Push to `main` → GitHub Action builds Docker image
- Pushes to container registry
- Triggers deploy on Railway or Fly.io (via webhook or CLI)

## Local Development

```bash
# Start all services
docker compose up --build

# Or run backend directly (for faster iteration)
cd backend && go run ./cmd/server/

# Frontend dev server with hot reload
cd frontend && npm run dev
```
