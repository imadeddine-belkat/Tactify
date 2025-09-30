package kafka

import (
	"context"
	"time"

	kafkaConfig "github.com/imadbelkat1/kafka/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer() *Producer {
	cfg := kafkaConfig.LoadConfig()

	var acks kafka.RequiredAcks
	switch cfg.KafkaAcks {
	case "0":
		acks = kafka.RequireNone
	case "1":
		acks = kafka.RequireOne
	default:
		acks = kafka.RequireAll
	}

	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.KafkaBroker),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: acks,
			BatchSize:    cfg.KafkaBatchSize,
			BatchTimeout: time.Duration(cfg.KafkaLingerMs) * time.Millisecond,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
		Time:  time.Now(),
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
