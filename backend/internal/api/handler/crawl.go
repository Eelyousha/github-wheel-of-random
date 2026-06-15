package handler

import (
	"github.com/gofiber/fiber/v2"
)

type CrawlHandler struct {
	triggerChan chan struct{}
}

func NewCrawlHandler(triggerChan chan struct{}) *CrawlHandler {
	return &CrawlHandler{triggerChan: triggerChan}
}

func (h *CrawlHandler) Trigger(c *fiber.Ctx) error {
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
