package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/identicalaffiliation/booking-service/notifications/config"
	"github.com/identicalaffiliation/booking-service/notifications/internal/adapters/broker"
	"github.com/identicalaffiliation/booking-service/notifications/internal/adapters/psql"
	"github.com/identicalaffiliation/booking-service/notifications/internal/controller"
	"github.com/identicalaffiliation/booking-service/notifications/pkg/logger"
	"github.com/identicalaffiliation/booking-service/notifications/pkg/pool"
	"github.com/identicalaffiliation/booking-service/notifications/pkg/worker"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	slogger, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to load logger", err)
		os.Exit(1)
	}

	ctx := context.Background()

	postgresPool, err, cleanup := pool.SetupPool(ctx, cfg)
	if err != nil {
		slogger.Error("failed to setup postgres pool", "error", err)
		os.Exit(1)
	}

	defer cleanup()

	repo := psql.NewNotificationsRepository(postgresPool)
	reader := broker.NewReader(cfg)

	defer func(reader *broker.KafkaReader) {
		err := reader.Close()
		if err != nil {
			slogger.Error("failed to close reader", "error", err)
		}
	}(reader)

	handler := controller.NewBookingsHandler(repo)
	consumer := worker.NewConsumer(reader, handler, slogger)
	
	consumer.Start(ctx)
}
