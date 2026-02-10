package main

import (
	"context"
	"testing"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	. "github.com/imadeddine-belkat/fpl-service/internal/services"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
)

func TestManagersApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	service := &ManagersApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}
	for n := 1; n < 26; n++ {
		t.Logf("Waiting for API to be ready... (attempt %d)", n+1)
		err := service.UpdateManager(ctx, 2839296, n)
		if err != nil {
			t.Fatalf("UpdateManager with API failed: %v", err)
		}
	}

	t.Logf("Real API test completed successfully")

}
