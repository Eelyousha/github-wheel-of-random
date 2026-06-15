package handler

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusHandler struct {
	pool       *pgxpool.Pool
	mu         sync.RWMutex
	lastCrawl  *time.Time
	isCrawling bool
}

func NewStatusHandler(pool *pgxpool.Pool) *StatusHandler {
	return &StatusHandler{pool: pool}
}

func (h *StatusHandler) SetCrawling(crawling bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.isCrawling = crawling
	if !crawling {
		now := time.Now()
		h.lastCrawl = &now
	}
}

func (h *StatusHandler) Status(c *fiber.Ctx) error {
	h.mu.RLock()
	isCrawling := h.isCrawling
	lastCrawl := h.lastCrawl
	h.mu.RUnlock()

	var count int
	err := h.pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		count = 0
	}

	return c.JSON(fiber.Map{
		"last_crawl":    lastCrawl,
		"project_count": count,
		"is_crawling":   isCrawling,
	})
}
