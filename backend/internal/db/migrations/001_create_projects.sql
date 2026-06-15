CREATE TABLE IF NOT EXISTS projects (
    id          BIGSERIAL PRIMARY KEY,
    full_name   TEXT UNIQUE NOT NULL,
    name        TEXT NOT NULL,
    owner       TEXT NOT NULL,
    description TEXT,
    stars       INTEGER NOT NULL DEFAULT 0,
    language    TEXT,
    topics      TEXT[],
    html_url    TEXT NOT NULL,
    avatar_url  TEXT,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_projects_language ON projects (language);
CREATE INDEX IF NOT EXISTS idx_projects_stars ON projects (stars DESC);
CREATE INDEX IF NOT EXISTS idx_projects_topics ON projects USING GIN (topics);
