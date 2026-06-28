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

	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/handlers"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/logger"
	"github.com/identicalaffiliation/booking-service/booking/internal/adapters/storage/psql"
	"github.com/identicalaffiliation/booking-service/booking/internal/application"
	"github.com/identicalaffiliation/booking-service/booking/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	ctx := context.Background()

	slogger, err := logger.NewLogger(cfg)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to load logger", err)
		os.Exit(1)
	}

	pool, err, cleanup := setupPool(ctx, cfg)
	if err != nil {
		slogger.Error("failed to setup pgx pool", "error", err)
	}

	defer cleanup()

	roomsRepo := psql.NewRoomsRepository(pool)
	schedulesRepo := psql.NewScheduleRepository(pool)

	roomsUsecase := application.NewRoomsUsecase(roomsRepo, slogger)
	schedulesUsecase := application.NewSchedulesUsecase(schedulesRepo, slogger)

	srv := setupServer(cfg, roomsUsecase, schedulesUsecase)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slogger.Debug("starting server..")
		if err := srv.Start(srv.Server.Addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogger.Error("failed to close server conn", "error", err)
		}
	}()

	<-signals

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slogger.Error("failed to shutdown server", "error", err)
	}

	slogger.Debug("server is stopped gracefully")
}

func setupPool(ctx context.Context, cfg *config.BookingConfig) (*pgxpool.Pool, error, func()) {
	pool, err := pgxpool.New(ctx, cfg.DB_URL)
	if err != nil {
		return nil, fmt.Errorf("open new pgx pool: %w", err), func() {}
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping pgx pool: %w", err), func() {}
	}

	pool.Config().MaxConns = cfg.MaxConns
	pool.Config().MaxConnLifetime = cfg.MaxLifetime

	return pool, nil, func() {
		pool.Close()
	}
}

func setupServer(cfg *config.BookingConfig, ru *application.RoomsUsecase, su *application.SchedulesUsecase) *echo.Echo {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IddleTimeout

	e.POST("/api/v1/rooms", handlers.CreateRoom(ru))
	e.POST("/api/v1/rooms/:roomId/schedule", handlers.CreateSchedule(su))

	return e
}
