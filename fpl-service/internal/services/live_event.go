package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/fpl_models"
)

type LiveEventApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *LiveEventApiService) UpdateLiveEvent(ctx context.Context, eventID int) error {
	liveEvent, err := s.GetLiveEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get live event data: %v", err)
	}

	if err := s.publishLiveEvent(ctx, liveEvent, eventID); err != nil {
		return fmt.Errorf("failed to publish live event data: %v", err)
	}

	return nil
}

func (s *LiveEventApiService) GetLiveEvent(ctx context.Context, eventID int) (*fpl_models.LiveEvent, error) {
	var liveEvent fpl_models.LiveEvent

	endpoint := fmt.Sprintf(s.Config.FplApi.LiveEvent, eventID)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &liveEvent); err != nil {
		return nil, fmt.Errorf("failed to fetch live event data: %v", err)
	}

	return &liveEvent, nil
}

func (s *LiveEventApiService) publishLiveEvent(ctx context.Context, liveEvent *fpl_models.LiveEvent, eventID int) error {
	liveEventTopic := s.Config.KafkaConfig.TopicsName.FplLiveEvent

	jobs := make(chan fpl_models.LiveElement, len(liveEvent.Elements))

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for element := range jobs {
				dto := fpl_models.LiveEventMessage{
					PlayerID: element.ID,
					Event:    eventID,
					SeasonID: s.Config.FplApi.CurrentSeasonID,
					Stats:    element.Stats,
					Explain:  element.Explain,
					Modified: false,
				}
				key := []byte(fmt.Sprintf("%d-%d", eventID, element.ID))
				value, err := json.Marshal(dto)
				if err != nil {
					continue
				}
				_ = s.Producer.Publish(ctx, liveEventTopic, key, value)

			}
		}()
	}

	// feed jobs
	for _, element := range liveEvent.Elements {
		jobs <- element
	}
	close(jobs)

	publishWg.Wait()
	return nil
}
