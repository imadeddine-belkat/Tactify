package tactify_kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	kafkaConfig "github.com/imadeddine-belkat/tactify-kafka/config"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
			MaxAttempts:  3,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	if topic == "" {
		return fmt.Errorf("kafka topic is empty")
	}
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

func (p *Producer) PublishWithProcess(ctx context.Context, model any, topic string, key []byte) error {
	marshaler := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}

	m, ok := model.(proto.Message)
	if !ok {
		return fmt.Errorf("model does not implement proto.Message")
	}

	out, err := marshaler.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal model: %v for topic: %s", err, topic)
	}

	err = p.Publish(ctx, topic, key, out)
	if err != nil {
		log.Printf("KAFKA ERROR: failed to publish: %v", err)
		return err
	}

	// This helps verify progress in the console
	fmt.Printf("Successfully buffered message for topic %s with key %s\n", topic, string(key))
	return nil
}
