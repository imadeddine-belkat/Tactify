package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
)

// statJob represents a single unit of work for the worker pool.
type statJob struct {
	label string
	item  *sofascore.TopTeamStatItem
}

type TopTeamsStatsService struct {
	Config   config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

// UpdateLeagueTopTeamsStats orchestrates the processing of all stat categories using a Worker Pool.
func (s *TopTeamsStatsService) UpdateLeagueTopTeamsStats(ctx context.Context, seasonId int, leagueId int) error {
	msg, err := s.GetTopTeamsStats(ctx, leagueId, seasonId)
	if err != nil {
		return err
	}
	if msg.TopTeams == nil {
		return fmt.Errorf("top teams data missing for league %d", leagueId)
	}

	topic := s.Config.KafkaConfig.TopicsName.SofascoreTopTeamsStats.Name

	// 1. Initialize the channel and the worker pool
	jobChan := make(chan statJob, 100)
	var wg sync.WaitGroup
	numWorkers := 5

	// 2. Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg, jobChan, topic, seasonId, leagueId)
	}

	// 3. Feed the channel (Producer)
	s.enqueueJobs(msg.TopTeams, jobChan)

	// 4. Close channel and wait for completion
	close(jobChan)
	wg.Wait()

	log.Printf("[INFO] Success: League %d (Season %d) processed via Worker Pool", leagueId, seasonId)
	return nil
}

// worker consumes jobs from the channel and publishes to Kafka by accessing fields directly.
func (s *TopTeamsStatsService) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan statJob, topic string, seasonId, leagueId int) {
	defer wg.Done()
	for job := range jobs {
		stat := job.item
		stat.SeasonId = int32(seasonId)
		stat.LeagueId = int32(leagueId)

		payload, err := json.Marshal(stat)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal %s: %v", job.label, err)
			continue
		}

		// Direct access to team ID without helper methods
		teamID := 0
		if stat.Team != nil {
			teamID = int(stat.Team.Id)
		}

		key := []byte(fmt.Sprintf("%s-%d-%d-%d", job.label, seasonId, leagueId, teamID))
		if err := s.Producer.Publish(ctx, topic, key, payload); err != nil {
			log.Printf("[ERROR] Kafka publish failed for %s: %v", job.label, err)
		}
	}
}

// enqueueJobs extracts all stats from the TopTeams message and sends them to the job channel.
func (s *TopTeamsStatsService) enqueueJobs(tt *sofascore.TopTeams, jobChan chan<- statJob) {
	categories := map[string][]*sofascore.TopTeamStatItem{
		"averageRating":         tt.AvgRating,
		"goalsScored":           tt.GoalsScored,
		"goalsConceded":         tt.GoalsConceded,
		"bigChances":            tt.BigChances,
		"bigChancesMissed":      tt.BigChancesMissed,
		"hitWoodWork":           tt.HitWoodWork,
		"yellowCards":           tt.YellowCards,
		"redCards":              tt.RedCards,
		"averageBallPossession": tt.AverageBallPossession,
		"accuratePasses":        tt.AccuratePasses,
		"accurateLongBalls":     tt.AccurateLongBalls,
		"accurateCrosses":       tt.AccurateCrosses,
		"shots":                 tt.Shots,
		"shotsOnTarget":         tt.ShotsOnTarget,
		"successfulDribbles":    tt.SuccessfulDribbles,
		"tackles":               tt.Tackles,
		"interceptions":         tt.Interceptions,
		"clearances":            tt.Clearances,
		"corners":               tt.Corners,
		"fouls":                 tt.Fouls,
		"penaltyGoals":          tt.PenaltyGoals,
		"penaltyGoalsConceded":  tt.PenaltyGoalsConceded,
		"cleanSheets":           tt.CleanSheets,
	}

	for label, items := range categories {
		for _, item := range items {
			if item != nil {
				jobChan <- statJob{label: label, item: item}
			}
		}
	}
}

// GetTopTeamsStats fetches the raw data from the SofaScore API.
func (s *TopTeamsStatsService) GetTopTeamsStats(ctx context.Context, leagueId int, seasonId int) (*sofascore.TopTeamsMessage, error) {
	topTeams := &sofascore.TopTeamsMessage{}
	endpoint := fmt.Sprintf(s.Config.SofascoreApi.TeamEndpoints.TopTeamsStats, leagueId, seasonId)
	if err := s.Client.GetAndUnmarshal(ctx, endpoint, topTeams); err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	return topTeams, nil
}
