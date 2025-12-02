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
)

func TestLeagueRepo(t *testing.T) {
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

	leagueRepo := sofascore_repositories.NewLeagueRepo(
		db,
		&sofascore_models.LeagueUniqueTournaments{},
	)

	h := sofascore_handler.NewHandler(
		cfg,
		&cfg.Kafka,
		nil,
		nil,
		leagueRepo,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	leagueIDsTopic := cfg.Kafka.TopicsName.SofascoreLeagueIDs
	leagueSeasonsTopic := cfg.Kafka.TopicsName.SofascoreLeagueSeasons

	h.Route(ctx, leagueIDsTopic)
	h.Route(ctx, leagueSeasonsTopic)

	log.Printf("✅ Sofascore Match indexer started, listening for %s...", leagueIDsTopic)
	log.Printf("✅ Sofascore Match indexer started, listening for %s...", leagueSeasonsTopic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Stopping...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("✅ Stopped")

}
