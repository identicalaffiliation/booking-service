package cluster

import (
	"fmt"
	"strings"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/segmentio/kafka-go"
)

func CreateKafkaTopics(cfg *config.BookingConfig) error {
	addr := fmt.Sprintf("%s:%d", cfg.KafkaConfig.Host, cfg.KafkaConfig.Port)
	fmt.Println(addr)
	conn, err := kafka.Dial(cfg.ConnectionType, addr)
	if err != nil {
		return fmt.Errorf("dial kafka: %w", err)
	}

	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("get cluster controller: %w", err)
	}

	controllerAddr := fmt.Sprintf("%s:%d", cfg.KafkaConfig.Host, controller.Port)
	controllerConn, err := kafka.Dial(cfg.ConnectionType, controllerAddr)
	if err != nil {
		return fmt.Errorf("dial cluster controller: %w", err)
	}

	defer controllerConn.Close()

	for _, topic := range cfg.Topics {
		err := controllerConn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     cfg.Partitions,
			ReplicationFactor: cfg.Replications,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("create kafka topic: %w", err)
			}
		}
	}

	return nil
}
