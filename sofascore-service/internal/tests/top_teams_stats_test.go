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

func TestTopTeamsStatsService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()
	service := &services.TopTeamsStatsService{
		Config:   *cfg,
		Client:   sofascore_api.NewSofascoreApiClient(cfg),
		Producer: kafka.NewProducer(),
	}

	ctx := context.Background()
	seasonID := cfg.MustGetSeasonID("PREMIERLEAGUE", "2425")
	leagueID := cfg.SofascoreApi.LeaguesID.PremierLeague

	log.Printf("Testing season %d, league %d", seasonID, leagueID)

	start := time.Now()
	if err := service.UpdateLeagueTopTeamsStats(ctx, seasonID, leagueID); err != nil {
		t.Fatalf("Error: %v", err)
	}

	log.Printf("Completed in %v", time.Since(start))
	t.Log("Test completed successfully")
}
