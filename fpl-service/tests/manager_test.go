package main

import (
	"context"
	"testing"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	. "github.com/imadbelkat1/fpl-service/internal/services"
	"github.com/imadbelkat1/kafka"
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

	err := service.UpdateManager(ctx, 2839296, 6)
	if err != nil {
		t.Fatalf("UpdateManager with API failed: %v", err)
	}

	t.Logf("Real API test completed successfully")

}
