package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"wheel-of-random/internal/api/handler"
	"wheel-of-random/internal/api/middleware"
)

func NewRouter(pool *pgxpool.Pool, crawlTrigger chan struct{}, statusHandler *handler.StatusHandler, adminKey string) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(middleware.CORS())

	projectHandler := handler.NewProjectHandler(pool)
	crawlHandler := handler.NewCrawlHandler(crawlTrigger, adminKey)

	api := app.Group("/api/v1")

	api.Get("/projects", projectHandler.List)
	api.Get("/projects/random", projectHandler.Random)
	api.Get("/filters", projectHandler.GetFilters)
	api.Post("/crawl", crawlHandler.Trigger)
	api.Get("/status", statusHandler.Status)

	return app
}
