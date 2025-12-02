package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	"golang.org/x/sync/errgroup"
)

type SeasonService struct {
	Config        *config.SofascoreConfig
	Client        *sofascore_api.SofascoreApiClient
	Producer      *kafka.Producer
	LeagueService *LeagueService
}

func (s *SeasonService) GetLeagueSeasons(ctx context.Context, leagueId int) (*sofascore_models.Seasons, error) {
	seasons := &sofascore_models.Seasons{}

	leagueSeason := s.Config.SofascoreApi.LeagueSeasonsIDs
	endpoint := fmt.Sprintf(leagueSeason, leagueId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, seasons); err != nil {
		return nil, fmt.Errorf("error getting league %d seasons: %w", leagueId, err)
	}

	return seasons, nil
}

func (s *SeasonService) UpdateLeaguesSeasons(ctx context.Context, leagueId int) error {
	leagueSeasonsTopic := s.Config.KafkaConfig.TopicsName.SofascoreLeagueSeasons

	leagueSeasons, err := s.GetLeagueSeasons(ctx, leagueId)
	if err != nil {
		return fmt.Errorf("error getting league %d seasons: %w", leagueId, err)
	}

	leagueSeasons.LeagueID = leagueId

	// Extract the starting year from config (e.g., "25" from "25/26")
	configYear := s.Config.SofascoreApi.CurrentYear
	if len(configYear) >= 2 {
		configYear = configYear[:2] // "25/26" → "25"
	}

	for i := range leagueSeasons.Seasons {
		season := &leagueSeasons.Seasons[i]
		seasonYear := season.Year

		// Check multiple formats:
		// - Exact match: "25/26" == "25/26"
		// - Full year: "2025" contains "25"
		// - Short year starts with: "25/26" starts with "25"
		isCurrent := seasonYear == s.Config.SofascoreApi.CurrentYear ||
			(len(seasonYear) == 4 && seasonYear[2:] == configYear) || // "2025" ends with "25"
			(len(seasonYear) >= 2 && seasonYear[:2] == configYear) // "25/26" starts with "25"

		if isCurrent {
			season.IsCurrent = true
		} else {
			season.IsCurrent = false
		}
	}

	value, err := json.Marshal(leagueSeasons)
	if err != nil {
		return fmt.Errorf("error marshaling league %d seasons: %w", leagueId, err)
	}

	key := []byte(fmt.Sprintf("%d", leagueId))

	if err = s.Producer.Publish(ctx, leagueSeasonsTopic, key, value); err != nil {
		return fmt.Errorf("error publishing league %d seasons: %w", leagueId, err)
	}

	return nil
}

func (s *SeasonService) UpdateAllLeaguesSeasons(ctx context.Context) error {
	leagueCountries, err := s.LeagueService.GetLeagueCountries(ctx)
	if err != nil {
		return fmt.Errorf("error getting country countries ids: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	for _, country := range leagueCountries.Categories {
		country := country

		uniqueTournament, err := s.LeagueService.GetLeagueInfo(ctx, country.ID)
		if err != nil {
			fmt.Printf("error getting country %d info: %s\n", country.ID, err)
			continue
		}

		for _, group := range uniqueTournament.Groups {
			for _, league := range group.UniqueTournament {
				league := league

				g.Go(func() error {
					if err := s.UpdateLeaguesSeasons(ctx, league.ID); err != nil {
						fmt.Printf("error updating league %d: %s\n", league.ID, err)
					}
					return nil
				})
			}
		}
	}

	return g.Wait()
}
