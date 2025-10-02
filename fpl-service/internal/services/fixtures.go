package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	"github.com/imadbelkat1/fpl-service/internal/models"
	"github.com/imadbelkat1/kafka"
)

type FixturesApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *FixturesApiService) GetFixtures(ctx context.Context) (*models.Fixtures, error) {
	var fixtures models.Fixtures

	fixturesEndpoint := s.Config.FplApi.Fixtures

	if err := s.Client.GetAndUnmarshal(ctx, fixturesEndpoint, &fixtures); err != nil {
		return nil, err
	}

	return &fixtures, nil
}

func (s *FixturesApiService) UpdateFixtures(ctx context.Context) error {

	start := time.Now()
	fixtures, err := s.GetFixtures(ctx)
	log.Printf("API fetch took: %v", time.Since(start))

	if err != nil {
		return fmt.Errorf("fetching fixtures data: %w", err)
	}

	start = time.Now()
	if err := s.publishFixtures(ctx, fixtures); err != nil {
		return fmt.Errorf("update fixtures: %w", err)
	}
	log.Printf("Publishing took: %v", time.Since(start))

	return nil
}

func (s *FixturesApiService) publishFixtures(ctx context.Context, fixtures *models.Fixtures) error {
	fixturesTopic := s.Config.KafkaConfig.TopicsName.FplFixtures

	jobs := make(chan models.Fixture, len(*fixtures))

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for element := range jobs {
				dto := models.FixtureDTO{
					Code:                 element.Code,
					Event:                element.Event,
					Finished:             element.Finished,
					FinishedProvisional:  element.FinishedProvisional,
					ID:                   element.ID,
					KickoffTime:          element.KickoffTime,
					Minutes:              element.Minutes,
					ProvisionalStartTime: element.ProvisionalStartTime,
					Started:              element.Started,
					TeamA:                element.TeamA,
					TeamAScore:           element.TeamAScore,
					TeamH:                element.TeamH,
					TeamHScore:           element.TeamHScore,
					TeamHDifficulty:      element.TeamHDifficulty,
					TeamADifficulty:      element.TeamADifficulty,
					PulseID:              element.PulseID,
				}
				value, err := json.Marshal(dto)
				if err != nil {
					continue
				}
				key := []byte(fmt.Sprintf("%d", element.ID))
				_ = s.Producer.Publish(ctx, fixturesTopic, key, value)
			}
		}()
	}

	for _, element := range *fixtures {
		jobs <- element
	}
	close(jobs)

	publishWg.Wait()
	return nil
}
