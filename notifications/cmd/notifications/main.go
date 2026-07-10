package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/identicalaffiliation/booking-service/notifications/config"
	"github.com/identicalaffiliation/booking-service/notifications/pkg/logger"
	"github.com/identicalaffiliation/booking-service/notifications/pkg/pool"
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

	postgresPool, err, cleanup := pool.SetupPool(context.Background(), cfg)
	if err != nil {
		slogger.Error("failed to setup postgres pool", "error", err)
		os.Exit(1)
	}

	defer cleanup()

}
