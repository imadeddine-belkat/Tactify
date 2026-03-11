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

type TeamApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *TeamApiService) getBootstrapData(ctx context.Context) (*fpl.BootstrapResponse, error) {
	var bootstrap *fpl.BootstrapResponse
	endpoint := s.Config.FplApi.Bootstrap

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, bootstrap); err != nil {
		return nil, err
	}
	return bootstrap, nil
}

func (s *TeamApiService) UpdateTeams(ctx context.Context) error {
	bootstrap, err := s.getBootstrapData(ctx)
	if err != nil {
		return fmt.Errorf("fetching bootstrap data: %w", err)
	}

	if err := s.publishTeams(ctx, bootstrap.GetTeams()); err != nil {
		return fmt.Errorf("publishing teams: %w", err)
	}

	return nil
}

func (s *TeamApiService) publishTeams(ctx context.Context, teams []*fpl.Team) error {
	teamsTopic := s.Config.KafkaConfig.TopicsName.FplTeams.Name

	jobs := make(chan *fpl.Team, len(teams))

	var publishWg sync.WaitGroup
	for i := 0; i < 10; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for team := range jobs {
				message := &fpl.TeamMessage{
					Team:     team,
					SeasonId: s.Config.FplApi.CurrentSeasonID,
				}
				key := []byte(fmt.Sprintf("%d", team.GetId()))
				err := s.Producer.PublishWithProcess(ctx, message, teamsTopic, key)
				if err != nil {
					fmt.Printf("failed to publish team %d: %v\n", team.GetId(), err)
				}
			}
		}()
	}

	for _, team := range teams {
		jobs <- team
	}
	close(jobs)

	publishWg.Wait()

	return s.Producer.Close()
}
