package fpl_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/imadeddine-belkat/fpl-service/config"
	fplProto "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
)

type FplApiClient struct {
	Config     config.Config
	HttpClient *http.Client
	UserAgent  string
}

func NewFplApiClient(cfg *config.Config) *FplApiClient {
	return &FplApiClient{
		Config:     *cfg,
		HttpClient: &http.Client{},
		UserAgent:  "FPL-Service-Client/1.0",
	}
}

func (c *FplApiClient) GetPlayersBootstrap(ctx context.Context) (*fplProto.PlayersBootstrap, error) {
	bootstrap, err := c.getBootstrapData(ctx)
	if err != nil {
		return nil, err
	}

	return &fplProto.PlayersBootstrap{Elements: bootstrap.GetElements()}, nil
}

func (c *FplApiClient) getBootstrapData(ctx context.Context) (*fplProto.BootstrapResponse, error) {
	var bootstrap fplProto.BootstrapResponse
	endpoint := c.Config.FplApi.Bootstrap

	if err := c.GetAndUnmarshal(ctx, endpoint, &bootstrap); err != nil {
		return nil, err
	}
	return &bootstrap, nil
}

func (c *FplApiClient) GetAndUnmarshal(ctx context.Context, endpoint string, result any) error {
	data, err := c.Get(ctx, endpoint)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("fpl-api: unmarshaling response: %w", err)
	}

	return nil
}

func (c *FplApiClient) Get(ctx context.Context, endpoint string) ([]byte, error) {
	baseUrl := c.Config.FplApi.BaseUrl

	url := fmt.Sprintf("%s%s", baseUrl, endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("fpl-api: creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fpl-api: making request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("fpl-api: Error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fpl-api: API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fpl-api: reading response body: %w", err)
	}

	return body, nil
}
