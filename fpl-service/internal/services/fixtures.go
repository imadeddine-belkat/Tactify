package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/shared/fpl_models"
)

type FixturesApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *FixturesApiService) GetFixtures(ctx context.Context) (*fpl_models.Fixtures, error) {
	var fixtures fpl_models.Fixtures

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

func (s *FixturesApiService) publishFixtures(ctx context.Context, fixtures *fpl_models.Fixtures) error {
	fixturesTopic := s.Config.KafkaConfig.TopicsName.FplFixtures
	//fixtureDetailsTopic := s.Config.KafkaConfig.TopicsName.FplFixtureDetails

	jobs := make(chan fpl_models.Fixture, len(*fixtures))

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for fixture := range jobs {
				fixtureMessage := fpl_models.FixtureMessage{
					Fixture:  fixture,
					SeasonID: s.Config.FplApi.CurrentSeasonID,
				}
				value, err := json.Marshal(fixtureMessage)
				if err != nil {
					continue
				}

				/*dtoStats := fpl_models.FixtureStatDTO{
					ID:          fixture.ID,
					FixtureStat: fixture.Stats,
				}
				valueStats, err := json.Marshal(dtoStats)
				if err != nil {
					continue
				}*/

				key := []byte(fmt.Sprintf("%d", fixture.ID))
				_ = s.Producer.Publish(ctx, fixturesTopic, key, value)
				//_ = s.Producer.Publish(ctx, fixtureDetailsTopic, key, valueStats)
			}

		}()
	}

	for _, fixture := range *fixtures {
		jobs <- fixture
	}
	close(jobs)

	publishWg.Wait()
	return nil
}
