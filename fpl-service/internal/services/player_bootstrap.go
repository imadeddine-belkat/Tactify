package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/imadeddine-belkat/fpl-service/config"
	api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	fplProto "github.com/imadeddine-belkat/tactify-protos/go/fpl/v1"
	kafkaProto "github.com/imadeddine-belkat/tactify-protos/go/kafka/v1"
)

type PlayerBootstrapService struct {
	Config   *config.Config
	Client   *api.FplApiClient
	Producer *kafka.Producer
}

func (b *PlayerBootstrapService) IngestPlayersBootstrap(ctx context.Context) error {
	bootstrap, err := b.Client.GetPlayersBootstrap(ctx)
	if err != nil {
		return fmt.Errorf("error fetching players bootstrap data: %w", err)
	}

	processedPlayers := b.processPlayerBootstrap(bootstrap)

	if err := b.publishPlayerBootstrap(ctx, processedPlayers); err != nil {
		return fmt.Errorf("publishing player bootstrap: %w", err)
	}

	return nil
}

func (b *PlayerBootstrapService) processPlayerBootstrap(bootstrap *fplProto.PlayersBootstrap) []*kafkaProto.PlayerBootstrap {
	log.Printf("Processing %d player bootstrap records...", len(bootstrap.GetElements()))

	processedPlayers := make([]*kafkaProto.PlayerBootstrap, len(bootstrap.GetElements()))
	for _, player := range bootstrap.GetElements() {
		playerBootstrap := &kafkaProto.PlayerBootstrap{
			SeasonCode:     b.Config.CurrentSeasonID,
			PlayerCode:     player.GetCode(),
			TeamCode:       player.GetTeamCode(),
			PlayerId:       player.GetId(),
			ElementTypeId:  player.GetElementType(),
			Status:         player.GetStatus(),
			NowCost:        player.GetNowCost(),
			Photo:          player.GetPhoto(),
			SquadNumber:    player.GetSquadNumber(),
			CanTransact:    player.GetCanTransact(),
			CanSelect:      player.GetCanSelect(),
			InDreamteam:    player.GetInDreamteam(),
			DreamteamCount: player.GetDreamteamCount(),
			Special:        player.GetSpecial(),
			Removed:        player.GetRemoved(),
			Unavailable:    player.GetUnavailable(),
		}

		processedPlayers = append(processedPlayers, playerBootstrap)
	}

	log.Printf("Processed all player bootstrap data")
	return processedPlayers
}

func (b *PlayerBootstrapService) publishPlayerBootstrap(ctx context.Context, bootstrap []*kafkaProto.PlayerBootstrap) error {
	topic := b.Config.KafkaConfig.TopicsName.PlayerBootstrap.Name
	var publishErrors []error

	log.Printf("Publishing %d player bootstrap records...", len(bootstrap))

	for _, player := range bootstrap {
		player := player

		key := []byte(fmt.Sprintf("player_bootstrap_%d", player.GetPlayerCode()))
		if err := b.Producer.PublishWithProcess(ctx, player, topic, key); err != nil {
			publishErrors = append(publishErrors, fmt.Errorf("player %d: %w", player.GetPlayerCode(), err))
		}
	}

	if len(publishErrors) > 0 {
		joined := errors.Join(publishErrors...)
		log.Printf("partial publish failure: %d/%d players failed\n%v", len(publishErrors), len(bootstrap), joined)
		return joined
	}

	return nil
}
