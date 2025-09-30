package services

import (
	"context"
	"fmt"

	sofascore_api "github.com/imadbelkat1/sofascore-service/internal/api"
	"github.com/imadbelkat1/sofascore-service/internal/models"
)

type EventsService struct {
	Client *sofascore_api.SofascoreApiClient
}

func (e *EventsService) GetRoundMatches(ctx context.Context, leagueId int, leagueName string, season int, round int) (*models.Events, error) {
	var Event models.Events
	leagueRoundMatches := e.Client.Config.SofascoreApi.LeagueEndpoints.LeagueRoundMatches

	endpoint := fmt.Sprintf(leagueRoundMatches, leagueId, season, round)

	if err := e.Client.GetAndUnmarshal(ctx, endpoint, &Event); err != nil {
		return nil, fmt.Errorf("fetching events data: %w", err)
	}

	return &Event, nil
}
