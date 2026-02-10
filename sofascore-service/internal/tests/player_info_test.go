package tests

import (
	"context"
	"testing"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"github.com/imadeddine-belkat/sofascore-service/internal/services"
)

func TestPlayersInfoService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()
	apiClient := sofascore_api.NewSofascoreApiClient(cfg)
	service := &services.PlayersService{
		Config: *cfg,
		Client: apiClient,
		Standing: &services.LeagueStandingService{
			Config:   cfg,
			Client:   apiClient,
			Producer: kafka.NewProducer(),
		},
		Producer: kafka.NewProducer(),
	}

	ctx := context.Background()
	leagues := map[string]int{
		"LALIGA":        cfg.SofascoreApi.LeaguesID.LaLiga,
		"PREMIERLEAGUE": cfg.SofascoreApi.LeaguesID.PremierLeague,
	}

	for leagueName, leagueID := range leagues {
		for _, seasonID := range cfg.AllSeasons(leagueName) {
			if err := service.UpdateLeaguePlayersInfo(ctx, seasonID, leagueID); err != nil {
				t.Fatalf("Error: %v", err)
			}
		}
	}

	t.Log("Test completed successfully")
}
