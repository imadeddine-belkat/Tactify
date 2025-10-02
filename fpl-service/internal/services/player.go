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

type PlayerApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *PlayerApiService) UpdatePlayers(ctx context.Context) error {
	playersBootstrap, err := s.getPlayersBootstrap(ctx)
	if err != nil {
		fmt.Errorf("fetching players bootstrap data: %w", err)
	}
	if err := s.publishPlayers(ctx, playersBootstrap); err != nil {
		return fmt.Errorf("publishing players data: %w", err)
	}
	return nil
}

func (s *PlayerApiService) publishPlayers(ctx context.Context, bootstrap *models.PlayersBootstrap) error {
	playersStatsTopic := s.Config.KafkaConfig.TopicsName.FplPlayersStats
	playersMatchStatsTopic := s.Config.KafkaConfig.TopicsName.FplPlayerMatchStats
	playersPastHistoryTopic := s.Config.KafkaConfig.TopicsName.FplPlayerHistoryStats

	jobHistory := make(chan models.PlayerBootstrap, len(bootstrap.PlayerBootstrap))

	var wg sync.WaitGroup
	for i := 0; i < s.Config.PublishWorkerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for player := range jobHistory {
				playerSummary, err := s.getPlayerSummary(ctx, player.ID)
				if err != nil {
					fmt.Printf("error fetching player summary for player ID %d: %v\n", player.ID, err)
					continue
				}
				playerHistory := playerSummary.PlayerHistory
				playerPast := playerSummary.PlayerPast

				valueSummary, err := json.Marshal(&player)
				if err != nil {
					fmt.Errorf("marshaling player summary for player ID %d: %w", player.ID, err)
					continue
				}

				valueHistory, err := json.Marshal(playerHistory)
				if err != nil {
					fmt.Errorf("marshaling player history for player ID %d: %w", player.ID, err)
					continue
				}

				valuePastHistory, err := json.Marshal(playerPast)
				if err != nil {
					fmt.Errorf("marshaling player past history for player ID %d: %w", player.ID, err)
					continue
				}

				playerKey := []byte(fmt.Sprintf("player_summary_%d", player.ID))
				playerHistoryKey := []byte(fmt.Sprintf("player_summary_%d", player.ID))
				playerPastHistoryKey := []byte(fmt.Sprintf("player_past_history_%d", player.ID))

				_ = s.Producer.Publish(ctx, playersStatsTopic, playerKey, valueSummary)
				_ = s.Producer.Publish(ctx, playersMatchStatsTopic, playerHistoryKey, valueHistory)
				_ = s.Producer.Publish(ctx, playersPastHistoryTopic, playerPastHistoryKey, valuePastHistory)

			}
		}()
	}

	for _, player := range bootstrap.PlayerBootstrap {
		jobHistory <- player
	}
	close(jobHistory)

	wg.Wait()
	return nil
}

func (s *PlayerApiService) getPlayerSummary(ctx context.Context, playerId int) (*models.Player, error) {
	var player models.Player
	playerSummary := s.Config.FplApi.PlayerSummary // /element-summary/%d/

	endpoint := fmt.Sprintf(playerSummary, playerId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &player); err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *PlayerApiService) getPlayersBootstrap(ctx context.Context) (*models.PlayersBootstrap, error) {
	bootstrap, err := s.getBootstrapData(ctx)
	if err != nil {
		return nil, err
	}

	return &models.PlayersBootstrap{PlayerBootstrap: bootstrap.Elements}, nil
}

func (s *PlayerApiService) getBootstrapData(ctx context.Context) (*models.BootstrapResponse, error) {
	var bootstrap models.BootstrapResponse
	endpoint := s.Config.FplApi.Bootstrap

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &bootstrap); err != nil {
		return nil, err
	}
	return &bootstrap, nil
}
