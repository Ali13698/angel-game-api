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

// NewKafkaProducer یک اتصال امن به Kafka ایجاد می‌کند
func NewKafkaProducer() *KafkaProducer {
	addr := os.Getenv("KAFKA_ADDR")
	if addr == "" {
		addr = "kafka-kub-service:9092" // آدرس پیش‌فرض از پنل تو
	}
	user := os.Getenv("KAFKA_USER") // مثلا admin
	pass := os.Getenv("KAFKA_PASS") // پسورد از سکرت

	w := &kafka.Writer{
		Addr:         kafka.TCP(addr),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		BatchTimeout: 10 * time.Millisecond,
	}

	// اگر پسورد تنظیم شده بود، احراز هویت SASL را فعال کن
	if pass != "" {
		mechanism := plain.Mechanism{
			Username: user,
			Password: pass,
		}
		w.Transport = &kafka.Transport{
			SASL: mechanism,
		}
	}

	return &KafkaProducer{writer: w}
}

// Publish رویداد را به Kafka ارسال می‌کند
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
