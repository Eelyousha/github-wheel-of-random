package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/exp/rand"
)

type ProjectHandler struct {
	pool *pgxpool.Pool
}

func NewProjectHandler(pool *pgxpool.Pool) *ProjectHandler {
	return &ProjectHandler{pool: pool}
}

func (h *ProjectHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	language := c.Query("lang")
	topic := c.Query("topic")
	minStars, _ := strconv.Atoi(c.Query("min_stars", "0"))

	if limit < 1 || limit > 100 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}

	where := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if language != "" {
		where = append(where, fmt.Sprintf("LOWER(language) = LOWER($%d)", argIdx))
		args = append(args, language)
		argIdx++
	}
	if topic != "" {
		where = append(where, fmt.Sprintf("$%d = ANY(topics)", argIdx))
		args = append(args, topic)
		argIdx++
	}
	if minStars > 0 {
		where = append(where, fmt.Sprintf("stars >= $%d", argIdx))
		args = append(args, minStars)
		argIdx++
	}

	offset := (page - 1) * limit

	query := fmt.Sprintf(
		`SELECT id, full_name, name, owner, description, stars, language, topics, html_url, avatar_url, last_updated
		 FROM projects WHERE %s ORDER BY stars DESC LIMIT $%d OFFSET $%d`,
		strings.Join(where, " AND "), argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	projects, err := scanProjects(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(projects)
}

func (h *ProjectHandler) Random(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	language := c.Query("lang")
	topic := c.Query("topic")
	minStars, _ := strconv.Atoi(c.Query("min_stars", "0"))

	if limit < 1 || limit > 100 {
		limit = 20
	}

	where := []string{"1=1"}
	args := []any{}
	argIdx := 1

	if language != "" {
		where = append(where, fmt.Sprintf("LOWER(language) = LOWER($%d)", argIdx))
		args = append(args, language)
		argIdx++
	}
	if topic != "" {
		where = append(where, fmt.Sprintf("$%d = ANY(topics)", argIdx))
		args = append(args, topic)
		argIdx++
	}
	if minStars > 0 {
		where = append(where, fmt.Sprintf("stars >= $%d", argIdx))
		args = append(args, minStars)
		argIdx++
	}

	query := fmt.Sprintf(
		`SELECT id, full_name, name, owner, description, stars, language, topics, html_url, avatar_url, last_updated
		 FROM projects WHERE %s ORDER BY RANDOM() LIMIT $%d`,
		strings.Join(where, " AND "), argIdx,
	)
	args = append(args, limit)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	projects, err := scanProjects(rows)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(projects)
}

type Filters struct {
	Languages []string `json:"languages"`
	Topics    []string `json:"topics"`
}

func (h *ProjectHandler) GetFilters(c *fiber.Ctx) error {
	f := Filters{}

	rows, err := h.pool.Query(context.Background(), "SELECT DISTINCT language FROM projects WHERE language IS NOT NULL ORDER BY language")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	for rows.Next() {
		var lang string
		if err := rows.Scan(&lang); err == nil {
			f.Languages = append(f.Languages, lang)
		}
	}

	rows2, err := h.pool.Query(context.Background(), "SELECT DISTINCT unnest(topics) AS topic FROM projects ORDER BY topic")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows2.Close()

	for rows2.Next() {
		var topic string
		if err := rows2.Scan(&topic); err == nil {
			f.Topics = append(f.Topics, topic)
		}
	}

	return c.JSON(f)
}

func scanProjects(rows pgx.Rows) ([]map[string]any, error) {
	var projects []map[string]any
	for rows.Next() {
		var (
			id          int64
			fullName    string
			name        string
			owner       string
			description *string
			stars       int
			language    *string
			topics      []string
			htmlURL     string
			avatarURL   *string
			lastUpdated time.Time
		)
		if err := rows.Scan(&id, &fullName, &name, &owner, &description, &stars, &language, &topics, &htmlURL, &avatarURL, &lastUpdated); err != nil {
			return nil, err
		}

		// shuffle topics to make the response non-deterministic
		rand.Shuffle(len(topics), func(i, j int) {
			topics[i], topics[j] = topics[j], topics[i]
		})

		projects = append(projects, map[string]any{
			"id":           id,
			"full_name":    fullName,
			"name":         name,
			"owner":        owner,
			"description":  description,
			"stars":        stars,
			"language":     language,
			"topics":       topics,
			"html_url":     htmlURL,
			"avatar_url":   avatarURL,
			"last_updated": lastUpdated,
		})
	}
	return projects, nil
}
