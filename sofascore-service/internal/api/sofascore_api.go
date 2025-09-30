package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/imadbelkat1/sofascore-service/config"
)

type SofascoreApiClient struct {
	Config     *config.Config
	HttpClient *http.Client
	UserAgent  string
}

func NewSofascoreApiClient() *SofascoreApiClient {
	return &SofascoreApiClient{
		HttpClient: &http.Client{},
		UserAgent:  "Sofascore-Service-Client/1.0",
	}
}

func (c *SofascoreApiClient) Get(ctx context.Context, endpoint string) ([]byte, error) {
	// Load base URL from config
	baseURL := c.Config.SofascoreApi.BaseURL

	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("sofascore-api: Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sofascore-api: API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("sofascore-api: reading response body: %w", err)
	}

	return body, nil
}

func (c *SofascoreApiClient) GetAndUnmarshal(ctx context.Context, endpoint string, target any) error {
	data, err := c.Get(ctx, endpoint)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("sofascore-api: unmarshaling response: %w", err)
	}

	return nil
}
