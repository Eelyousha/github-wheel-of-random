package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	SupabaseURL   string
	SupabaseKey   string
	GitHubToken   string
	BackendPort   string
	CrawlInterval time.Duration
}

func Load() (*Config, error) {
	intervalStr := getEnv("CRAWL_INTERVAL", "6h")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CRAWL_INTERVAL: %w", err)
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	// If SUPABASE_URL doesn't have a scheme, build a proper postgres connection string
	if !strings.HasPrefix(supabaseURL, "postgres://") && !strings.HasPrefix(supabaseURL, "postgresql://") {
		supabaseURL = buildConnectionString(supabaseURL, supabaseKey)
	}

	return &Config{
		SupabaseURL:   supabaseURL,
		SupabaseKey:   supabaseKey,
		GitHubToken:   getEnv("GITHUB_TOKEN", ""),
		BackendPort:   getEnv("BACKEND_PORT", "8000"),
		CrawlInterval: interval,
	}, nil
}

// buildConnectionString constructs a postgres:// URL from Supabase project URL and key.
// Accepts formats:
//
//	db.xxxxx.supabase.co
//	https://db.xxxxx.supabase.co
func buildConnectionString(projectURL, password string) string {
	host := strings.TrimPrefix(projectURL, "https://")
	host = strings.TrimSuffix(host, "/")

	// Default Supabase connection pooler port
	u := url.URL{
		Scheme: "postgresql",
		Host:   host + ":5432",
		Path:   "/postgres",
	}

	if password != "" {
		u.User = url.UserPassword("postgres", password)
	}

	// Append SSL mode
	u.RawQuery = "sslmode=require"

	return u.String()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
