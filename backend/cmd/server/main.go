package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"wheel-of-random/config"
	"wheel-of-random/internal/api"
	"wheel-of-random/internal/api/handler"
	"wheel-of-random/internal/crawler"
	"wheel-of-random/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.Connect(ctx, cfg.SupabaseURL)
	if err != nil {
		log.Fatalf("db.Connect: %v", err)
	}
	defer database.Close()

	if err := database.RunMigrations(ctx); err != nil {
		log.Fatalf("db.RunMigrations: %v", err)
	}

	statusH := handler.NewStatusHandler(database.Pool)
	crawlTrigger := make(chan struct{}, 1)

	scheduler := crawler.NewScheduler(
		database.Pool,
		cfg.GitHubToken,
		cfg.CrawlInterval,
		crawlTrigger,
		func() { statusH.SetCrawling(true) },
		func() { statusH.SetCrawling(false) },
	)
	scheduler.Start(ctx)
	defer scheduler.Stop()

	app := api.NewRouter(database.Pool, crawlTrigger, statusH, cfg.AdminKey)

	go func() {
		if err := app.Listen(":" + cfg.BackendPort); err != nil {
			log.Fatalf("app.Listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	app.Shutdown()
	cancel()
}
