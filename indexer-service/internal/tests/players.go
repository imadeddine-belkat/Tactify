package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imadbelkat1/indexer-service/config"
	"github.com/imadbelkat1/indexer-service/internal/fpl_repositories"
	"github.com/imadbelkat1/indexer-service/internal/handler"
	"github.com/imadbelkat1/shared/models"
	_ "github.com/lib/pq"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.FplDatabase,
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

	// Initialize player repo
	playerRepo := fpl_repositories.NewPlayerRepo(
		db,
		&models.PlayerBootstrap{},
		&models.PlayerHistory{},
		&models.PlayerPastHistory{},
	)

	// Initialize handler with player repo only
	h := handler.NewHandler(
		cfg,
		&cfg.Kafka,
		playerRepo, // playerRepo
		nil,        // teamRepo
		nil,        // fixtureRepo
		nil,        // managerRepo
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming all three player topics
	h.Route(ctx, cfg.Kafka.TopicsName.FplPlayersBootstrap)
	h.Route(ctx, cfg.Kafka.TopicsName.FplPlayersStats)
	h.Route(ctx, cfg.Kafka.TopicsName.FplPlayerMatchStats)
	h.Route(ctx, cfg.Kafka.TopicsName.FplPlayerHistoryStats)

	log.Println("✅ Player indexer started, listening for player data...")
	log.Println("   - Player bootstrap topic:", cfg.Kafka.TopicsName.FplPlayersBootstrap)
	log.Println("   - Player stats topic:", cfg.Kafka.TopicsName.FplPlayersStats)
	log.Println("   - Player match stats topic:", cfg.Kafka.TopicsName.FplPlayerMatchStats)
	log.Println("   - Player past history topic:", cfg.Kafka.TopicsName.FplPlayerHistoryStats)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Stopping...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("✅ Stopped")
}
