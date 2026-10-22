package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/angel/angel/internal/events"
)

func main() {
	fmt.Println("Starting Kafka connection test...")
	kafka := events.NewKafkaProducer()
	defer kafka.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// یک پیام تستی می‌فرستیم
	err := kafka.Publish(ctx, "system.test", "startup", map[string]string{
		"status": "kafka_connected_successfully",
		"time":   time.Now().String(),
	})
	if err != nil {
		log.Fatalf("❌ Failed to publish message: %v", err)
	}

	fmt.Println("✅ Successfully sent message to Kafka!")
}
