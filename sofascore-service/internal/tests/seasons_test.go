package tests

import (
	"context"
	"testing"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"github.com/imadeddine-belkat/sofascore-service/internal/services"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
)

func TestSeasonService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()

	service := &services.SeasonService{
		Config:   cfg,
		Client:   sofascore_api.NewSofascoreApiClient(cfg),
		Producer: kafka.NewProducer(),
		LeagueService: &services.LeagueService{
			Config:   cfg,
			Client:   sofascore_api.NewSofascoreApiClient(cfg),
			Producer: kafka.NewProducer(),
		},
	}

	premierLeagueID := cfg.SofascoreApi.LeaguesID.PremierLeague

	/*if err := service.UpdateAllLeaguesSeasons(ctx); err != nil {
		t.Fatal(err)
	}*/

	if err := service.UpdateLeaguesSeasons(ctx, premierLeagueID); err != nil {
		t.Fatal(err)
	}

	t.Log("Test completed successfully")
}
