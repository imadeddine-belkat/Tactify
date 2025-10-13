package sofascore_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/imadbelkat1/indexer-service/config"
	"github.com/imadbelkat1/indexer-service/internal/sofascore_repositories"
	"github.com/imadbelkat1/kafka"
	kafkaConfig "github.com/imadbelkat1/kafka/config"
	"github.com/imadbelkat1/shared/sofascore_models"
)

type Handler struct {
	config      *config.IndexerConfig
	kafkaConfig *kafkaConfig.KafkaConfig
	consumers   map[string]*kafka.Consumer
	teamRepo    *sofascore_repositories.TeamRepo
}

func NewHandler(
	cfg *config.IndexerConfig,
	kafkaCfg *kafkaConfig.KafkaConfig,
	teamRepo *sofascore_repositories.TeamRepo,
) *Handler {
	h := &Handler{
		config:      cfg,
		kafkaConfig: kafkaCfg,
		consumers:   make(map[string]*kafka.Consumer),
		teamRepo:    teamRepo,
	}

	if teamRepo != nil {
		h.consumers[kafkaCfg.TopicsName.SofascoreTopTeamsStats] = kafka.NewConsumer(
			kafkaCfg,
			kafkaCfg.TopicsName.SofascoreTopTeamsStats,
			kafkaCfg.ConsumersGroupID.SofascoreTeamOverallStats,
		)
	}

	return h
}

type HandlerFunc func(ctx context.Context)

func (h *Handler) Route(ctx context.Context, topic string) {

	handlers := map[string]HandlerFunc{
		h.kafkaConfig.TopicsName.SofascoreTopTeamsStats: h.handleTopTeamsStats,
	}

	if fn, ok := handlers[topic]; ok {
		go fn(ctx)
	} else {
		log.Printf("[WARN] No handler found for topic=%s\n", topic)
	}
}

// Generic batch processor - handles all the common batching logic
func batchProcess[T any, K comparable](
	ctx context.Context,
	consumer *kafka.Consumer,
	batchSize int,
	flushInterval time.Duration,
	topicName string,
	getKey func(T) K,
	process func(T) error,
) {

	messages, errors := consumer.Subscribe(ctx)
	batch := make(map[K]T)
	flushTicker := time.NewTicker(flushInterval)
	defer flushTicker.Stop()

	flushBatch := func(logContext string) {
		if len(batch) == 0 {
			return
		}

		for _, item := range batch {
			if err := process(item); err != nil {
				log.Printf("Error %s for %s: %v\n", logContext, topicName, err)
			}
		}
		batch = make(map[K]T)
	}

	for {
		select {
		case msg := <-messages:

			var item T
			if err := json.Unmarshal(msg.Value, &item); err != nil {
				continue
			}

			batch[getKey(item)] = item

			if len(batch) >= batchSize {
				flushBatch("inserting")
			}

		case <-flushTicker.C:
			flushBatch("flushing")

		case err := <-errors:
			if err != nil {
				log.Printf("Error consuming %s message: %v\n", topicName, err)
			}

		case <-ctx.Done():
			flushBatch("inserting on shutdown")
			return
		}
	}
}

// Generic batch processor with slice conversion - for handlers that need to convert map to slice
func batchProcessWithSlice[T any, K comparable](
	ctx context.Context,
	consumer *kafka.Consumer,
	batchSize int,
	flushInterval time.Duration,
	topicName string,
	getKey func(T) K,
	processBatch func([]T) error,
) {
	messages, errors := consumer.Subscribe(ctx)
	batch := make(map[K]T)
	flushTicker := time.NewTicker(flushInterval)
	defer flushTicker.Stop()

	flushBatch := func(logContext string) {
		if len(batch) == 0 {
			return
		}

		items := make([]T, 0, len(batch))
		for _, item := range batch {
			items = append(items, item)
		}

		if err := processBatch(items); err != nil {
			log.Printf("Error %s batch for %s: %v\n", logContext, topicName, err)
		}
		batch = make(map[K]T)
	}

	for {
		select {
		case msg := <-messages:
			var item T
			if err := json.Unmarshal(msg.Value, &item); err != nil {
				log.Printf("Error unmarshaling %s message: %v, raw: %s\n", topicName, err, string(msg.Value))
				continue
			}
			batch[getKey(item)] = item

			if len(batch) >= batchSize {
				flushBatch("inserting")
			}

		case <-flushTicker.C:
			flushBatch("flushing")

		case err := <-errors:
			if err != nil {
				log.Printf("Error consuming %s message: %v\n", topicName, err)
			}

		case <-ctx.Done():
			flushBatch("inserting on shutdown")
			return
		}
	}
}

func (h *Handler) handleTopTeamsStats(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.SofascoreTopTeamsStats],
		1,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.SofascoreTopTeamsStats,
		func(topT sofascore_models.TopTeamsMessage) string {
			return fmt.Sprintf("%d", time.Now().UnixNano())
		},
		h.teamRepo.InsertTeamOverallStats,
	)

}
