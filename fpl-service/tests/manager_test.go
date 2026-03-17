package main

import (
	"context"
	"log"
	"testing"

	"github.com/imadeddine-belkat/fpl-service/config"
	fpl_api "github.com/imadeddine-belkat/fpl-service/internal/api"
	. "github.com/imadeddine-belkat/fpl-service/internal/services"
	kafka "github.com/imadeddine-belkat/tactify-kafka"
)

func TestManagersApiService(t *testing.T) {
	ctx := context.Background()
	if testing.Short() {
		t.Skip("Skipping real API test")
	}

	producer := kafka.NewProducer()
	defer func(producer *kafka.Producer) {
		err := producer.Close()
		if err != nil {
			log.Printf("Error closing producer: %v", err)
			return
		}
	}(producer)

	service := &ManagersApiService{
		Config:   config.LoadConfig(),
		Client:   fpl_api.NewFplApiClient(config.LoadConfig()),
		Producer: producer,
	}
	for n := 1; n < 29; n++ {
		t.Logf("Waiting for API to be ready... (attempt %d)", n+1)
		err := service.UpdateManager(ctx, 2839296, n)
		if err != nil {
			t.Fatalf("UpdateManager with API failed: %v", err)
		}
	}

	t.Logf("Real API test completed successfully")

}
