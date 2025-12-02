package services

import (
	"context"

	"github.com/imadeddine-belkat/kafka"
	"github.com/imadeddine-belkat/sofascore-service/config"
	sofascore_api "github.com/imadeddine-belkat/sofascore-service/internal/api"
)

type PlayersService struct {
	Config   config.SofascoreConfig
	Client   *sofascore_api.SofascoreApiClient
	Producer *kafka.Producer
}

// get unique tour id -> fetch its standing -> get team ids from it -> fetch players

func (p *PlayersService) GetPlayerInfo(ctx context.Context, playerId int) error {
	return nil
}
