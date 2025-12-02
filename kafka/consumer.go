package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/imadeddine-belkat/kafka/config"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.KafkaConfig, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{cfg.KafkaBroker},
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       1,                      // Changed: Process messages immediately (was 10e3)
			MaxBytes:       10e6,                   // Keep: Max 10MB per fetch
			MaxWait:        100 * time.Millisecond, // Added: Max wait time 100ms
			CommitInterval: time.Second,            // Added: Auto-commit every second
			StartOffset:    kafka.LastOffset,       // Added: Start from latest (or use kafka.FirstOffset for all messages)
		}),
	}
}

func (c *Consumer) Subscribe(ctx context.Context) (<-chan kafka.Message, <-chan error) {
	messages := make(chan kafka.Message, 100)
	errors := make(chan error, 10)

	go func() {
		defer close(messages)
		defer close(errors)

		for {
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // Context cancelled
				}
				select {
				case errors <- err:
				case <-ctx.Done():
					return
				default:
					// Drop error if channel full
				}
				continue
			}

			select {
			case messages <- msg:
				// Auto-commit the message after successful delivery to channel
				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					select {
					case errors <- fmt.Errorf("committing message: %w", err):
					default:
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return messages, errors
}

// CommitMessage commits a single message (kept for backward compatibility)
func (c *Consumer) CommitMessage(ctx context.Context, msg kafka.Message) error {
	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("committing message: %w", err)
	}
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
