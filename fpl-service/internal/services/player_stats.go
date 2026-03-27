package services

import (
	"context"
	"fmt"
	"log"

	"github.com/imadeddine-belkat/fpl-service/config"
	api "github.com/imadeddine-belkat/fpl-service/internal/api"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
	kafkaProto "github.com/imadeddine-belkat/tactify-protos/go/kafka/v1"
	plProto "github.com/imadeddine-belkat/tactify-protos/go/pl/v1"
)

type PlayerStatsService struct {
	Config   *config.Config
	Client   *api.PlApiClient
	Producer *kafka.Producer
}
