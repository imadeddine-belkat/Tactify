package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	playerService "github.com/imadbelkat1/fpl-service/internal/services"
	"github.com/imadbelkat1/kafka"
)

func TestPlayersApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	service := &playerService.PlayerApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}

	time.Sleep(200 * time.Millisecond)

	log.Println("Calling FPL API...")
	start := time.Now()
	err := service.UpdatePlayers(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("UpdatePlayer with API failed: %v", err)
	}

	log.Printf("Publishing completed in: %v", elapsed)
	t.Log("Real API test completed successfully")
}
