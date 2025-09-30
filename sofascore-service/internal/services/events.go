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

func (e *EventsService) GetRoundMatches(leagueId int, leagueName string, season int, round int) (*models.Events, error) {
	ctx := context.Background()
	var Event models.Events

	endpoint := fmt.Sprintf(e.Client.Config.SofascoreApi.LeagueEndpoints.LeagueRoundMatches, leagueId, season, round)

	if err := e.Client.GetAndUnmarshal(ctx, endpoint, &Event); err != nil {
		return nil, fmt.Errorf("fetching events data: %w", err)
	}

	return &Event, nil
}
