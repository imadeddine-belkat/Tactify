package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/shared/sofascore_models"
	"github.com/imadbelkat1/sofascore-service/config"
	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
	"golang.org/x/sync/errgroup"
)

type LeagueStandingService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (l *LeagueStandingService) UpdateLeagueStanding(ctx context.Context, seasonId int, leagueId int) error {
	standing, err := l.GetLeagueStanding(ctx, seasonId, leagueId)
	if err != nil {
		return fmt.Errorf("getting league standing: %w", err)
	}

	if err := l.PublishLeagueStanding(ctx, seasonId, leagueId, standing); err != nil {
		return fmt.Errorf("publishing league standing: %w", err)
	}

	return nil
}

func (t *LeagueStandingService) GetLeagueStanding(ctx context.Context, seasonId int, leagueId int) (*sofascore_models.Standings, error) {
	var standing *sofascore_models.Standings

	leagueStandingEndpoint := t.Config.SofascoreApi.LeagueEndpoints.LeagueSeasonStandings // /unique-tournament/%d/season/%d/standings/total
	endpoint := fmt.Sprintf(leagueStandingEndpoint, leagueId, seasonId)

	if err := t.Client.GetAndUnmarshal(ctx, endpoint, &standing); err != nil {
		return nil, fmt.Errorf("fetching league standing data: %w", err)
	}

	return standing, nil
}

func (t *LeagueStandingService) PublishLeagueStanding(ctx context.Context, seasonId int, leagueId int, standing *sofascore_models.Standings) error {
	leagueStandingTopic := t.Config.KafkaConfig.TopicsName.SofascoreLeagueStandings

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	for _, s := range standing.Standings {
		for _, row := range s.Rows {
			r := row
			g.Go(func() error {
				leagueStanding := &sofascore_models.StandingMessage{
					SeasonID: seasonId,
					LeagueID: leagueId,
					Row:      r,
				}
				value, err := json.Marshal(leagueStanding)
				if err != nil {
					return fmt.Errorf("marshaling league standing message: %w", err)
				}

				key := []byte(fmt.Sprintf("%d-%d", leagueId, r.Team.ID)) // Unique key per team in league

				return t.Producer.Publish(ctx, leagueStandingTopic, key, value)

			})
		}
	}

	return g.Wait()
}
