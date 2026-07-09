package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	kafkaWriter *kafka.Writer
}

func NewKafkaWriter(cfg *config.BookingConfig) *Writer {
	addr := fmt.Sprintf("%s:%d", cfg.KafkaConfig.Host, cfg.KafkaConfig.Port)
	return &Writer{
		kafkaWriter: &kafka.Writer{
			Addr:         kafka.TCP(addr),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Compression:  kafka.Snappy,
		},
	}
}

func (w *Writer) send(ctx context.Context, topic string, key, value []byte) error {
	return w.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (w *Writer) SendJSON(ctx context.Context, topic, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal value: %w", err)
	}

	return w.send(ctx, topic, []byte(key), data)
}

func (w *Writer) Close() error {
	return w.kafkaWriter.Close()
}
