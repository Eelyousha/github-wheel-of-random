package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GitHubRepo struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Name        string `json:"name"`
	Owner       struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	Description *string  `json:"description"`
	Stars       int      `json:"stargazers_count"`
	Language    *string  `json:"language"`
	Topics      []string `json:"topics"`
	HTMLURL     string   `json:"html_url"`
}

type searchResponse struct {
	Items      []GitHubRepo `json:"items"`
	TotalCount int          `json:"total_count"`
}

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}
}

func (c *Client) Search(ctx context.Context, query string, perPage, page int) ([]GitHubRepo, int, error) {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=%d&page=%d", query, perPage, page)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "wheel-of-random/1.0")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("GitHub API status %d: %s", resp.StatusCode, string(body))
	}

	var sr searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}

	return sr.Items, sr.TotalCount, nil
}

func (c *Client) SearchTopRepos(ctx context.Context, minStars int, maxPages, perPage int) ([]GitHubRepo, error) {
	query := fmt.Sprintf("stars:>%d", minStars)
	var all []GitHubRepo

	for page := 1; page <= maxPages; page++ {
		select {
		case <-ctx.Done():
			return all, ctx.Err()
		default:
		}

		items, total, err := c.Search(ctx, query, perPage, page)
		if err != nil {
			return all, fmt.Errorf("page %d: %w", page, err)
		}

		all = append(all, items...)

		// Stop if total items are exhausted
		if page*perPage >= total || len(items) < perPage {
			break
		}
	}

	return all, nil
}
