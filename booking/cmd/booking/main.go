package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
	"github.com/identicalaffiliation/booking-service/booking/pkg/generator"
	"github.com/identicalaffiliation/booking-service/booking/pkg/httpserver"
	"github.com/identicalaffiliation/booking-service/booking/pkg/logger"
	"github.com/identicalaffiliation/booking-service/booking/pkg/pool"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slogger, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to load logger", err)
		os.Exit(1)
	}

	postgresPool, err, cleanup := pool.SetupPool(ctx, cfg)
	if err != nil {
		slogger.Error("failed to setup pgx pool", "error", err)
	}

	defer cleanup()

	roomsRepo := psql.NewRoomsRepository(postgresPool)
	schedulesRepo := psql.NewScheduleRepository(postgresPool)
	slotsRepo := psql.NewSlotsRepository(postgresPool)

	rooms := usecase.NewRoomsUsecase(roomsRepo, slogger)
	slots := usecase.NewSlotsUsecase(slotsRepo, schedulesRepo, slogger, cfg)
	schedules := usecase.NewSchedulesUsecase(schedulesRepo, slogger, slots)

	srv := httpserver.SetupServer(cfg, rooms, schedules)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	generator.StartSlotGenerator(ctx, slots, slogger, cfg)

	go func() {
		slogger.Debug("starting server..")
		if err := srv.Start(srv.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogger.Error("failed to close server conn", "error", err)
		}
	}()

	<-signals

	ctx, cancel = context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slogger.Error("failed to shutdown server", "error", err)
	}

	slogger.Debug("server is stopped gracefully")
}
