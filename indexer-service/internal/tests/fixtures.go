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

	// Initialize ONLY fixture repo
	fixtureRepo := fpl_repositories.NewFixtureRepo(db, &models.Fixture{})

	// Initialize handler with nil for other repos
	h := handler.NewHandler(
		cfg,
		&cfg.Kafka,
		nil, // playerRepo
		nil, // teamRepo
		fixtureRepo,
		nil, // managerRepo
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming fixtures ONLY
	h.Route(ctx, cfg.Kafka.TopicsName.FplFixtures)

	log.Println("✅ Fixture indexer started, listening for fixtures...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Stopping...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("✅ Stopped")
}
