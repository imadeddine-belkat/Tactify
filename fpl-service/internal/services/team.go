package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	"github.com/imadbelkat1/fpl-service/internal/models"
	"github.com/imadbelkat1/kafka"
)

type TeamApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *TeamApiService) getBootstrapData(ctx context.Context) (*models.BootstrapResponse, error) {
	var bootstrap models.BootstrapResponse
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

func (s *TeamApiService) publishTeams(ctx context.Context, teams []models.Team) error {
	teamsTopic := s.Config.KafkaConfig.TopicsName.FplTeams

	jobs := make(chan models.Team, len(teams))

	var publishWg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for team := range jobs {
				dto := models.TeamDTO{
					ID:                  team.ID,
					Code:                team.Code,
					Name:                team.Name,
					ShortName:           team.ShortName,
					Strength:            team.Strength,
					Form:                team.Form,
					Position:            team.Position,
					Points:              team.Points,
					Played:              team.Played,
					Win:                 team.Win,
					Draw:                team.Draw,
					Loss:                team.Loss,
					TeamDivision:        team.TeamDivision,
					Unavailable:         team.Unavailable,
					StrengthOverallHome: team.StrengthOverallHome,
					StrengthOverallAway: team.StrengthOverallAway,
					StrengthAttackHome:  team.StrengthAttackHome,
					StrengthAttackAway:  team.StrengthAttackAway,
					StrengthDefenceHome: team.StrengthDefenceHome,
					StrengthDefenceAway: team.StrengthDefenceAway,
				}
				value, err := json.Marshal(dto)
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
