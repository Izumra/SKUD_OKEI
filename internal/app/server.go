package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/Izumra/SKUD_OKEI/internal/http"
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers"
	"github.com/gofiber/fiber/v3"
)

type Services struct {
	AuthService    controllers.AuthService
	PersonsService controllers.PersonsService
}

type Server struct {
	app      *fiber.App
	logger   *slog.Logger
	services *Services
}

func NewServer(
	logger *slog.Logger,
	services *Services,
) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout: time.Second * 10,
	})

	http.RegistrHandlers(
		app,
		services.AuthService,
		services.PersonsService,
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
	log.Println(addr)
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
