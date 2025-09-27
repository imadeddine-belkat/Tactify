package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	"github.com/imadbelkat1/fpl-service/internal/models"

	fixtureProducer "github.com/imadbelkat1/kafka"
)

type FixturesApiService struct {
	Client *fpl_api.FplApiClient
}

func (s *FixturesApiService) UpdateFixtures() error {
	var fixtures models.Fixtures
	producer := fixtureProducer.NewProducer()
	ctx := context.Background()

	cfg := config.LoadConfig()
	endpoint := cfg.FplApiFixtures

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &fixtures); err != nil {
		return err
	}

	for _, f := range fixtures {
		// Separate stats from fixture before marshaling
		fixtureStatsJSON, err := json.Marshal(f.Stats)

		fixtureBytes, err := json.Marshal(f)
		if err != nil {
			return fmt.Errorf("failed to marshal fixture with ID: %d: %v", f.ID, err)
		}

		var newFixture map[string]interface{}
		err = json.Unmarshal(fixtureBytes, &newFixture)
		if err != nil {
			return fmt.Errorf("failed to unmarshal fixture with ID: %d: %v", f.ID, err)
		}

		delete(newFixture, "stats")
		fixtureJSON, err := json.Marshal(newFixture)

		err = fixtureProducer.Publish(ctx, cfg.FplFixturesTopic, []byte(fmt.Sprintf("%d", f.ID)), fixtureJSON)
		if err != nil {
			return fmt.Errorf("failed to publish fixture with ID: %d to Kafka: %v", f.ID, err)
		}
	}
	return nil
}
