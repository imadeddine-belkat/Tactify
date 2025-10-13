package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/shared/sofascore_models"
	"github.com/imadbelkat1/sofascore-service/config"
	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
)

type TopTeamsStatsService struct {
	Config   config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

// TeamStatter interface
type TeamStatter interface {
	GetTeamID() int
	SetSeasonID(seasonId int)
	SetLeagueID(leagueId int)
}

func publishStats[T TeamStatter](
	ctx context.Context,
	producer *kafka.Producer,
	topic string,
	seasonId int,
	leagueId int,
	stats *[]T,
	statType string,
) error {
	for _, stat := range *stats {
		stat.SetSeasonID(seasonId)
		stat.SetLeagueID(leagueId)

		statsBytes, err := json.Marshal(stat)
		if err != nil {
			log.Printf("Error marshalling %s stats: %v", statType, err)
			continue
		}

		teamID := stat.GetTeamID()
		key := []byte(fmt.Sprintf("%s-%d-%d-%d", statType, seasonId, leagueId, teamID))

		if err := producer.Publish(ctx, topic, key, statsBytes); err != nil {
			log.Printf("Error producing %s stats to Kafka: %v", statType, err)
		} else {
			log.Printf("Produced %s stats for league %d, season %d, for team ID %d to topic %s", statType, leagueId, seasonId, teamID, topic)
		}
	}

	return nil
}

// publishStatsConcurrent wraps publishStats for concurrent execution
func publishStatsConcurrent[T TeamStatter](
	ctx context.Context,
	wg *sync.WaitGroup,
	producer *kafka.Producer,
	topic string,
	seasonId int,
	leagueId int,
	stats *[]T,
	statType string,
) {
	defer wg.Done()
	if err := publishStats(ctx, producer, topic, seasonId, leagueId, stats, statType); err != nil {
		log.Printf("Error in concurrent publish for %s: %v", statType, err)
	}
}

func (s *TopTeamsStatsService) GetTopTeamsStats(ctx context.Context, leagueId int, seasonId int) (*sofascore_models.TopTeamsMessage, error) {
	var topTeams *sofascore_models.TopTeamsMessage

	topTeamsStats := s.Config.SofascoreApi.TeamEndpoints.TopTeamsStats
	endpoint := fmt.Sprintf(topTeamsStats, leagueId, seasonId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &topTeams); err != nil {
		return nil, err
	}

	return topTeams, nil
}

// UpdateLeagueTopTeamsStats Concurrent
func (s *TopTeamsStatsService) UpdateLeagueTopTeamsStats(ctx context.Context, seasonId int, leagueId int) error {

	topTeamsStats, err := s.GetTopTeamsStats(ctx, leagueId, seasonId)
	if err != nil {
		return fmt.Errorf("getting top teams stats: %w", err)
	}

	topic := s.Config.KafkaConfig.TopicsName.SofascoreTopTeamsStats
	var wg sync.WaitGroup

	wg.Add(23)
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.AverageRating, "averageRating")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.GoalsScored, "goalsScored")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.GoalsConceded, "goalsConceded")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.BigChances, "bigChances")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.BigChancesMissed, "bigChancesMissed")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.HitWoodWork, "hitWoodWork")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.YellowCards, "yellowCards")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.RedCards, "redCards")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.AverageBallPossession, "averageBallPossession")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.AccuratePasses, "accuratePasses")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.AccurateLongBalls, "accurateLongBalls")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.AccurateCrosses, "accurateCrosses")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.Shots, "shots")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.ShotsOnTarget, "shotsOnTarget")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.SuccessfulDribbles, "successfulDribbles")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.Tackles, "tackles")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.Interceptions, "interceptions")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.Clearances, "clearances")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.Corners, "corners")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.Fouls, "fouls")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.PenaltyGoals, "penaltyGoals")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.PenaltyGoalsConceded, "penaltyGoalsConceded")
	go publishStatsConcurrent(ctx, &wg, s.Producer, topic, seasonId, leagueId, &topTeamsStats.TopTeams.CleanSheets, "cleanSheets")

	wg.Wait()

	log.Printf("Successfully processed all top teams stats concurrently for league %d, season %d", leagueId, seasonId)
	return nil
}
