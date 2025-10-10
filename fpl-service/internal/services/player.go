package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadbelkat1/fpl-service/config"
	fpl_api "github.com/imadbelkat1/fpl-service/internal/api"
	"github.com/imadbelkat1/kafka"
	"github.com/imadbelkat1/shared/models"
)

type PlayerApiService struct {
	Config   *config.FplConfig
	Client   *fpl_api.FplApiClient
	Producer *kafka.Producer
}

func (s *PlayerApiService) UpdatePlayers(ctx context.Context) error {
	playersBootstrap, err := s.getPlayersBootstrap(ctx)
	if err != nil {
		return fmt.Errorf("fetching players bootstrap data: %w", err)
	}

	// Step 1: Publish bootstrap data first
	if err := s.publishPlayersBootstrap(ctx, playersBootstrap); err != nil {
		return fmt.Errorf("publishing players bootstrap: %w", err)
	}

	// Step 2: Publish detailed histories
	if err := s.publishPlayersHistory(ctx, playersBootstrap); err != nil {
		return fmt.Errorf("publishing players history: %w", err)
	}

	return nil
}

// publishPlayersBootstrap publishes only the bootstrap/stats data
func (s *PlayerApiService) publishPlayersBootstrap(ctx context.Context, bootstrap *models.PlayersBootstrap) error {
	playersBootstrapTopic := s.Config.KafkaConfig.TopicsName.FplPlayersBootstrap

	log.Printf("Publishing %d player bootstrap records...", len(bootstrap.PlayerBootstrap))

	for _, player := range bootstrap.PlayerBootstrap {
		playerBootstrap := models.PlayerBootstrapMessage{
			Player:   player,
			SeasonID: s.Config.FplApi.CurrentSeasonID,
		}
		valueSummary, err := json.Marshal(&playerBootstrap)
		if err != nil {
			fmt.Printf("error marshaling player bootstrap for player ID %d: %v\n", player.ID, err)
			continue
		}
		playerKey := []byte(fmt.Sprintf("player_%d", player.ID))
		if err := s.Producer.Publish(ctx, playersBootstrapTopic, playerKey, valueSummary); err != nil {
			fmt.Printf("error publishing player bootstrap for player ID %d: %v\n", player.ID, err)
		}
	}

	log.Printf("Published all player bootstrap data")
	return nil
}

// publishPlayersHistory fetches and publishes match history and past seasons
func (s *PlayerApiService) publishPlayersHistory(ctx context.Context, bootstrap *models.PlayersBootstrap) error {
	playersMatchStatsTopic := s.Config.KafkaConfig.TopicsName.FplPlayerMatchStats
	playersPastHistoryTopic := s.Config.KafkaConfig.TopicsName.FplPlayerHistoryStats

	log.Printf("Fetching and publishing history for %d players...", len(bootstrap.PlayerBootstrap))

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

				// Publish player match history (gameweek stats)
				if len(playerSummary.PlayerHistory) > 0 {
					matchStatsPayload := models.PlayerHistoryMessage{
						PlayerID: player.ID,
						SeasonID: s.Config.FplApi.CurrentSeasonID,
						History:  playerSummary.PlayerHistory,
					}
					valueHistory, err := json.Marshal(matchStatsPayload)
					if err != nil {
						fmt.Printf("error marshaling player history for player ID %d: %v\n", player.ID, err)
						continue
					}
					playerHistoryKey := []byte(fmt.Sprintf("player_history_%d", player.ID))
					if err := s.Producer.Publish(ctx, playersMatchStatsTopic, playerHistoryKey, valueHistory); err != nil {
						fmt.Printf("error publishing player history for player ID %d: %v\n", player.ID, err)
					}
				}

				// Publish player past seasons
				if len(playerSummary.PlayerPast) > 0 {
					// Map season names to IDs
					for i := range playerSummary.PlayerPast {
						seasonID := s.mapSeasonNameToID(playerSummary.PlayerPast[i].SeasonName)
						playerSummary.PlayerPast[i].SeasonId = seasonID
					}

					pastHistoryPayload := models.PlayerPastHistoryMessage{
						PlayerCode:        player.Code,
						PlayerPastHistory: playerSummary.PlayerPast,
					}

					valuePastHistory, err := json.Marshal(pastHistoryPayload)
					if err != nil {
						fmt.Printf("error marshaling player past history for player ID %d: %v\n", player.ID, err)
						continue
					}
					playerPastHistoryKey := []byte(fmt.Sprintf("player_past_%d", player.Code))
					if err := s.Producer.Publish(ctx, playersPastHistoryTopic, playerPastHistoryKey, valuePastHistory); err != nil {
						fmt.Printf("error publishing player past history for player code %d: %v\n", player.Code, err)
					}
				}
			}
		}()
	}

	for _, player := range bootstrap.PlayerBootstrap {
		jobHistory <- player
	}
	close(jobHistory)

	wg.Wait()
	log.Printf("Published all player history data")
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

func (s *PlayerApiService) mapSeasonNameToID(seasons string) int {
	switch seasons {
	case "2025/26":
		return s.Config.FplApi.Season2526
	case "2024/25":
		return s.Config.FplApi.Season2425
	case "2023/24":
		return s.Config.FplApi.Season2324
	case "2022/23":
		return s.Config.FplApi.Season2223
	case "2021/22":
		return s.Config.FplApi.Season2122
	case "2020/21":
		return s.Config.FplApi.Season2021
	case "2019/20":
		return s.Config.FplApi.Season1920
	case "2018/19":
		return s.Config.FplApi.Season1819
	case "2017/18":
		return s.Config.FplApi.Season1718
	case "2016/17":
		return s.Config.FplApi.Season1617
	case "2015/16":
		return s.Config.FplApi.Season1516
	case "2014/15":
		return s.Config.FplApi.Season1415
	case "2013/14":
		return s.Config.FplApi.Season1314
	case "2012/13":
		return s.Config.FplApi.Season1213
	case "2011/12":
		return s.Config.FplApi.Season1112
	case "2010/11":
		return s.Config.FplApi.Season1011
	default:
		return 0 // Unknown season
	}
}
