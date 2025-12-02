package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
)

type LeagueService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
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

func (l *LeagueService) UpdateLeagueIDs(ctx context.Context) error {
	leagueIdsTopic := l.Config.KafkaConfig.TopicsName.SofascoreLeagueIDs

	leagueCountries, err := l.GetLeagueCountries(ctx)
	if err != nil {
		return fmt.Errorf("error getting country countries ids: %w", err)
	}

	for _, country := range leagueCountries.Categories {
		uniqueTournament := &sofascore_models.LeagueUniqueTournaments{}

		uniqueTournament, err = l.GetLeagueInfo(ctx, country.ID)
		if err != nil {
			return fmt.Errorf("error getting country info id: %d to update tournaments: %w", country.ID, err)
		}

		value, err := json.Marshal(uniqueTournament)
		if err != nil {
			return fmt.Errorf("error marshalling country uniqueTournament id: %d, %w", country.ID, err)
		}

		key := []byte(fmt.Sprintf("%d", country.ID))

		if err = l.Producer.Publish(ctx, leagueIdsTopic, key, value); err != nil {
			return fmt.Errorf("error publishing league ids for country id: %d: %w", country.ID, err)
		}
	}

	return nil
}
