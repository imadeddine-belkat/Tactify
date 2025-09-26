package kafka

import (
	"context"
	"tactify/kafka/config"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.Config, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{cfg.KafkaBroker},
			Topic:   topic,
		}),
	}
}

func (c *Consumer) Subscribe(ctx context.Context) (<-chan kafka.Message, <-chan error) {
	messages := make(chan kafka.Message, 10000) // Large buffer
	errors := make(chan error, 100)

	go func() {
		defer close(messages)
		defer close(errors)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					select {
					case errors <- err:
					default: // Drop error if channel full
					}
					continue
				}

				select {
				case messages <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return messages, errors
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
