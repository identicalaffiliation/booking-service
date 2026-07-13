package worker

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/booking-service/notifications/internal/adapters/broker"
	"github.com/identicalaffiliation/booking-service/notifications/internal/controller"
	"github.com/identicalaffiliation/booking-service/notifications/internal/ports"
)

type Consumer struct {
	reader  *broker.KafkaReader
	handler *controller.BookingsHandler
	log     ports.Logger
}

func NewConsumer(r *broker.KafkaReader, h *controller.BookingsHandler, l ports.Logger) *Consumer {
	return &Consumer{
		reader:  r,
		handler: h,
		log:     l,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.log.Error("failed to read kafka message", "error", err)
			continue
		}

		messageID, err := uuid.Parse(string(message.Key))
		if err != nil {
			c.log.Error("failed to parse message key", "error", err)
			continue
		}

		if err := c.handler.HandleBookingMessage(ctx, message.Value); err != nil {
			c.log.Error("failed to handle message", "messageId", messageID, "error", err)
		}
	}
}
