package tests

import (
	"context"
	"testing"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"github.com/imadeddine-belkat/sofascore-service/internal/services"
)

func TestMatchLineupService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()
	service := &services.MatchLineupService{
		Event: &services.EventsService{
			Config: cfg,
			Client: sofascore_api.NewSofascoreApiClient(cfg),
		},
		Config:   cfg,
		Client:   sofascore_api.NewSofascoreApiClient(cfg),
		Producer: kafka.NewProducer(),
	}

	ctx := context.Background()
	leagues := map[string]int{
		"LALIGA":        cfg.SofascoreApi.LeaguesID.LaLiga,
		"PREMIERLEAGUE": cfg.SofascoreApi.LeaguesID.PremierLeague,
	}

	for leagueName, leagueID := range leagues {
		for _, seasonID := range cfg.AllSeasons(leagueName) {
			for round := 1; round <= 38; round++ {
				t.Logf("Fetching %s SeasonID: %d, Round: %d", leagueName, seasonID, round)
				if err := service.UpdatePlayersStats(ctx, seasonID, leagueID, round); err != nil {
					t.Fatalf("Error: %v", err)
				}
			}
		}
	}

	t.Log("Test completed successfully")
}
