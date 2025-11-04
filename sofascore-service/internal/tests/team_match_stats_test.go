package tests

import (
	"context"
	"log"
	"reflect"
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

	log.Println("Calling FPL API...")

	var laLigaSeasonIDs []int
	var premierLeagueSeasonIDs []int
	var leagueIDs []int

	ligaSeason := reflect.ValueOf(service.Config.SofascoreApi.LaLigaSeasonsIDs)
	for i := 0; i < ligaSeason.NumField(); i++ {
		laLigaSeasonIDs = append(laLigaSeasonIDs, int(ligaSeason.Field(i).Int()))
	}

	plSeason := reflect.ValueOf(service.Config.SofascoreApi.PremierLeagueSeasonIDs)
	for i := 0; i < plSeason.NumField(); i++ {
		premierLeagueSeasonIDs = append(premierLeagueSeasonIDs, int(plSeason.Field(i).Int()))
	}

	league := reflect.ValueOf(service.Config.SofascoreApi.LeaguesID)
	for i := 0; i < league.NumField(); i++ {
		leagueIDs = append(leagueIDs, int(league.Field(i).Int()))
	}

	start := time.Now()
	for _, leagueId := range leagueIDs {
		if leagueId == service.Config.SofascoreApi.LeaguesID.LaLiga {
			for _, seasonId := range laLigaSeasonIDs {
				for round := 1; round <= 38; round++ {
					log.Printf("Processing round %d", round)
					err := service.UpdateLeagueMatchStats(ctx, seasonId, leagueId, round)
					if err != nil {
						t.Fatalf("UpdateLeagueMatchStats failed: %v", err)
					}
				}
			}
		} else if leagueId == service.Config.SofascoreApi.LeaguesID.PremierLeague {
			for _, seasonId := range premierLeagueSeasonIDs {
				for round := 1; round <= 38; round++ {
					log.Printf("Processing round %d", round)
					err := service.UpdateLeagueMatchStats(ctx, seasonId, leagueId, round)
					if err != nil {
						t.Fatalf("UpdateLeagueMatchStats failed: %v", err)
					}
				}
			}
		}
	}

	elapsed := time.Since(start)

	log.Printf("Publishing completed in: %v", elapsed)
	t.Log("Sofascore API test completed successfully")
}
