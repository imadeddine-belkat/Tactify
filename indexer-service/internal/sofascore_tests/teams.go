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
	"github.com/imadbelkat1/indexer-service/internal/sofascore_handler"
	"github.com/imadbelkat1/indexer-service/internal/sofascore_repositories"
	"github.com/imadbelkat1/shared/sofascore_models"
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
	teamRepo := sofascore_repositories.NewTeamRepo(db, &sofascore_models.TopTeamsMessage{})

	// Initialize sofascore_handler with nil for other repos
	h := sofascore_handler.NewHandler(
		cfg,
		&cfg.Kafka,
		teamRepo, // teamRepo
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topic := cfg.Kafka.TopicsName.SofascoreTopTeamsStats

	h.Route(ctx, topic)
	log.Printf("✅ Handler goroutine launched for topic: %s", topic)

	log.Printf("✅ Sofascore Team indexer started, listening for %s...", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("🛑 Stopping...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("✅ Stopped")
}
