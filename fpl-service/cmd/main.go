package main

import (
	"context"
	"log"

	"github.com/imadeddine-belkat/fpl-service/config"
	api "github.com/imadeddine-belkat/fpl-service/internal/api"
	"github.com/imadeddine-belkat/fpl-service/internal/services"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()
	client := api.NewFplApiClient(cfg)
	producer := kafka.NewProducer()

	service := &services.PlayerBootstrapService{
		Config:   cfg,
		Client:   client,
		Producer: producer,
	}

	if err := service.IngestPlayersBootstrap(ctx); err != nil {
		log.Fatalf("Error ingesting player bootstrap data: %v", err)
	}

	log.Println("Player bootstrap ingestion completed successfully")
}
