package tests

import (
	"context"
	"testing"

	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/sofascore-service/config"
	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
	leagueStandingService "github.com/imadbelkat1/sofascore-service/internal/services"
)

func TestLeagueStandingService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	service := &leagueStandingService.LeagueStandingService{
		Config:   config.LoadConfig(),
		Client:   sofascore_api.NewSofascoreApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}
	seasonId := service.Config.SofascoreApi.SeasonsIDs.PremierLeague2526
	leagueId := service.Config.SofascoreApi.LeaguesIDs.PremierLeague

	t.Log("Calling Sofascore API...")
	err := service.UpdateLeagueStanding(ctx, seasonId, leagueId)
	if err != nil {
		t.Fatalf("GetLeagueStanding Update failed: %v", err)
	}

	t.Logf("Fetched standings for league %d, season %d", leagueId, seasonId)
	t.Log("Sofascore API Publish test completed successfully")

}
