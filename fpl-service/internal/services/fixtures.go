package services

import (
	"context"
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
	toDelete := []string{"stats"}

	jobs := make(chan models.Fixture, len(*fixtures))
	fixturesChan := make(chan config.ProcessedModel, len(*fixtures))

	var deleteWg sync.WaitGroup
	for i := 0; i < s.Config.DeleteWorkerCount; i++ {
		deleteWg.Add(1)
		go func() {
			defer deleteWg.Done()
			for element := range jobs {
				processed, err := s.Config.ProcessDelete(element, toDelete)
				if err != nil {
					continue
				}
				fixturesChan <- config.ProcessedModel{ID: element.ID, Data: processed}
			}
		}()
	}

	go func() {
		deleteWg.Wait()
		close(fixturesChan)
	}()

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for element := range fixturesChan {
				key := []byte(fmt.Sprintf("%d", element.ID))
				_ = s.Producer.Publish(ctx, fixturesTopic, key, element.Data)
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
