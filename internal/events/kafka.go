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
	pass := os.Getenv("KAFKA") // اینجا اسم متغیر رو دقیقاً KAFKA گذاشتم

	w := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		BatchTimeout: 10 * time.Millisecond,
	}

	// اگر پسورد تنظیم شده بود، احراز هویت SASL را فعال کن
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
	return k.writer.Close()
}
