package services

import (
	"context"
	"fmt"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
	"golang.org/x/sync/errgroup"
)

type SeasonService struct {
	Config        *config.SofascoreConfig
	Client        *sofascore_api.SofascoreApiClient
	Producer      *kafka.Producer
	LeagueService *LeagueService
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

		uniqueTournament, err := s.LeagueService.GetLeagueInfo(ctx, int(country.Id))
		if err != nil {
			fmt.Printf("error getting country %d info: %s\n", country.Id, err)
			continue
		}

		for _, group := range uniqueTournament.Groups {
			for _, league := range group.UniqueTournaments {
				league := league

				g.Go(func() error {
					if err := s.UpdateLeaguesSeasons(ctx, int(league.Id)); err != nil {
						fmt.Printf("error updating league %d: %s\n", league.Id, err)
					}
					return nil
				})
			}
		}
	}

	return g.Wait()
}

func (s *SeasonService) UpdateLeaguesSeasons(ctx context.Context, leagueId int) error {
	leagueSeasonsTopic := s.Config.KafkaConfig.TopicsName.SofascoreLeagueSeasons.Name

	leagueSeasons, err := s.GetLeagueSeasons(ctx, leagueId)
	if err != nil {
		return fmt.Errorf("error getting league %d seasons: %w", leagueId, err)
	}

	// Extract the starting year from config (e.g., "25" from "25/26")
	configYear := s.Config.SofascoreApi.CurrentYear
	if len(configYear) >= 2 {
		configYear = configYear[:2] // "25/26" → "25"
	}

	for i := range leagueSeasons {
		season := &leagueSeasons[i]
		leagueSeasons[i].LeagueId = int32(leagueId)
		seasonYear := (*season).Year

		// Check multiple formats:
		// - Exact match: "25/26" == "25/26"
		// - Full year: "2025" contains "25"
		// - Short year starts with: "25/26" starts with "25"
		isCurrent := seasonYear == s.Config.SofascoreApi.CurrentYear ||
			(len(seasonYear) == 4 && seasonYear[2:] == configYear) || // "2025" ends with "25"
			(len(seasonYear) >= 2 && seasonYear[:2] == configYear) // "25/26" starts with "25"

		if isCurrent {
			(*season).IsCurrent = true
		} else {
			(*season).IsCurrent = false
		}
	}

	if err := s.publishLeagueSeasons(ctx, leagueId, leagueSeasons); err != nil {
		return fmt.Errorf("error publishing league %d seasons to topic %s: %w", leagueId, leagueSeasonsTopic, err)
	}

	return nil
}

func (s *SeasonService) GetLeagueSeasons(ctx context.Context, leagueId int) ([]*sofascore.Season, error) {
	var seasons []*sofascore.Season

	leagueSeason := s.Config.SofascoreApi.LeagueSeasonsIDs
	endpoint := fmt.Sprintf(leagueSeason, leagueId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, seasons); err != nil {
		return nil, fmt.Errorf("error getting league %d seasons: %w", leagueId, err)
	}

	return seasons, nil
}

func (s *SeasonService) publishLeagueSeasons(ctx context.Context, leagueId int, seasons []*sofascore.Season) error {
	leagueSeasonsTopic := s.Config.KafkaConfig.TopicsName.SofascoreLeagueSeasons.Name

	key := []byte(fmt.Sprintf("%d", leagueId))

	if err := s.Producer.PublishWithProcess(ctx, seasons, leagueSeasonsTopic, key); err != nil {
		return fmt.Errorf("error publishing league %d seasons: %w", leagueId, err)
	}

	return nil
}
