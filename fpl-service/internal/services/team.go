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

type TeamApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *TeamApiService) getBootstrapData(ctx context.Context) (*fpl_models.BootstrapResponse, error) {
	var bootstrap fpl_models.BootstrapResponse
	endpoint := s.Config.FplApi.Bootstrap

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &bootstrap); err != nil {
		return nil, err
	}
	return &bootstrap, nil
}

func (s *TeamApiService) UpdateTeams(ctx context.Context) error {
	bootstrap, err := s.getBootstrapData(ctx)
	if err != nil {
		return fmt.Errorf("fetching bootstrap data: %w", err)
	}

	if err := s.publishTeams(ctx, bootstrap.Teams); err != nil {
		return fmt.Errorf("publishing teams: %w", err)
	}

	return nil
}

func (s *TeamApiService) publishTeams(ctx context.Context, teams []fpl_models.Team) error {
	teamsTopic := s.Config.KafkaConfig.TopicsName.FplTeams

	jobs := make(chan fpl_models.Team, len(teams))

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for team := range jobs {
				message := fpl_models.TeamMessage{
					Team:     team,
					SeasonID: s.Config.FplApi.Season2526,
				}
				value, err := json.Marshal(message)
				if err != nil {
					continue
				}
				key := []byte(fmt.Sprintf("%d", team.ID))
				_ = s.Producer.Publish(ctx, teamsTopic, key, value)
			}
		}()
	}

	for _, team := range teams {
		jobs <- team
	}
	close(jobs)

	publishWg.Wait()

	return nil
}
