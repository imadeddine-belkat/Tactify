package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imadeddine-belkat/indexer-service/internal/sofascore_repositories"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	_ "github.com/lib/pq"

	"github.com/imadeddine-belkat/indexer-service/config"
	"github.com/imadeddine-belkat/indexer-service/internal/db/connection"
	"github.com/imadeddine-belkat/indexer-service/internal/fpl_handler"
	"github.com/imadeddine-belkat/indexer-service/internal/fpl_repositories"
	"github.com/imadeddine-belkat/indexer-service/internal/sofascore_handler"
	"github.com/imadeddine-belkat/shared/fpl_models"
)

func main() {
	// 1. Load Config & DB
	cfg := config.LoadConfig()
	kafkaCfg := &cfg.Kafka

	// Connect to database
	fplDb, err := connection.NewRepository(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.FplDatabase,
		"disable",
	)
	if err != nil {
		log.Fatal("Failed to connect to fpl database:", err)
	}

	sofascoreDb, err := connection.NewRepository(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.SofascoreDatabase,
		"disable",
	)
	if err != nil {
		log.Fatal("Failed to connect to sofascore database:", err)
	}

	defer fplDb.Close()
	defer sofascoreDb.Close()

	log.Println("Database connected")

	// 2. Initialize Repositories
	fplFixtureRepo := fpl_repositories.NewFixtureRepo(
		fplDb.DB(),
		&fpl_models.Fixture{},
	)

	FplTeamRepo := fpl_repositories.NewTeamRepo(
		fplDb.DB(),
		&fpl_models.Team{},
	)

	fplPlayerRepo := fpl_repositories.NewPlayerRepo(
		fplDb.DB(),
		&fpl_models.PlayerBootstrap{},
		&fpl_models.PlayerHistory{},
		&fpl_models.PlayerPastHistory{},
	)

	FplManagerRepo := fpl_repositories.NewManagerRepo(
		fplDb.DB(),
		&fpl_models.EntryMessage{},
		&fpl_models.EntryEventPicksMessage{},
		&fpl_models.EntryTransfersMessage{},
		&fpl_models.EntryHistoryMessage{},
	)

	sofascoreTeamReop := sofascore_repositories.NewTeamRepo(
		sofascoreDb.DB(),
		&sofascore_models.StandingMessage{},
		&sofascore_models.TeamOverallStatsMessage{},
		&sofascore_models.MatchStatsMessage{},
	)

	sofacorePlayerRepo := sofascore_repositories.NewPlayerRepo(
		sofascoreDb.DB(),
		&sofascore_models.PlayerMessage{},
	)

	sofascoreMatchRepo := sofascore_repositories.NewMatchRepo(
		sofascoreDb.DB(),
		&sofascore_models.Event{},
	)

	sofascoreLeagueRepo := sofascore_repositories.NewLeagueRepo(
		sofascoreDb.DB(),
		&sofascore_models.LeagueUniqueTournaments{},
	)

	// 3. Initialize Handler
	FplHandler := fpl_handler.NewHandler(
		cfg,
		kafkaCfg,
		fplPlayerRepo,
		FplTeamRepo,
		fplFixtureRepo,
		FplManagerRepo,
	)

	sofascoreHandler := sofascore_handler.NewHandler(
		cfg,
		kafkaCfg,
		sofascoreTeamReop,
		sofacorePlayerRepo,
		sofascoreMatchRepo,
		sofascoreLeagueRepo,
	)

	// 4. Context that listens for OS signals (SIGINT, SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 5. Start Consumers in Background Goroutines
	// The Route method uses 'go fn(ctx)' internally, so these are non-blocking calls
	log.Println("🚀 Starting Kafka consumers...")

	// FPL Topics
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplFixtures.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplTeams.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplPlayersBootstrap.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplPlayersStats.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplPlayerMatchStats.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplPlayerHistoryStats.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplLiveEvent.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplEntry.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplEntryPicks.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplEntryTransfers.Name)
	FplHandler.Route(ctx, kafkaCfg.TopicsName.FplEntryHistory.Name)

	// Sofascore Topics
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascoreLeagueStandings.Name)
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascoreTeamOverallStats.Name)
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascoreTeamMatchStats.Name)
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascorePlayerInfo.Name)
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascoreLeagueRoundMatches.Name)
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascoreLeagueIDs.Name)
	sofascoreHandler.Route(ctx, kafkaCfg.TopicsName.SofascoreLeagueSeasons.Name)

	log.Println("✅ All Kafka consumers started.")

	// 6. BLOCK the main thread until an OS signal is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This line blocks forever until you press Ctrl+C or K8s stops the pod
	sig := <-quit

	log.Printf("⚠️  Signal %v received. Shutting down gracefully...", sig)

	// 7. Cancel context to stop consumers
	cancel()

	time.Sleep(2 * time.Second)
	log.Println("🏁 Shutdown complete.")
}
