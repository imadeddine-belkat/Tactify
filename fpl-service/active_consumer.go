package main

import (
	"context"
	"fmt"

	"github.com/imadbelkat1/kafka"
	kafkaConfig "github.com/imadbelkat1/kafka/config"
)

func main() {
	kafkaCfg := kafkaConfig.LoadConfig()

	// Get all available topics
	topics := []string{
		kafkaCfg.TopicsName.FplTeams,
		kafkaCfg.TopicsName.FplFixtures,
		kafkaCfg.TopicsName.FplPlayerMatchStats,
		kafkaCfg.TopicsName.FplLiveEvent,
		// Add other topics here as needed
	}

	consumerGrp := map[string]string{
		kafkaCfg.TopicsName.FplTeams:            kafkaCfg.ConsumersGroupID.Teams,
		kafkaCfg.TopicsName.FplFixtures:         kafkaCfg.ConsumersGroupID.Fixtures,
		kafkaCfg.TopicsName.FplPlayerMatchStats: kafkaCfg.ConsumersGroupID.PlayersStats,
		kafkaCfg.TopicsName.FplLiveEvent:        kafkaCfg.ConsumersGroupID.Live,
	}

	for _, topic := range topics {
		go func() {
			ctx := context.Background()
			consumer := kafka.NewConsumer(kafkaCfg, topic, consumerGrp[topic])
			messages, errors := consumer.Subscribe(ctx)

			fmt.Printf("Starting to listen on topic: %s\n", topic)

			for {
				select {
				case msg := <-messages:
					fmt.Printf("Topic [%s] - Received message: key=%s, value=%s\n",
						topic, string(msg.Key), string(msg.Value))

				case err := <-errors:
					if err != nil {
						fmt.Printf("Topic [%s] - Consumer error: %v\n", topic, err)
					}

				case <-ctx.Done():
					fmt.Printf("Topic [%s] - Consumer stopped\n", topic)
					return
				}
			}
		}() // Adjust index as needed
	}

	// Keep program running to listen for messages
	select {}
}
