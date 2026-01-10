package tests

import (
	"context"
	"log"
	"testing"
	"time"

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
	service := &services.PlayersService{
		Config: *cfg,
		Client: sofascore_api.NewSofascoreApiClient(cfg),
		Standing: &services.LeagueStandingService{
			Config:   cfg,
			Client:   sofascore_api.NewSofascoreApiClient(cfg),
			Producer: kafka.NewProducer(),
		},
		Producer: kafka.NewProducer(),
	}

	ctx := context.Background()
	seasonID := cfg.MustGetSeasonID("PREMIERLEAGUE", "2526")
	leagueID := cfg.SofascoreApi.LeaguesID.PremierLeague

	log.Printf("Testing season %d, league %d", seasonID, leagueID)

	start := time.Now()
	if err := service.UpdateLeaguePlayersInfo(ctx, seasonID, leagueID); err != nil {
		t.Fatalf("Error: %v", err)
	}

	log.Printf("Completed in %v", time.Since(start))
	t.Log("Test completed successfully")
}
