package services

import (
	"context"
	"fmt"
	"log"

	"github.com/imadeddine-belkat/fpl-service/config"
	api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	fplProto "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
	kafkaProto "github.com/imadeddine-belkat/tactify-protos/go/kafka/v1"
)

type PlayerStaticService struct {
	Config   *config.Config
	Client   *api.FplApiClient
	Producer *kafka.Producer
}

func (s *PlayerStaticService) GetAndPublishPlayersStatic(ctx context.Context) error {
	playersBootstrap, err := s.Client.GetPlayersBootstrap(ctx)
	if err != nil {
		return fmt.Errorf("fetching players bootstrap data: %w", err)
	}

	// Step 1: Publish core player data first
	if err := s.publishPlayersCore(ctx, playersBootstrap); err != nil {
		return fmt.Errorf("publishing players core: %w", err)
	}

	// Step 2: Publish player seasons data
	if err := s.publishPlayerSeasons(ctx, playersBootstrap); err != nil {
		return fmt.Errorf("publishing player seasons: %w", err)
	}

	return nil
}

// publishPlayersBootstrap publishes only the bootstrap/stats data
func (s *PlayerStaticService) publishPlayersCore(ctx context.Context, bootstrap *fplProto.PlayersBootstrap) error {
	playersCoreTopic := s.Config.KafkaConfig.TopicsName.PlayersCore.Name

	log.Printf("Publishing %d player core records...", len(bootstrap.GetElements()))

	for _, player := range bootstrap.GetElements() {

		playerCore := &kafkaProto.PlayerCore{
			Code:       player.GetCode(),
			FirstName:  player.GetFirstName(),
			SecondName: player.GetSecondName(),
			WebName:    player.GetWebName(),
			BirthDate:  player.GetBirthDate(),
			OptaCode:   player.GetOptaCode(),
		}

		key := []byte(fmt.Sprintf("player_%d", player.GetCode()))

		if err := s.Producer.PublishWithProcess(ctx, playerCore, playersCoreTopic, key); err != nil {
			fmt.Printf("error publishing player core for player ID %d: %v\n", player.GetCode(), err)
		}
	}

	log.Printf("Published all player core data")
	return nil
}

func (s *PlayerStaticService) publishPlayerSeasons(ctx context.Context, bootstrap *fplProto.PlayersBootstrap) error {
	playersSeasonsTopic := s.Config.KafkaConfig.TopicsName.PlayersSeasons.Name

	log.Printf("Publishing %d player seasons records...", len(bootstrap.GetElements()))

	for _, player := range bootstrap.GetElements() {
		playerSeasons := &kafkaProto.PlayerSeason{
			PlayerCode: player.GetCode(),
			SeasonCode: s.Config.CurrentSeasonID,
			PlayerId:   player.GetId(),
		}

		key := []byte(fmt.Sprintf("player_%d", player.GetCode()))
		if err := s.Producer.PublishWithProcess(ctx, playerSeasons, playersSeasonsTopic, key); err != nil {
			fmt.Printf("error publishing player seasons for player ID %d: %v\n", player.GetCode(), err)
		}
	}

	log.Printf("Published all player seasons data")
	return nil
}
