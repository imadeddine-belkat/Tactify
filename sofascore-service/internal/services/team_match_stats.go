package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
)

type TeamMatchStatsService struct {
	Event    *EventsService
	Config   config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (s *TeamMatchStatsService) GetTeamMatchStats(ctx context.Context, matchId int) (*sofascore_models.MatchStats, error) {
	teamMatchStats := &sofascore_models.MatchStats{}

	matchStats := s.Config.SofascoreApi.TeamEndpoints.TeamMatchStats
	endpoint := fmt.Sprintf(matchStats, matchId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, teamMatchStats); err != nil {
		return nil, err
	}

	return teamMatchStats, nil
}

func (s *TeamMatchStatsService) UpdateLeagueMatchStats(ctx context.Context, seasonId int, leagueId int, round int) error {
	teamMatchStatsTopic := s.Config.KafkaConfig.TopicsName.SofascoreTeamMatchStats

	log.Printf("Starting UpdateLeagueMatchStats for league=%d, season=%d, round=%d", leagueId, seasonId, round)

	roundMatches, err := s.Event.GetRoundMatches(ctx, seasonId, leagueId, round)
	if err != nil {
		return fmt.Errorf("getting round matches: %w", err)
	}

	matchCount := len(roundMatches.Events)
	log.Printf("Found %d matches to process", matchCount)

	if matchCount == 0 {
		log.Printf("No matches found for league %d, season %d, round %d", leagueId, seasonId, round)
		return nil
	}

	jobs := make(chan sofascore_models.Event, matchCount)

	// Use optimal worker count: min(configured workers, number of matches)
	workerCount := s.Config.PublishWorkerCount
	if matchCount < workerCount {
		workerCount = matchCount
	}

	log.Printf("Starting %d workers for %d matches", workerCount, matchCount)

	var statsWg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		workerID := i
		statsWg.Add(1)
		go func() {
			defer statsWg.Done()

			processedMatches := 0
			for match := range jobs {
				processedMatches++

				select {
				case <-ctx.Done():
					log.Printf("Worker %d cancelled: %v", workerID, ctx.Err())
					return
				default:
				}

				matchStats, err := s.GetTeamMatchStats(ctx, match.ID)
				if err != nil {
					log.Printf("Worker %d: error fetching stats for match %d: %v", workerID, match.ID, err)
					continue
				}

				// Skip matches without statistics
				if len(matchStats.MatchPeriods) == 0 {
					log.Printf("Worker %d: No statistics available for match %d (not played yet)", workerID, match.ID)
					continue
				}

				publishCount := 0
				for _, period := range matchStats.MatchPeriods {
					for _, group := range period.Groups {
						for _, stat := range group.StatsItems {
							// Check context before each publish
							select {
							case <-ctx.Done():
								log.Printf("Worker %d cancelled during publish", workerID)
								return
							default:
							}

							msg := &sofascore_models.MatchStatsMessage{
								SeasonId:   seasonId,
								LeagueId:   leagueId,
								MatchID:    match.ID,
								Event:      round,
								HomeTeamID: match.HomeTeam.ID,
								AwayTeamID: match.AwayTeam.ID,
								GroupName:  group.GroupName,
								MatchStatistics: sofascore_models.StatsMessage{
									Period:         period.Period,
									Key:            stat.Key,
									Name:           stat.Name,
									HomeValue:      stat.HomeValue,
									AwayValue:      stat.AwayValue,
									HomeTotal:      stat.HomeTotal,
									AwayTotal:      stat.AwayTotal,
									CompareCode:    stat.CompareCode,
									StatisticsType: stat.StatisticsType,
									RenderType:     stat.RenderType,
								},
							}

							value, err := json.Marshal(msg)
							if err != nil {
								log.Printf("Worker %d: error marshalling stats for match %d: %v", workerID, match.ID, err)
								continue
							}

							key := []byte(fmt.Sprintf("%d-%s-%s", match.ID, period.Period, stat.Name))

							if err := s.Producer.Publish(ctx, teamMatchStatsTopic, key, value); err != nil {
								log.Printf("Worker %d: error publishing stat for match %d: %v", workerID, match.ID, err)
								continue
							}

							publishCount++
						}
					}
				}

				log.Printf("Worker %d: Published %d stats for match %d", workerID, publishCount, match.ID)
			}

			if processedMatches > 0 {
				log.Printf("Worker %d finished, processed %d matches", workerID, processedMatches)
			}
		}()
	}

	// Enqueue all matches
	for _, match := range roundMatches.Events {
		jobs <- match
	}
	close(jobs)

	statsWg.Wait()
	log.Printf("Successfully updated team match stats for league %d, season %d, round %d", leagueId, seasonId, round)

	return nil
}
