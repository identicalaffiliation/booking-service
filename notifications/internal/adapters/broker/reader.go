package broker

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/booking-service/notifications/config"
	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	reader *kafka.Reader
}

func NewReader(cfg *config.NotificationsConfig) *KafkaReader {
	brokers := []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}
	return &KafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       "",
			GroupID:     "notifications-service",
			GroupTopics: cfg.Topics,
		}),
	}
}

func (r *KafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return r.reader.ReadMessage(ctx)
}

func (r *KafkaReader) Close() error {
	return r.reader.Close()
}
