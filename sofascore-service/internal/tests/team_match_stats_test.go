package tests

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"github.com/imadeddine-belkat/sofascore-service/internal/services"
)

func TestTeamMatchStatsService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping API test")
	}

	cfg := config.LoadConfig()
	service := &services.TeamMatchStatsService{
		Event: &services.EventsService{
			Config: cfg,
			Client: sofascore_api.NewSofascoreApiClient(cfg),
		},
		Config:   *cfg,
		Client:   sofascore_api.NewSofascoreApiClient(cfg),
		Producer: kafka.NewProducer(),
	}

	type workItem struct {
		seasonID int
		leagueID int
		round    int
	}

	leagues := map[string]int{
		"LALIGA":        cfg.SofascoreApi.LeaguesID.LaLiga,
		"PREMIERLEAGUE": cfg.SofascoreApi.LeaguesID.PremierLeague,
	}

	// Build work queue
	var work []workItem
	for leagueName, leagueID := range leagues {
		for _, seasonID := range cfg.AllSeasons(leagueName) {
			for round := 1; round <= 38; round++ {
				work = append(work, workItem{seasonID, leagueID, round})
			}
		}
	}

	// Process with workers
	start := time.Now()
	const numWorkers = 10

	workChan := make(chan workItem, len(work))
	errChan := make(chan error, len(work))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for item := range workChan {
				log.Printf("[Worker %d] Season %d, Round %d", workerID, item.seasonID, item.round)
				if err := service.UpdateLeagueMatchStats(context.Background(), item.seasonID, item.leagueID, item.round); err != nil {
					errChan <- fmt.Errorf("season %d, round %d: %w", item.seasonID, item.round, err)
				}
			}
		}(i)
	}

	for _, item := range work {
		workChan <- item
	}
	close(workChan)
	wg.Wait()
	close(errChan)

	// Check errors
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	log.Printf("Completed in %v. Processed %d items", time.Since(start), len(work))

	if len(errs) > 0 {
		t.Fatalf("Failed with %d errors: %v", len(errs), errs[0])
	}

	t.Log("Test completed successfully")
}
