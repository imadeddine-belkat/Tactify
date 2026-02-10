package sofascore_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/imadeddine-belkat/sofascore-service/config"
)

type SofascoreApiClient struct {
	Config         config.SofascoreConfig
	HttpClient     *http.Client
	UserAgent      string
	UseBrowserOnly bool
}

func NewSofascoreApiClient(cfg *config.SofascoreConfig) *SofascoreApiClient {
	client := &SofascoreApiClient{
		Config:    *cfg,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
		HttpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}

	log.Println("✓ Initialized Direct Client (Using local IP)")
	return client
}

func (c *SofascoreApiClient) Get(ctx context.Context, endpoint string) ([]byte, error) {
	baseURL := strings.TrimRight(c.Config.SofascoreApi.BaseURL, "/")
	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	// If browser-only mode is forced in config
	if c.UseBrowserOnly {
		return c.getWithBrowser(ctx, url)
	}

	// 1. Try Standard HTTP Request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	c.setHeaders(req)

	resp, err := c.HttpClient.Do(req)

	// 2. Check for blocks or errors
	// If we get blocked (403/429) or service unavailable (503), switch to Browser
	if err != nil || resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode == http.StatusServiceUnavailable {

		if resp != nil {
			log.Printf("⚠ HTTP Request blocked (Status: %d). Falling back to Browser...", resp.StatusCode)
			resp.Body.Close()
		} else {
			log.Printf("⚠ HTTP Request failed (%v). Falling back to Browser...", err)
		}

		return c.getWithBrowser(ctx, url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("api status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *SofascoreApiClient) getWithBrowser(ctx context.Context, url string) ([]byte, error) {
	// Add a small random delay to look more human
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // Set to false if you want to see the browser open
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent(c.UserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancel := context.WithTimeout(browserCtx, 45*time.Second)
	defer cancel()

	var responseBody string

	// Navigate and extract JSON from <pre> tag (common in raw JSON views) or body
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Evaluate(`document.querySelector('pre') ? document.querySelector('pre').textContent : document.body.textContent`, &responseBody),
	)

	if err != nil {
		return nil, fmt.Errorf("browser error: %w", err)
	}

	return []byte(strings.TrimSpace(responseBody)), nil
}

func (c *SofascoreApiClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://www.sofascore.com/")
	req.Header.Set("Origin", "https://www.sofascore.com")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	// Random sleep to mimic human behavior
	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
}

func (c *SofascoreApiClient) GetAndUnmarshal(ctx context.Context, endpoint string, target any) error {
	data, err := c.Get(ctx, endpoint)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}
