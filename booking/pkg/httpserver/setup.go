package httpserver

import (
	"fmt"

	"github.com/identicalaffiliation/booking-service/booking/config"
	"github.com/identicalaffiliation/booking-service/booking/internal/controller"
	middlewares "github.com/identicalaffiliation/booking-service/booking/internal/controller/middleware"
	"github.com/identicalaffiliation/booking-service/booking/internal/domain"
	"github.com/identicalaffiliation/booking-service/booking/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupServer(
	cfg *config.BookingConfig,
	ru *usecase.RoomsUsecase,
	su *usecase.SchedulesUsecase,
	au *usecase.AuthUsecase,
) *echo.Echo {
	e := echo.New()
	e.Server.Addr = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout

	e.HTTPErrorHandler = controller.HTTPErrorHandler()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())

	api := e.Group("/api/v1")
	api.POST("/sign-up", controller.Registration(au))
	api.POST("/sign-in", controller.Login(au))

	private := api.Group("", middlewares.AuthMiddleware(cfg.JwtSecret, cfg.AccessTokenConfig.IssuedBy))

	// user booking routes
	// ...

	admin := private.Group("", middlewares.RoleMiddleware(domain.Admin))

	// admin room routes
	admin.POST("/rooms", controller.CreateRoom(ru))
	admin.GET("/rooms/:roomId", controller.GetRoom(ru))
	admin.GET("/rooms", controller.GetRooms(ru))
	admin.DELETE("/rooms/:roomId", controller.DeleteRoom(ru))

	// admin schedule routes
	admin.POST("/rooms/:roomId/schedule", controller.CreateSchedule(su))

	return e
}
