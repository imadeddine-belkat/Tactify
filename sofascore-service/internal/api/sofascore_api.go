package sofascore_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/imadbelkat1/sofascore-service/config"
)

type SofascoreApiClient struct {
	Config         config.SofascoreConfig
	HttpClient     *http.Client
	UserAgent      string
	UseBrowserOnly bool
}

func NewSofascoreApiClient(cfg *config.SofascoreConfig) *SofascoreApiClient {
	return &SofascoreApiClient{
		Config:     *cfg,
		HttpClient: &http.Client{},
		UserAgent:  "Sofascore-Service-Client/1.0",
	}
}

// Get attempts HTTP request first, falls back to browser if 403
func (c *SofascoreApiClient) Get(ctx context.Context, endpoint string) ([]byte, error) {
	baseURL := strings.TrimRight(c.Config.SofascoreApi.BaseURL, "/")
	url := fmt.Sprintf("%s%s", baseURL, endpoint)
	log.Println("Request URL:", url)

	if c.UseBrowserOnly {
		return c.getWithBrowser(ctx, url)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: creating request: %w", err)
	}

	ua := c.UserAgent
	if ua == "" {
		ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Referer", "https://www.sofascore.com/")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: making request: %w", err)
	}
	defer resp.Body.Close()

	// If 403, fallback to browser method
	if resp.StatusCode == http.StatusForbidden {
		return c.getWithBrowser(ctx, url)
	}

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sofascore-api: API request failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	return io.ReadAll(resp.Body)
}

// getWithBrowser uses headless Chrome to bypass bot detection
func (c *SofascoreApiClient) getWithBrowser(ctx context.Context, url string) ([]byte, error) {
	// Create browser context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	browserCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	timeoutCtx, cancel := context.WithTimeout(browserCtx, 30*time.Second)
	defer cancel()

	var responseBody string
	var finalURL string

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Location(&finalURL),
		chromedp.Evaluate(`document.querySelector('pre') ? document.querySelector('pre').textContent : document.body.textContent`, &responseBody),
	)

	if err != nil {
		return nil, fmt.Errorf("sofascore-api: browser request failed: %w", err)
	}

	// Trim whitespace
	responseBody = strings.TrimSpace(responseBody)

	// Validate it's JSON
	if !strings.HasPrefix(responseBody, "{") && !strings.HasPrefix(responseBody, "[") {
		return nil, fmt.Errorf("sofascore-api: browser returned non-JSON content: %s", responseBody[:min(100, len(responseBody))])
	}

	return []byte(responseBody), nil
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (c *SofascoreApiClient) GetAndUnmarshal(ctx context.Context, endpoint string, target any) error {
	data, err := c.Get(ctx, endpoint)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, target); err != nil {
		log.Printf("Unmarshal error: %v", err)
		log.Printf("Full data causing error: %s", string(data))
		return fmt.Errorf("sofascore-api: unmarshaling response: %w", err)
	}

	return nil
}
