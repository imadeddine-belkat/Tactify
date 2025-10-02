package services

import (
	"context"
	"fmt"
	"log"

	"github.com/imadbelkat1/sofascore-service/config"
	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
	"github.com/imadbelkat1/sofascore-service/internal/models"
)

type EventsService struct {
	Config *config.SofascoreConfig
	Client *sofascore_api.SofascoreApiClient
}

func (e *EventsService) GetRoundMatches(ctx context.Context, leagueId int, season int, round int) (*models.Events, error) {
	var Event models.Events
	leagueRoundMatches := e.Config.SofascoreApi.LeagueEndpoints.LeagueRoundMatches //unique-tournament/%d/season/%d/events/round/%d
	log.Printf("leagueRoundMatches raw: %q", leagueRoundMatches)

	endpoint := fmt.Sprintf(leagueRoundMatches, leagueId, season, round)
	log.Println(endpoint)

	if err := e.Client.GetAndUnmarshal(ctx, endpoint, &Event); err != nil {
		return nil, fmt.Errorf("fetching events data: %w", err)
	}

	return &Event, nil
}
