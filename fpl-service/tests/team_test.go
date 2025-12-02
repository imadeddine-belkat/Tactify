package main

import (
	"context"
	"log"
	"testing"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	"github.com/imadeddine-belkat/kafka"

	teamService "github.com/imadeddine-belkat/fpl-service/internal/services"
)

func TestTeamApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	// Setup service with real client
	service := &teamService.TeamApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}

	// Test with real API
	log.Println("Calling FPL API...")
	err := service.UpdateTeams(ctx)
	if err != nil {
		log.Fatalf("UpdateTeams with API failed: %v", err)
	}

	log.Println("Real API test completed successfully")
}
