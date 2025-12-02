package sofascore_tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/imadeddine-belkat/indexer-service/config"
	"github.com/imadeddine-belkat/indexer-service/internal/sofascore_handler"
	"github.com/imadeddine-belkat/indexer-service/internal/sofascore_repositories"
	"github.com/imadeddine-belkat/shared/sofascore_models"
	_ "github.com/lib/pq"
)

func TestTeamRepo(t *testing.T) {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.SofascoreDatabase,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✅ Database connected")

	// Initialize ONLY fixture repo
	teamRepo := sofascore_repositories.NewTeamRepo(db,
		&sofascore_models.StandingMessage{},
		&sofascore_models.TeamOverallStatsMessage{},
		&sofascore_models.MatchStatsMessage{},
	)

	// Initialize sofascore_handler with nil for other repos
	h := sofascore_handler.NewHandler(
		cfg,
		&cfg.Kafka,
		teamRepo, // teamRepo
		nil,
		nil,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	leagueStandingTopic := cfg.Kafka.TopicsName.SofascoreLeagueStandings
	overallStatsTopic := cfg.Kafka.TopicsName.SofascoreTeamOverallStats
	matchStatsTopic := cfg.Kafka.TopicsName.SofascoreTeamMatchStats

	h.Route(ctx, leagueStandingTopic)
	h.Route(ctx, overallStatsTopic)
	h.Route(ctx, matchStatsTopic)

	log.Printf("✅ Sofascore Team indexer started, listening for %s...", leagueStandingTopic)
	log.Printf("✅ Sofascore Team indexer started, listening for %s...", overallStatsTopic)
	log.Printf("✅ Sofascore Team indexer started, listening for %s...", matchStatsTopic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Stopping...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("✅ Stopped")
}
