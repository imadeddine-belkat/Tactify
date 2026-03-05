package main

import (
	"context"
	"log"
	"testing"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	live_event_service "github.com/imadeddine-belkat/fpl-service/internal/services"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
)

func TestLiveEventApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	producer := kafka.NewProducer()
	defer producer.Close()

	// Setup service with real client
	service := &live_event_service.LiveEventApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: producer,
	}

	// Test with real API
	log.Println("Calling FPL API...")
	for i := 1; i <= 29; i++ {
		err := service.UpdateLiveEvent(ctx, i)
		if err != nil {
			t.Fatalf("UpdateLiveEvent(Test) with API failed: %v", err)
		}
	}

	t.Log("Real API test completed successfully")
}
