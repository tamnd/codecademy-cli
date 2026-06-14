// Package codecademy scrapes the Codecademy course catalog.
package codecademy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Config controls the HTTP client behaviour.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://www.codecademy.com",
		Rate:      300 * time.Millisecond,
		Timeout:   30 * time.Second,
		Retries:   3,
		UserAgent: "codecademy-cli/0.1 (github.com/tamnd/codecademy-cli)",
	}
}

// Client fetches Codecademy data.
type Client struct {
	cfg     Config
	http    *http.Client
	last    time.Time
	buildID string
}

// NewClient creates a Client from cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

var buildIDRe = regexp.MustCompile(`"buildId"\s*:\s*"([^"]+)"`)

type catalogResponse struct {
	PageProps struct {
		InitialCatalogResults struct {
			TotalResults int          `json:"totalResults"`
			TotalPages   int          `json:"totalPages"`
			PageItems    []catalogItem `json:"pageItems"`
		} `json:"initialCatalogResults"`
	} `json:"pageProps"`
}

type catalogItem struct {
	Slug              string `json:"slug"`
	Title             string `json:"title"`
	Type              string `json:"type"`
	Difficulty        string `json:"difficulty"`
	LessonCount       int    `json:"lessonCount"`
	TimeToComplete    int    `json:"timeToComplete"`
	Pro               bool   `json:"pro"`
	GrantsCertificate bool   `json:"grantsCertificate"`
	ShortDescription  string `json:"shortDescription"`
	URLPath           string `json:"urlPath"`
}

// ListCourses fetches all catalog items across all pages.
func (c *Client) ListCourses(ctx context.Context, limit int) ([]Course, error) {
	if err := c.ensureBuildID(ctx); err != nil {
		return nil, err
	}
	var all []Course
	rank := 1
	for page := 1; ; page++ {
		items, total, err := c.fetchPage(ctx, page)
		if err != nil {
			return nil, err
		}
		for _, it := range items {
			url := c.cfg.BaseURL + it.URLPath
			if it.URLPath == "" {
				url = c.cfg.BaseURL + "/learn/" + it.Slug
			}
			all = append(all, Course{
				Rank:             rank,
				Slug:             it.Slug,
				Title:            it.Title,
				Type:             it.Type,
				Difficulty:       it.Difficulty,
				LessonCount:      it.LessonCount,
				TimeToComplete:   it.TimeToComplete, // hours
				Pro:              it.Pro,
				Certificate:      it.GrantsCertificate,
				ShortDescription: it.ShortDescription,
				URL:              url,
			})
			rank++
			if limit > 0 && len(all) >= limit {
				return all, nil
			}
		}
		if rank > total || len(items) == 0 {
			break
		}
	}
	return all, nil
}

// Search returns courses whose title or slug matches query (case-insensitive).
func (c *Client) Search(ctx context.Context, query string, limit int) ([]Course, error) {
	all, err := c.ListCourses(ctx, 0)
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var out []Course
	for _, course := range all {
		if strings.Contains(strings.ToLower(course.Title), q) ||
			strings.Contains(strings.ToLower(course.Slug), q) ||
			strings.Contains(strings.ToLower(course.ShortDescription), q) {
			out = append(out, course)
			if limit > 0 && len(out) >= limit {
				break
			}
		}
	}
	for i := range out {
		out[i].Rank = i + 1
	}
	return out, nil
}

func (c *Client) ensureBuildID(ctx context.Context) error {
	if c.buildID != "" {
		return nil
	}
	body, err := c.fetch(ctx, "/catalog")
	if err != nil {
		return fmt.Errorf("discovering build ID: %w", err)
	}
	m := buildIDRe.FindSubmatch(body)
	if m == nil {
		return fmt.Errorf("could not find Next.js buildId in /catalog")
	}
	c.buildID = string(m[1])
	return nil
}

func (c *Client) fetchPage(ctx context.Context, page int) ([]catalogItem, int, error) {
	path := fmt.Sprintf("/_next/data/%s/catalog.json?page=%d", c.buildID, page)
	body, err := c.fetch(ctx, path)
	if err != nil {
		return nil, 0, err
	}
	var resp catalogResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, 0, fmt.Errorf("parsing catalog page %d: %w", page, err)
	}
	r := resp.PageProps.InitialCatalogResults
	return r.PageItems, r.TotalResults, nil
}

func (c *Client) fetch(ctx context.Context, path string) ([]byte, error) {
	url := c.cfg.BaseURL + path
	var last error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}
		c.pace()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", c.cfg.UserAgent)
		resp, err := c.http.Do(req)
		if err != nil {
			last = err
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			last = fmt.Errorf("HTTP %d", resp.StatusCode)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
		}
		return io.ReadAll(resp.Body)
	}
	return nil, fmt.Errorf("all retries failed for %s: %w", url, last)
}

func (c *Client) pace() {
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}
