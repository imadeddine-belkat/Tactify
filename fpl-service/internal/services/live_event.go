package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	"github.com/imadbelkat1/fpl-service/internal/models"
	"github.com/imadbelkat1/kafka"
)

type LiveEventApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *LiveEventApiService) GetLiveEvent(ctx context.Context, eventID string) (*models.LiveEvent, error) {
	var liveEvent models.LiveEvent

	endpoint := fmt.Sprintf(s.Config.FplApi.LiveEvent, eventID)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &liveEvent); err != nil {
		return nil, fmt.Errorf("failed to fetch live event data: %v", err)
	}

	return &liveEvent, nil
}

func (s *LiveEventApiService) UpdateLiveEvent(ctx context.Context, eventID string) error {
	liveEvent, err := s.GetLiveEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get live event data: %v", err)
	}

	if err := s.publishLiveEvent(ctx, liveEvent, eventID); err != nil {
		return fmt.Errorf("failed to publish live event data: %v", err)
	}

	return nil
}

func (s *LiveEventApiService) publishLiveEvent(ctx context.Context, liveEvent *models.LiveEvent, eventID string) error {
	liveEventTopic := s.Config.KafkaConfig.TopicsName.FplLiveEvent
	gameweek, err := strconv.Atoi(eventID)
	if err != nil {
		return fmt.Errorf("invalid event ID: %v", err)
	}

	toDelete := []string{"explain", "modified"}

	jobs := make(chan models.LiveElement, len(liveEvent.Elements))
	elements := make(chan config.ProcessedModel, len(liveEvent.Elements))

	// delete stage
	var deleteWg sync.WaitGroup
	for i := 0; i < s.Config.DeleteWorkerCount; i++ {
		deleteWg.Add(1)
		go func() {
			defer deleteWg.Done()
			for element := range jobs {
				element.Gameweek = gameweek
				processed, err := s.Config.ProcessDelete(element, toDelete)
				if err != nil {
					continue
				}
				elements <- config.ProcessedModel{ID: element.ID, Data: processed}
			}
		}()
	}

	go func() {
		deleteWg.Wait()
		close(elements)
	}()

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for element := range elements {
				key := []byte(fmt.Sprintf("%d-%d", gameweek, element.ID))
				_ = s.Producer.Publish(ctx, liveEventTopic, key, element.Data)

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
