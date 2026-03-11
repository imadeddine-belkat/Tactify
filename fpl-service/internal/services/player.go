package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	fpl "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
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
func (s *PlayerApiService) publishPlayersBootstrap(ctx context.Context, bootstrap *fpl.PlayersBootstrap) error {
	playersBootstrapTopic := s.Config.KafkaConfig.TopicsName.FplPlayersBootstrap.Name

	log.Printf("Publishing %d player bootstrap records...", len(bootstrap.GetElements()))

	for _, player := range bootstrap.Elements {
		playerBootstrap := fpl.PlayerBootstrapMessage{
			Player:   player,
			SeasonId: s.Config.FplApi.CurrentSeasonID,
		}
		valueSummary, err := json.Marshal(&playerBootstrap)
		if err != nil {
			fmt.Printf("error marshaling player bootstrap for player ID %d: %v\n", player.GetId(), err)
			continue
		}
		playerKey := []byte(fmt.Sprintf("player_%d", player.GetId()))
		if err := s.Producer.Publish(ctx, playersBootstrapTopic, playerKey, valueSummary); err != nil {
			fmt.Printf("error publishing player bootstrap for player ID %d: %v\n", player.GetId(), err)
		}
	}

	log.Printf("Published all player bootstrap data")
	return nil
}

// publishPlayersHistory fetches and publishes match history and past seasons
func (s *PlayerApiService) publishPlayersHistory(ctx context.Context, bootstrap *fpl.PlayersBootstrap) error {
	playersMatchStatsTopic := s.Config.KafkaConfig.TopicsName.FplPlayerMatchStats.Name
	playersPastHistoryTopic := s.Config.KafkaConfig.TopicsName.FplPlayerHistoryStats.Name

	log.Printf("Fetching and publishing history for %d players...", len(bootstrap.GetElements()))

	jobHistory := make(chan *fpl.PlayerBootstrap, len(bootstrap.GetElements()))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for player := range jobHistory {
				playerSummary, err := s.getPlayerSummary(ctx, player.Id)
				if err != nil {
					fmt.Printf("error fetching player summary for player ID %d: %v\n", player.Id, err)
					continue
				}

				// Publish player match history (gameweek stats)
				if len(playerSummary.History) > 0 {
					matchStatsPayload := fpl.PlayerHistoryMessage{
						PlayerId: player.GetId(),
						SeasonId: s.Config.FplApi.CurrentSeasonID,
						History:  playerSummary.GetHistory(),
					}
					valueHistory, err := json.Marshal(&matchStatsPayload)
					if err != nil {
						fmt.Printf("error marshaling player history for player ID %d: %v\n", player.Id, err)
						continue
					}
					playerHistoryKey := []byte(fmt.Sprintf("player_history_%d", player.Id))
					if err := s.Producer.Publish(ctx, playersMatchStatsTopic, playerHistoryKey, valueHistory); err != nil {
						fmt.Printf("error publishing player history for player ID %d: %v\n", player.Id, err)
					}
				}

				// Publish player past seasons
				if len(playerSummary.GetHistoryPast()) > 0 {
					// Map season names to IDs
					for i, past := range playerSummary.GetHistoryPast() {
						seasonID := s.Config.MapSeasonNameToID(past.SeasonName)
						playerSummary.HistoryPast[i].SeasonId = seasonID
					}

					pastHistoryPayload := fpl.PlayerPastHistoryMessage{
						ElementCode: player.GetCode(),
						PastHistory: playerSummary.GetHistoryPast(),
					}

					valuePastHistory, err := json.Marshal(&pastHistoryPayload)
					if err != nil {
						fmt.Printf("error marshaling player past history for player ID %d: %v\n", player.GetId(), err)
						continue
					}
					playerPastHistoryKey := []byte(fmt.Sprintf("player_past_%d", player.GetCode()))
					if err := s.Producer.Publish(ctx, playersPastHistoryTopic, playerPastHistoryKey, valuePastHistory); err != nil {
						fmt.Printf("error publishing player past history for player code %d: %v\n", player.GetCode(), err)
					}
				}
			}
		}()
	}

	for _, player := range bootstrap.GetElements() {
		jobHistory <- player
	}
	close(jobHistory)

	wg.Wait()
	log.Printf("Published all player history data")
	return nil
}

func (s *PlayerApiService) getPlayerSummary(ctx context.Context, playerId int32) (*fpl.Player, error) {
	var player fpl.Player
	playerSummary := s.Config.FplApi.PlayerSummary // /element-summary/%d/

	endpoint := fmt.Sprintf(playerSummary, playerId)

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &player); err != nil {
		return nil, err
	}

	return &player, nil
}

func (s *PlayerApiService) getPlayersBootstrap(ctx context.Context) (*fpl.PlayersBootstrap, error) {
	bootstrap, err := s.getBootstrapData(ctx)
	if err != nil {
		return nil, err
	}

	return &fpl.PlayersBootstrap{Elements: bootstrap.GetElements()}, nil
}

func (s *PlayerApiService) getBootstrapData(ctx context.Context) (*fpl.BootstrapResponse, error) {
	var bootstrap fpl.BootstrapResponse
	endpoint := s.Config.FplApi.Bootstrap

	if err := s.Client.GetAndUnmarshal(ctx, endpoint, &bootstrap); err != nil {
		return nil, err
	}
	return &bootstrap, nil
}
