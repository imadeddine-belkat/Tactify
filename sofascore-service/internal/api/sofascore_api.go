package sofascore_api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/imadbelkat1/sofascore-service/config"
	"golang.org/x/net/proxy"
)

type SofascoreApiClient struct {
	Config         config.SofascoreConfig
	HttpClient     *http.Client
	UserAgent      string
	UseBrowserOnly bool

	// Tor circuit rotation rate limiting
	lastRotation time.Time
	rotationMu   sync.Mutex
}

func NewSofascoreApiClient(cfg *config.SofascoreConfig) *SofascoreApiClient {
	client := &SofascoreApiClient{
		Config:    *cfg,
		UserAgent: "Sofascore-Service-Client/1.0",
	}

	// Configure HTTP client based on Tor settings
	if cfg.Tor.Enabled {
		httpClient, err := createTorClient(cfg.Tor)
		if err != nil {
			log.Fatalf("Failed to create Tor client: %v", err)
		}
		client.HttpClient = httpClient
		log.Println("✓ Tor client initialized (SOCKS5 proxy active)")
	} else {
		client.HttpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
		log.Println("✓ Direct HTTP client initialized (no Tor)")
	}

	return client
}

// createTorClient builds an HTTP client that routes through Tor's SOCKS5 proxy
func createTorClient(torCfg config.TorConfig) (*http.Client, error) {
	// Create SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", torCfg.SocksAddr, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("creating SOCKS5 dialer: %w", err)
	}

	// Custom transport with SOCKS5
	transport := &http.Transport{
		Dial:                dialer.Dial,
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		DisableKeepAlives:   false,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Tor is slower
	}, nil
}

// RotateCircuit forces Tor to use a new exit node (new IP)
// Rate-limited to once per 10 seconds (Tor's limitation)
func (c *SofascoreApiClient) RotateCircuit(ctx context.Context) error {
	if !c.Config.Tor.Enabled {
		return fmt.Errorf("Tor is not enabled in config")
	}

	c.rotationMu.Lock()
	defer c.rotationMu.Unlock()

	// Tor rate-limits NEWNYM to once per 10 seconds
	timeSinceLastRotation := time.Since(c.lastRotation)
	if timeSinceLastRotation < 10*time.Second {
		waitTime := 10*time.Second - timeSinceLastRotation
		log.Printf("⏳ Waiting %.1fs before rotating (Tor rate limit)...", waitTime.Seconds())
		time.Sleep(waitTime)
	}

	controlAddr := c.Config.Tor.ControlAddr
	password := c.Config.Tor.Password

	// Connect to Tor control port
	conn, err := net.DialTimeout("tcp", controlAddr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("connecting to Tor control port: %w", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Authenticate (if password is set)
	if password != "" {
		authCmd := fmt.Sprintf("AUTHENTICATE \"%s\"\r\n", password)
		if _, err := conn.Write([]byte(authCmd)); err != nil {
			return fmt.Errorf("sending AUTHENTICATE: %w", err)
		}

		resp, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading auth response: %w", err)
		}

		if !strings.HasPrefix(resp, "250") {
			return fmt.Errorf("authentication failed: %s", strings.TrimSpace(resp))
		}
	}

	// Send NEWNYM signal to rotate circuit
	if _, err := conn.Write([]byte("SIGNAL NEWNYM\r\n")); err != nil {
		return fmt.Errorf("sending NEWNYM signal: %w", err)
	}

	resp, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("reading NEWNYM response: %w", err)
	}

	if !strings.HasPrefix(resp, "250") {
		return fmt.Errorf("NEWNYM failed: %s", strings.TrimSpace(resp))
	}

	c.lastRotation = time.Now()
	log.Println("✓ Tor circuit rotated (new IP)")

	// Wait for Tor to establish new circuit
	time.Sleep(2 * time.Second)

	return nil
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

	// If 403 and Tor is enabled, rotate and retry ONCE
	if resp.StatusCode == http.StatusForbidden {
		if c.Config.Tor.Enabled {
			log.Println("⚠ 403 detected, rotating Tor circuit...")
			if err := c.RotateCircuit(ctx); err != nil {
				log.Printf("Circuit rotation failed: %v, falling back to browser", err)
			} else {
				// Retry request with new IP (ONE retry only to avoid infinite loops)
				log.Println("🔄 Retrying request with new IP...")
				return c.retryGet(ctx, url)
			}
		}

		// Fallback to browser if not using Tor or rotation failed
		log.Println("→ Falling back to browser mode")
		return c.getWithBrowser(ctx, url)
	}

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sofascore-api: API request failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	return io.ReadAll(resp.Body)
}

// retryGet performs a single retry without recursive rotation
func (c *SofascoreApiClient) retryGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: creating retry request: %w", err)
	}

	// Set same headers
	ua := c.UserAgent
	if ua == "" {
		ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://www.sofascore.com/")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: retry request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("sofascore-api: retry failed with status %d: %s", resp.StatusCode, string(errBody))
	}

	return io.ReadAll(resp.Body)
}

// getWithBrowser uses headless Chrome to bypass bot detection
func (c *SofascoreApiClient) getWithBrowser(ctx context.Context, url string) ([]byte, error) {
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

	responseBody = strings.TrimSpace(responseBody)

	if !strings.HasPrefix(responseBody, "{") && !strings.HasPrefix(responseBody, "[") {
		return nil, fmt.Errorf("sofascore-api: browser returned non-JSON content: %s", responseBody[:min(100, len(responseBody))])
	}

	return []byte(responseBody), nil
}

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
