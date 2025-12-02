package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kafka "github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
)

type EventsService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (e *EventsService) GetRoundMatches(ctx context.Context, seasonId, leagueId, round int) (*sofascore_models.Events, error) {
	event := &sofascore_models.Events{}
	leagueRoundMatches := e.Config.SofascoreApi.LeagueEndpoints.LeagueRoundMatches //unique-tournament/%d/seasonId/%d/events/round/%d
	log.Printf("leagueRoundMatches raw: %q", leagueRoundMatches)

	endpoint := fmt.Sprintf(leagueRoundMatches, leagueId, seasonId, round)
	log.Println(endpoint)

	if err := e.Client.GetAndUnmarshal(ctx, endpoint, event); err != nil {
		return nil, fmt.Errorf("fetching events data: %w", err)
	}

	return event, nil
}

func (e *EventsService) UpdateRoundMatches(ctx context.Context, seasonId, leagueId, round int) error {
	roundMatches, err := e.GetRoundMatches(ctx, seasonId, leagueId, round)
	if err != nil {
		return fmt.Errorf("error fetching events data: %w", err)
	}

	for _, event := range roundMatches.Events {
		err := e.publishRoundMatches(ctx, &event)
		if err != nil {
			return fmt.Errorf("error publishing event: %w", err)
		}
	}

	return nil
}

func (e *EventsService) publishRoundMatches(ctx context.Context, event *sofascore_models.Event) error {
	roundMatchesTopic := e.Config.KafkaConfig.TopicsName.SofascoreLeagueRoundMatches

	value, err := json.Marshal(&event)
	if err != nil {
		return fmt.Errorf("error marshalling event id %d: %w", event.ID, err)
	}

	key := []byte(fmt.Sprintf("%d-%s-%s", event.ID, event.Season.Name, event.Tournament.Name))

	if err := e.Producer.Publish(ctx, roundMatchesTopic, key, value); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}
