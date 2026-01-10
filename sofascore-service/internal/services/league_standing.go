package services

import (
	"context"
	"fmt"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"golang.org/x/sync/errgroup"
)

type LeagueStandingService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (t *LeagueStandingService) UpdateLeagueStanding(ctx context.Context, seasonId, leagueId int) error {
	standing, err := t.GetLeagueStanding(ctx, seasonId, leagueId)
	if err != nil {
		return fmt.Errorf("getting league standing: %w", err)
	}

	if err := t.publishLeagueStanding(ctx, seasonId, leagueId, standing); err != nil {
		return fmt.Errorf("publishing league standing: %w", err)
	}

	return nil
}

func (t *LeagueStandingService) GetLeagueStanding(ctx context.Context, seasonId, leagueId int) (*sofascore_models.Standings, error) {
	standing := &sofascore_models.Standings{}

	leagueStandingEndpoint := t.Config.SofascoreApi.LeagueEndpoints.LeagueSeasonStandings // /unique-tournament/%d/season/%d/standings/total
	endpoint := fmt.Sprintf(leagueStandingEndpoint, leagueId, seasonId)

	if err := t.Client.GetAndUnmarshal(ctx, endpoint, standing); err != nil {
		return nil, fmt.Errorf("fetching league standing data: %w", err)
	}

	return standing, nil
}

func (t *LeagueStandingService) publishLeagueStanding(ctx context.Context, seasonId int, leagueId int, standing *sofascore_models.Standings) error {
	leagueStandingTopic := t.Config.KafkaConfig.TopicsName.SofascoreLeagueStandings.Name

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	for _, s := range standing.Standings {
		for _, row := range s.Rows {
			row := row
			g.Go(func() error {
				leagueStanding := &sofascore_models.StandingMessage{
					SeasonID: seasonId,
					LeagueID: leagueId,
					Row:      row,
				}
				key := []byte(fmt.Sprintf("%d-%d", leagueId, row.Team.ID))

				return t.Producer.PublishWithProcess(ctx, leagueStanding, leagueStandingTopic, key)

			})
		}
	}

	return g.Wait()
}
