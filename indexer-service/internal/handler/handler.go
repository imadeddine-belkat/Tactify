package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/imadbelkat1/indexer-service/config"
	"github.com/imadbelkat1/indexer-service/internal/fpl_repositories"
	"github.com/imadbelkat1/kafka"
	kafkaConfig "github.com/imadbelkat1/kafka/config"
	"github.com/imadbelkat1/shared/models"
)

type Handler struct {
	config      *config.IndexerConfig
	kafkaConfig *kafkaConfig.KafkaConfig
	consumers   map[string]*kafka.Consumer
	playerRepo  *fpl_repositories.PlayerRepo
	teamRepo    *fpl_repositories.TeamRepo
	fixtureRepo *fpl_repositories.FixtureRepo
	managerRepo *fpl_repositories.ManagerRepo
}

func NewHandler(
	config *config.IndexerConfig,
	kafkaConfig *kafkaConfig.KafkaConfig,
	playerRepo *fpl_repositories.PlayerRepo,
	teamRepo *fpl_repositories.TeamRepo,
	fixtureRepo *fpl_repositories.FixtureRepo,
	managerRepo *fpl_repositories.ManagerRepo,
) *Handler {
	h := &Handler{
		config:      config,
		kafkaConfig: kafkaConfig,
		playerRepo:  playerRepo,
		teamRepo:    teamRepo,
		fixtureRepo: fixtureRepo,
		managerRepo: managerRepo,
		consumers:   make(map[string]*kafka.Consumer),
	}

	// Pre-create consumers only for non-nil repositories
	if fixtureRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplFixtures] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplFixtures,
			kafkaConfig.ConsumersGroupID.FplFixtures,
		)
	}

	if teamRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplTeams] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplTeams,
			kafkaConfig.ConsumersGroupID.FplTeams,
		)
	}

	if playerRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplPlayersBootstrap] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayersBootstrap,
			kafkaConfig.ConsumersGroupID.FplPlayers,
		)

		h.consumers[kafkaConfig.TopicsName.FplPlayersStats] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayersStats,
			kafkaConfig.ConsumersGroupID.FplPlayersStats,
		)

		h.consumers[kafkaConfig.TopicsName.FplPlayerMatchStats] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayerMatchStats,
			kafkaConfig.ConsumersGroupID.FplPlayersStats,
		)

		h.consumers[kafkaConfig.TopicsName.FplPlayerHistoryStats] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplPlayerHistoryStats,
			kafkaConfig.ConsumersGroupID.FplPlayersStats,
		)

		h.consumers[kafkaConfig.TopicsName.FplLiveEvent] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplLiveEvent,
			kafkaConfig.ConsumersGroupID.FplLive,
		)
	}

	if managerRepo != nil {
		h.consumers[kafkaConfig.TopicsName.FplEntry] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntry,
			kafkaConfig.ConsumersGroupID.FplEntries,
		)

		h.consumers[kafkaConfig.TopicsName.FplEntryPicks] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntryPicks,
			kafkaConfig.ConsumersGroupID.FplEntriesPicks,
		)

		h.consumers[kafkaConfig.TopicsName.FplEntryTransfers] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntryTransfers,
			kafkaConfig.ConsumersGroupID.FplEntriesTransfers,
		)

		h.consumers[kafkaConfig.TopicsName.FplEntryHistory] = kafka.NewConsumer(
			kafkaConfig,
			kafkaConfig.TopicsName.FplEntryHistory,
			kafkaConfig.ConsumersGroupID.FplEntriesHistory,
		)
	}

	return h
}

type HandlerFunc func(ctx context.Context)

func (h *Handler) Route(ctx context.Context, topic string) {
	handlers := map[string]HandlerFunc{
		h.kafkaConfig.TopicsName.FplFixtures:           h.handleFixtures,
		h.kafkaConfig.TopicsName.FplTeams:              h.handleTeams,
		h.kafkaConfig.TopicsName.FplPlayersBootstrap:   h.handlePlayerBootstrap,
		h.kafkaConfig.TopicsName.FplPlayersStats:       h.handlePlayerStats,
		h.kafkaConfig.TopicsName.FplPlayerMatchStats:   h.handlePlayerMatchStats,
		h.kafkaConfig.TopicsName.FplPlayerHistoryStats: h.handlePlayerPastHistory,
		h.kafkaConfig.TopicsName.FplLiveEvent:          h.handlePlayerExplain,
		h.kafkaConfig.TopicsName.FplEntry:              h.handleManagers,
		h.kafkaConfig.TopicsName.FplEntryPicks:         h.handleManagerPicks,
		h.kafkaConfig.TopicsName.FplEntryTransfers:     h.handleManagerTransfers,
		h.kafkaConfig.TopicsName.FplEntryHistory:       h.handleManagerHistory,
	}

	if fn, ok := handlers[topic]; ok {
		go fn(ctx)
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

func (h *Handler) handleFixtures(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplFixtures],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplFixtures,
		func(f models.FixtureMessage) int { return f.Fixture.ID },
		h.fixtureRepo.InsertFixtures,
	)
}

func (h *Handler) handleTeams(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplTeams],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplTeams,
		func(t models.TeamMessage) int { return t.Team.ID },
		h.teamRepo.InsertTeams,
	)
}

func (h *Handler) handlePlayerBootstrap(ctx context.Context) {
	consumer := h.consumers[h.kafkaConfig.TopicsName.FplPlayersBootstrap]
	messages, errors := consumer.Subscribe(ctx)
	batch := make(map[int]models.PlayerBootstrapMessage)

	flushTicker := time.NewTicker(h.config.FlushInterval)
	verifyTicker := time.NewTicker(30 * time.Second)
	defer flushTicker.Stop()
	defer verifyTicker.Stop()

	totalReceived := 0
	totalProcessed := 0

	flushBatch := func(logContext string) {
		if len(batch) == 0 {
			return
		}

		players := make([]models.PlayerBootstrapMessage, 0, len(batch))
		for _, p := range batch {
			players = append(players, p)
		}

		if err := h.playerRepo.InsertPlayerBootstrapComplete(players); err != nil {
			log.Printf("❌ Error %s player bootstrap batch: %v", logContext, err)
		} else {
			totalProcessed += len(players)
			if logContext == "inserting" {
				log.Printf("✅ Batch processed. Total received: %d, Total processed: %d", totalReceived, totalProcessed)
			} else if logContext == "flushing" {
				log.Printf("✅ Flush completed. Total received: %d, Total processed: %d", totalReceived, totalProcessed)
			}
		}
		batch = make(map[int]models.PlayerBootstrapMessage)
	}

	for {
		select {
		case msg := <-messages:
			var player models.PlayerBootstrapMessage
			if err := json.Unmarshal(msg.Value, &player); err != nil {
				log.Printf("Error unmarshaling player bootstrap message: %v, raw: %s\n", err, string(msg.Value))
				continue
			}
			batch[player.Player.ID] = player
			totalReceived++

			if len(batch) >= h.config.BatchSize {
				flushBatch("inserting")
			}

		case <-flushTicker.C:
			flushBatch("flushing")

		case <-verifyTicker.C:
			if h.playerRepo != nil {
				count, err := h.playerRepo.CountPlayers()
				if err != nil {
					log.Printf("⚠️  Error counting players: %v", err)
				} else {
					log.Printf("📊 Players in database: %d (received: %d, processed: %d)", count, totalReceived, totalProcessed)
				}
			}

		case err := <-errors:
			if err != nil {
				log.Println("Error consuming player bootstrap message:", err)
			}

		case <-ctx.Done():
			flushBatch("inserting on shutdown")
			if h.playerRepo != nil {
				count, _ := h.playerRepo.CountPlayers()
				log.Printf("🏁 Final stats - Players in DB: %d, Received: %d, Processed: %d", count, totalReceived, totalProcessed)
			}
			return
		}
	}
}

func (h *Handler) handlePlayerStats(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplPlayersStats],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplPlayersStats,
		func(p models.PlayerBootstrapMessage) int { return p.Player.ID },
		h.playerRepo.InsertPlayerBootstrapComplete,
	)
}

func (h *Handler) handlePlayerMatchStats(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplPlayerMatchStats],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplPlayerMatchStats,
		func(msg models.PlayerHistoryMessage) int { return msg.PlayerID },
		func(histMsg models.PlayerHistoryMessage) error {
			return h.playerRepo.InsertPlayerGameweekStats(histMsg)
		},
	)
}

func (h *Handler) handlePlayerPastHistory(ctx context.Context) {

	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplPlayerHistoryStats],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplPlayerHistoryStats,
		func(p models.PlayerPastHistoryMessage) int { return p.PlayerCode },
		h.playerRepo.InsertPlayerPastSeasons,
	)
}

func (h *Handler) handlePlayerExplain(ctx context.Context) {
	batchProcessWithSlice(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplLiveEvent],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplLiveEvent,
		func(p models.LiveEventMessage) int { return p.PlayerID },
		h.playerRepo.InsertPlayerGameweekExplain,
	)
}

func (h *Handler) handleManagers(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntry],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntry,
		func(e models.EntryMessage) int { return e.Entry.ID },
		func(e models.EntryMessage) error {
			return h.managerRepo.InsertManagerInfo(&e)
		},
	)
}

func (h *Handler) handleManagerPicks(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntryPicks],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntryPicks,
		func(p models.EntryEventPicksMessage) int { return p.EntryId },
		func(p models.EntryEventPicksMessage) error {
			return h.managerRepo.InsertManagerPicks(&p)
		},
	)
}

func (h *Handler) handleManagerTransfers(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntryTransfers],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntryTransfers,
		func(t models.EntryTransfersMessage) int { return t.EntryId },
		func(t models.EntryTransfersMessage) error {
			return h.managerRepo.InsertManagerTransfers(&t)
		},
	)
}

func (h *Handler) handleManagerHistory(ctx context.Context) {
	batchProcess(
		ctx,
		h.consumers[h.kafkaConfig.TopicsName.FplEntryHistory],
		h.config.BatchSize,
		h.config.FlushInterval,
		h.kafkaConfig.TopicsName.FplEntryHistory,
		func(hi models.EntryHistoryMessage) int { return hi.EntryId },
		func(hi models.EntryHistoryMessage) error {
			return h.managerRepo.InsertManagerFullHistory(&hi)
		},
	)
}
