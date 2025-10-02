package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/sofascore-service/config"
	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
	teamMatchStatservice "github.com/imadbelkat1/sofascore-service/internal/services"
)

func TestTeamMatchStatsService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	service := &teamMatchStatservice.TeamMatchStatsService{
		Event: &teamMatchStatservice.EventsService{
			Config: config.LoadConfig(),
			Client: sofascore_api.NewSofascoreApiClient(config.LoadConfig()),
		},
		Config:   *config.LoadConfig(),
		Client:   sofascore_api.NewSofascoreApiClient(config.LoadConfig()),
		Producer: kafka.NewProducer(),
	}

	time.Sleep(200 * time.Millisecond)

	log.Println("Calling FPL API...")

	seasonId := service.Config.SofascoreApi.SeasonsIDs.Laliga2425
	leagueId := service.Config.SofascoreApi.LeaguesIDs.LaLiga
	log.Println(seasonId)
	log.Println(leagueId)
	//log.Println(round)

	start := time.Now()
	for round := 1; round <= 38; round++ {
		log.Printf("Processing round %d", round)
		err := service.UpdateLeagueMatchStats(ctx, seasonId, leagueId, round)
		if err != nil {
			t.Fatalf("UpdateLeagueMatchStats failed: %v", err)
		}
	}
	elapsed := time.Since(start)

	log.Printf("Publishing completed in: %v", elapsed)
	t.Log("Sofascore API test completed successfully")
}
