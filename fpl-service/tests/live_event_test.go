package main

import (
	"context"
	"log"
	"testing"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	live_event_service "github.com/imadbelkat1/fpl-service/internal/services"
	"github.com/imadbelkat1/kafka"
)

func TestLiveEventApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	// Setup service with real client
	service := &live_event_service.LiveEventApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}

	// Test with real API
	log.Println("Calling FPL API...")
	err := service.UpdateLiveEvent(ctx, 7)
	if err != nil {
		t.Fatalf("UpdateLiveEvent(Test) with API failed: %v", err)
	}

	t.Log("Real API test completed successfully")
}
