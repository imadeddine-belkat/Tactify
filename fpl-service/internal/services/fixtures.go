package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
)

type FixturesApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *FixturesApiService) GetFixtures(ctx context.Context) ([]*fpl.Fixture, error) {
	var fixtures []*fpl.Fixture

	fixturesEndpoint := s.Config.FplApi.Fixtures

	if err := s.Client.GetAndUnmarshal(ctx, fixturesEndpoint, &fixtures); err != nil {
		return nil, err
	}

	return fixtures, nil
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

func (s *FixturesApiService) publishFixtures(ctx context.Context, fixtures []*fpl.Fixture) error {
	fixturesTopic := s.Config.KafkaConfig.TopicsName.FplFixtures.Name

	jobs := make(chan *fpl.Fixture, len(fixtures))

	var publishWg sync.WaitGroup
	for i := 0; i < 10; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for fixture := range jobs {
				fixtureMessage := &fpl.FixtureMessage{
					Fixture:  fixture,
					SeasonId: s.Config.FplApi.CurrentSeasonID,
				}

				key := []byte(fmt.Sprintf("%d", fixture.Id))
				err := s.Producer.PublishWithProcess(ctx, fixtureMessage, fixturesTopic, key)
				if err != nil {
					fmt.Printf("Failed to publish fixture message for fixture ID %d: %v\n", fixture.Id, err)
				}
			}

		}()
	}

	for _, fixture := range fixtures {
		jobs <- fixture
	}
	close(jobs)

	publishWg.Wait()
	return s.Producer.Close()
}
