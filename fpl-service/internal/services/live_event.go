package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
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

func (s *LiveEventApiService) GetLiveEvent(ctx context.Context, eventID int) (*fpl.LiveEvent, error) {
	var liveEvent fpl.LiveEvent

	endpoint := fmt.Sprintf(s.Config.FplApi.LiveEvent, eventID)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &liveEvent); err != nil {
		return nil, fmt.Errorf("failed to fetch live event data: %v", err)
	}

	return &liveEvent, nil
}

func (s *LiveEventApiService) publishLiveEvent(ctx context.Context, liveEvent *fpl.LiveEvent, eventID int) error {
	liveEventTopic := s.Config.KafkaConfig.TopicsName.FplLiveEvent.Name

	jobs := make(chan *fpl.LiveElement, len(liveEvent.Elements))

	var publishWg sync.WaitGroup
	for i := 0; i < 10; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for element := range jobs {
				dto := &fpl.LiveEventMessage{
					PlayerId: element.Id,
					Event:    int32(eventID),
					SeasonId: s.Config.FplApi.CurrentSeasonID,
					Stats:    element.Stats,
					Explain:  element.Explain,
					Modified: false,
				}
				key := []byte(fmt.Sprintf("%d-%d", eventID, element.Id))
				err := s.Producer.PublishWithProcess(ctx, dto, liveEventTopic, key)
				if err != nil {
					fmt.Printf("failed to publish live event element: %v for player ID: %d and event ID: %d\n", err, element.Id, eventID)
				}

			}
		}()
	}

	// feed jobs
	for _, element := range liveEvent.Elements {
		jobs <- element
	}
	close(jobs)

	publishWg.Wait()
	return s.Producer.Close()
}
