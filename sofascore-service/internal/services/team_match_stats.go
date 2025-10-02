package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/sofascore-service/config"
	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
	"github.com/imadbelkat1/sofascore-service/internal/models"
)

type TeamMatchStatsService struct {
	Event    *EventsService
	Config   config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (s *TeamMatchStatsService) GetTeamMatchStats(ctx context.Context, matchId int) (*models.MatchStats, error) {
	var teamMatchStats models.MatchStats

	matchStats := s.Config.SofascoreApi.TeamEndpoints.TeamMatchStats
	endpoint := fmt.Sprintf(matchStats, matchId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &teamMatchStats); err != nil {
		return nil, err
	}

	return &teamMatchStats, nil
}

func (s *TeamMatchStatsService) UpdateLeagueMatchStats(ctx context.Context, seasonId int, leagueId int, round int) error {
	var matchRound = round
	teamMatchStats := s.Config.KafkaConfig.TopicsName.SofascoreTeamMatchStats

	roundMatches, err := s.Event.GetRoundMatches(ctx, leagueId, seasonId, round)
	if err != nil {
		return fmt.Errorf("getting round matches: %w", err)
	}

	jobs := make(chan models.Event, len(roundMatches.Events))

	var statsWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		statsWg.Add(1)
		go func() {
			defer statsWg.Done()
			for match := range jobs {
				matchStats, err := s.GetTeamMatchStats(ctx, match.ID)
				matchStats.MatchID = match.ID
				matchStats.Event = matchRound
				matchStats.HomeTeamID = match.HomeTeam.ID
				matchStats.AwayTeamID = match.AwayTeam.ID
				if err != nil {
					fmt.Printf("error fetching team match stats for match ID %d: %v\n", match.ID, err)
					continue
				}
				value, err := json.Marshal(matchStats)
				if err != nil {
					fmt.Printf("error marshalling team match stats for match ID %d: %v\n", match.ID, err)
					continue
				}
				key := []byte(fmt.Sprintf("%d", match.ID))

				if err := s.Producer.Publish(ctx, teamMatchStats, key, value); err != nil {
					fmt.Printf("error publishing team match stats for match ID %d: %v\n", match.ID, err)
					continue
				}
			}
		}()
	}

	for _, match := range roundMatches.Events {
		jobs <- match
	}
	close(jobs)

	statsWg.Wait()
	log.Printf("Successfully updated team match stats for league %d, season %d, round %d", leagueId, seasonId, round)

	return nil
}
