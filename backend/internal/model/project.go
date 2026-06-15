package model

import "time"

type Project struct {
	ID          int64     `json:"id"`
	FullName    string    `json:"full_name"`
	Name        string    `json:"name"`
	Owner       string    `json:"owner"`
	Description *string   `json:"description"`
	Stars       int       `json:"stars"`
	Language    *string   `json:"language"`
	Topics      []string  `json:"topics"`
	HTMLURL     string    `json:"html_url"`
	AvatarURL   *string   `json:"avatar_url"`
	LastUpdated time.Time `json:"last_updated"`
}

type ProjectFilter struct {
	Language string   `json:"language"`
	Topics   []string `json:"topics"`
	MinStars int      `json:"min_stars"`
	Limit    int      `json:"limit"`
	Page     int      `json:"page"`
}

type Status struct {
	LastCrawl    *time.Time `json:"last_crawl"`
	ProjectCount int        `json:"project_count"`
	IsCrawling   bool       `json:"is_crawling"`
}

type Filters struct {
	Languages []string `json:"languages"`
	Topics    []string `json:"topics"`
}
