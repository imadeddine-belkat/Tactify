package main

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	playerService "github.com/imadeddine-belkat/fpl-service/internal/services"
	"github.com/imadeddine-belkat/kafka"
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
		log.Fatalf("UpdatePlayer with API failed: %v", err)
	}

	log.Printf("Publishing completed in: %v", elapsed)
	log.Printf("Real API test completed successfully")
}
