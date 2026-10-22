package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer() *KafkaProducer {
	addr := os.Getenv("KAFKA_ADDR")
	user := os.Getenv("KAFKA_USER")
	pass := os.Getenv("KAFKA")

	if addr == "" {
		slog.Warn("KAFKA_ADDR not set, running without Kafka")
		return &KafkaProducer{writer: nil}
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		BatchTimeout: 10 * time.Millisecond,
	}

	if pass != "" {
		w.Transport = &kafka.Transport{
			SASL: plain.Mechanism{
				Username: user,
				Password: pass,
			},
		}
	}

	slog.Info("Kafka producer initialized", "addr", addr)
	return &KafkaProducer{writer: w}
}

func (k *KafkaProducer) Publish(ctx context.Context, topic string, key string, payload any) error {
	if k.writer == nil {
		return nil // Fail-safe: if Kafka is not configured, just ignore
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	err = k.writer.WriteMessages(ctx, msg)
	if err != nil {
		slog.Error("Kafka publish failed", "topic", topic, "error", err)
		return err
	}

	slog.Info("Kafka event published", "topic", topic, "key", key)
	return nil
}

func (k *KafkaProducer) Close() error {
	if k.writer == nil {
		return nil
	}
	return k.writer.Close()
}
