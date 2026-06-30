package httpserver

import (
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/controller"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
	"github.com/labstack/echo/v4"
)

func SetupServer(cfg *config.BookingConfig, ru *usecase.RoomsUsecase, su *usecase.SchedulesUsecase) *echo.Echo {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IddleTimeout

	e.POST("/api/v1/rooms", controller.CreateRoom(ru))
	e.POST("/api/v1/rooms/:roomId/schedule", controller.CreateSchedule(su))

	return e
}
