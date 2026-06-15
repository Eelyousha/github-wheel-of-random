package crawler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron        *cron.Cron
	jobID       cron.EntryID
	github      *Client
	pool        *pgxpool.Pool
	interval    time.Duration
	triggerChan chan struct{}
	onStart     func()
	onFinish    func()
	crawling    bool
}

func NewScheduler(pool *pgxpool.Pool, token string, interval time.Duration, triggerChan chan struct{}, onStart, onFinish func()) *Scheduler {
	return &Scheduler{
		cron:        cron.New(),
		github:      NewClient(token),
		pool:        pool,
		interval:    interval,
		triggerChan: triggerChan,
		onStart:     onStart,
		onFinish:    onFinish,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	cronExpr := fmt.Sprintf("@every %s", s.interval.String())
	id, err := s.cron.AddFunc(cronExpr, func() {
		s.crawl(ctx)
	})
	if err != nil {
		log.Printf("scheduler: failed to add cron job: %v", err)
		return
	}
	s.jobID = id
	s.cron.Start()

	// Listen for manual triggers
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-s.triggerChan:
				s.crawl(ctx)
			}
		}
	}()

	log.Printf("scheduler: started, interval=%s", s.interval)
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) crawl(ctx context.Context) {
	if s.crawling {
		log.Println("crawler: already running, skipping")
		return
	}
	s.crawling = true
	if s.onStart != nil {
		s.onStart()
	}
	defer func() {
		s.crawling = false
		if s.onFinish != nil {
			s.onFinish()
		}
	}()

	log.Println("crawler: starting...")

	repos, err := s.github.SearchTopRepos(ctx, 1000, 10, 100)
	if err != nil {
		log.Printf("crawler: search failed: %v", err)
		return
	}

	log.Printf("crawler: fetched %d repos", len(repos))

	saved := 0
	for _, repo := range repos {
		_, err := s.pool.Exec(ctx, `
			INSERT INTO projects (full_name, name, owner, description, stars, language, topics, html_url, avatar_url, last_updated)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
			ON CONFLICT (full_name) DO UPDATE SET
				description = EXCLUDED.description,
				stars = EXCLUDED.stars,
				language = EXCLUDED.language,
				topics = EXCLUDED.topics,
				avatar_url = EXCLUDED.avatar_url,
				last_updated = NOW()
		`, repo.FullName, repo.Name, repo.Owner.Login, repo.Description, repo.Stars, repo.Language, repo.Topics, repo.HTMLURL, repo.Owner.AvatarURL)
		if err != nil {
			log.Printf("crawler: upsert error for %s: %v", repo.FullName, err)
			continue
		}
		saved++
	}

	log.Printf("crawler: done, saved/updated %d projects", saved)
}
