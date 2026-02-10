package services

import (
	"context"
	"fmt"
	"log"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"github.com/imadeddine-belkat/tactify-protos/sofascore_models"
)

type TeamOverallStatsService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
	Standing *LeagueStandingService
}

func (o *TeamOverallStatsService) GetTeamOverallStats(ctx context.Context, teamId, leagueId, seasonId int) (*sofascore_models.TeamOverallStatsMessage, error) {
	teamOverallStats := &sofascore_models.TeamOverallStatsMessage{}

	teamStats := o.Config.SofascoreApi.TeamEndpoints.TeamOverallStats
	endpoint := fmt.Sprintf(teamStats, teamId, leagueId, seasonId) //team/%d/unique-tournament/%d/season/%d/statistics/overall

	if err := o.Client.GetAndUnmarshal(ctx, endpoint, teamOverallStats); err != nil {
		return nil, err
	}

	return teamOverallStats, nil
}

func (o *TeamOverallStatsService) UpdateTeamOverallStats(ctx context.Context, teamId, leagueId, seasonId int) error {
	teamOverallStatsTopic := o.Config.KafkaConfig.TopicsName.SofascoreTeamOverallStats.Name

	teamOverallStats, err := o.GetTeamOverallStats(ctx, teamId, leagueId, seasonId)
	if err != nil {
		return err
	}

	teamOverallStats.TeamID = teamId
	teamOverallStats.LeagueID = leagueId
	teamOverallStats.SeasonID = seasonId

	key := []byte(fmt.Sprintf("%d_%d_%d", teamId, leagueId, seasonId))

	if err = o.Producer.PublishWithProcess(ctx, teamOverallStats, teamOverallStatsTopic, key); err != nil {
		return fmt.Errorf("failed to publish teamOverallStats to Kafka for team: %d, season: %d, league: %d: %w", teamId, seasonId, leagueId, err)
	}

	return nil
}

func (o *TeamOverallStatsService) UpdateAllTeamsOverallStats(ctx context.Context, leagueId, seasonId int) error {
	var teamsIds []int

	//get team ids from league standing
	standing, err := o.Standing.GetLeagueStanding(ctx, seasonId, leagueId)
	if err != nil {
		return fmt.Errorf("failed to get league standing for LeagueID %d: %w", seasonId, err)
	}

	for _, st := range standing.Standings {
		for _, row := range st.Rows {
			teamsIds = append(teamsIds, row.Team.ID)
		}
	}

	for _, teamId := range teamsIds {
		err = o.UpdateTeamOverallStats(ctx, teamId, leagueId, seasonId)
		if err != nil {
			log.Printf("failed to update teamOverallStats for TeamID %d, LeagueID %d, SeasonID %d: %w", teamId, leagueId, seasonId, err)
		}

	}

	return nil
}
