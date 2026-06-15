package handler

import (
	"github.com/gofiber/fiber/v2"
)

type CrawlHandler struct {
	triggerChan chan struct{}
	adminKey    string
}

func NewCrawlHandler(triggerChan chan struct{}, adminKey string) *CrawlHandler {
	return &CrawlHandler{triggerChan: triggerChan, adminKey: adminKey}
}

func (h *CrawlHandler) Trigger(c *fiber.Ctx) error {
	if h.adminKey != "" && c.Query("key") != h.adminKey {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or missing admin key",
		})
	}

	select {
	case h.triggerChan <- struct{}{}:
		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "crawl triggered",
		})
	default:
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "crawl already in progress",
		})
	}
}
