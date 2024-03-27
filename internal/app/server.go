package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Izumra/SKUD_OKEI/internal/http"
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/fiber/v2"
)

type Services struct {
	AuthService    controllers.AuthService
	PersonsService controllers.PersonsService
	EventsService  controllers.EventsService
}

type Server struct {
	app      *fiber.App
	logger   *slog.Logger
	services *Services
}

func NewServer(
	logger *slog.Logger,
	sessionStorage auth.SessionStorage,
	services *Services,
) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout: time.Second * 10,
	})

	http.RegistrHandlers(
		app,
		sessionStorage,
		services.AuthService,
		services.PersonsService,
		services.EventsService,
	)

	return &Server{
		app,
		logger,
		services,
	}
}

func (s *Server) Launch(ctx context.Context, port int) {
	op := "app.Server.Launch"
	logger := s.logger.With(slog.String("op", op))

	addr := fmt.Sprintf(":%d", port)

	err := s.app.Listen(addr)
	if err != nil {
		logger.Error("Occured the error while launching the server", slog.String("addr", addr))
		os.Exit(1)
	}
}

func (s *Server) Stop() {
	op := "app.Server.Stop"
	logger := s.logger.With(slog.String("op", op))

	err := s.app.Shutdown()
	if err != nil {
		logger.Error("Occured the error while stopping the server")
		os.Exit(1)
	}

	logger.Info("Server was gracefully sutting down")
}
