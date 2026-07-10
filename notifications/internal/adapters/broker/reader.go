package broker

import (
	"fmt"
	"strings"

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
			Brokers: brokers,
			Topic:   strings.Join(cfg.Topics, ","),
		}),
	}
}
