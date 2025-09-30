package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	fixutreService "github.com/imadbelkat1/fpl-service/internal/services"
	"github.com/imadbelkat1/kafka"
)

func TestFixturesApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	service := &fixutreService.FixturesApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}

	time.Sleep(200 * time.Millisecond)

	log.Println("Calling FPL API...")
	start := time.Now()
	err := service.UpdateFixtures(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("UpdateFixtures with API failed: %v", err)
	}

	log.Printf("Publishing completed in: %v", elapsed)
	t.Log("Real API test completed successfully")
}
