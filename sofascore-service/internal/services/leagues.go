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

type LeagueService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (l *LeagueService) UpdateLeagueIDs(ctx context.Context) error {
	leagueCountries, err := l.GetLeagueCountries(ctx)
	if err != nil {
		return fmt.Errorf("error getting country countries ids: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	for _, country := range leagueCountries.Categories {
		country := country
		g.Go(func() error {
			uniqueTournament := &sofascore_models.LeagueUniqueTournaments{}

			uniqueTournament, err = l.GetLeagueInfo(ctx, country.ID)
			if err != nil {
				return fmt.Errorf("error getting country info id: %d to update tournaments: %w", country.ID, err)
			}
			return l.publishLeagueInfo(ctx, country.ID, uniqueTournament)
		})
	}

	return g.Wait()
}

func (l *LeagueService) GetLeagueCountries(ctx context.Context) (*sofascore_models.LeagueCategories, error) {
	leagueCategories := &sofascore_models.LeagueCategories{}

	endpoint := l.Config.SofascoreApi.LeagueCountriesIDs
	if err := l.Client.GetAndUnmarshal(ctx, endpoint, leagueCategories); err != nil {
		return nil, fmt.Errorf("error getting league countries ids: %w", err)
	}

	return leagueCategories, nil
}

func (l *LeagueService) GetLeagueInfo(ctx context.Context, countryId int) (*sofascore_models.LeagueUniqueTournaments, error) {
	league := &sofascore_models.LeagueUniqueTournaments{}

	endpoint := l.Config.SofascoreApi.LeagueCountryLeagueIDs
	countryLeagues := fmt.Sprintf(endpoint, countryId)

	if err := l.Client.GetAndUnmarshal(ctx, countryLeagues, league); err != nil {
		return nil, fmt.Errorf("error getting league league ids: %w", err)
	}

	return league, nil
}

func (l *LeagueService) publishLeagueInfo(ctx context.Context, countryId int, uniqueTournament *sofascore_models.LeagueUniqueTournaments) error {
	leagueIdsTopic := l.Config.KafkaConfig.TopicsName.SofascoreLeagueIDs.Name

	key := []byte(fmt.Sprintf("%d", countryId))

	if err := l.Producer.PublishWithProcess(ctx, uniqueTournament, leagueIdsTopic, key); err != nil {
		return fmt.Errorf("error publishing league ids for country id: %d: %w", countryId, err)
	}

	return nil
}
