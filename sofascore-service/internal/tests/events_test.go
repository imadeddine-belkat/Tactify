package tests

import (
	"context"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	eventService "github.com/imadeddine-belkat/sofascore-service/internal/services"
)

func TestEventsService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()
	service := &eventService.EventsService{
		Config:   cfg,
		Client:   sofascore_api.NewSofascoreApiClient(cfg),
		Producer: kafka.NewProducer(),
	}

	leagues := map[string]int{
		"LALIGA":        cfg.SofascoreApi.LeaguesID.LaLiga,
		"PREMIERLEAGUE": cfg.SofascoreApi.LeaguesID.PremierLeague,
	}

	ctx := context.Background()
	g, gCtx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(5) // Max 5 concurrent requests

	for leagueName, leagueID := range leagues {
		for _, seasonID := range cfg.AllSeasons(leagueName) {
			for round := 1; round <= 38; round++ {
				// Capture variables
				leagueName, leagueID, seasonID, round := leagueName, leagueID, seasonID, round

				g.Go(func() error {
					if err := sem.Acquire(gCtx, 1); err != nil {
						return err
					}
					defer sem.Release(1)

					t.Logf("Fetching %s SeasonID: %d, Round: %d", leagueName, seasonID, round)

					// Retry up to 3 times
					var err error
					for attempt := 1; attempt <= 3; attempt++ {
						err = service.UpdateRoundMatches(gCtx, seasonID, leagueID, round)
						if err == nil {
							return nil
						}
						if attempt < 3 {
							time.Sleep(2 * time.Second)
						}
					}

					t.Errorf("Failed %s SeasonID %d, Round %d: %v", leagueName, seasonID, round, err)
					return err
				})
			}
		}
	}

	if err := g.Wait(); err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	t.Log("Test completed successfully")
}
