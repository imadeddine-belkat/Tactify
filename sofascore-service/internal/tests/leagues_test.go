package tests

import (
	"context"
	"testing"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"github.com/imadeddine-belkat/sofascore-service/internal/services"
)

func TestLeagues(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()

	service := &services.LeagueService{
		Config:   cfg,
		Client:   sofascore_api.NewSofascoreApiClient(cfg),
		Producer: kafka.NewProducer(),
	}

	if err := service.UpdateLeagueIDs(ctx); err != nil {
		t.Fatal(err)
	}

	t.Log("Test completed successfully")
}
