package services

import (
	"context"
	"fmt"
	"log"

	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	sofascore "github.com/imadeddine-belkat/tactify-protos/go/sofascore/v1"
)

type EventsService struct {
	Config   *config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

func (e *EventsService) UpdateRoundMatches(ctx context.Context, seasonId, leagueId, round int) error {
	roundMatches, err := e.GetRoundMatches(ctx, seasonId, leagueId, round)
	if err != nil {
		return fmt.Errorf("error fetching events data: %w", err)
	}

	for _, event := range roundMatches {
		err := e.publishRoundMatches(ctx, event)
		if err != nil {
			return fmt.Errorf("error publishing event: %w", err)
		}
	}

	return nil
}

func (e *EventsService) GetRoundMatches(ctx context.Context, seasonId, leagueId, round int) ([]*sofascore.Event, error) {
	var events []*sofascore.Event
	leagueRoundMatches := e.Config.SofascoreApi.LeagueEndpoints.LeagueRoundMatches //unique-tournament/%d/seasonId/%d/events/round/%d
	log.Printf("leagueRoundMatches raw: %q", leagueRoundMatches)

	endpoint := fmt.Sprintf(leagueRoundMatches, leagueId, seasonId, round)
	log.Println(endpoint)

	if err := e.Client.GetAndUnmarshal(ctx, endpoint, events); err != nil {
		return nil, fmt.Errorf("fetching events data: %w", err)
	}

	return events, nil
}

func (e *EventsService) publishRoundMatches(ctx context.Context, event *sofascore.Event) error {
	roundMatchesTopic := e.Config.KafkaConfig.TopicsName.SofascoreLeagueRoundMatches.Name

	key := []byte(fmt.Sprintf("%d-%s-%s", event.Id, event.Season.Name, event.Tournament.Name))

	if err := e.Producer.PublishWithProcess(ctx, event, roundMatchesTopic, key); err != nil {
		return fmt.Errorf("error publishing event: %w", err)
	}

	return nil
}
