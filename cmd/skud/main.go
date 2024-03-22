package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Izumra/SKUD_OKEI/internal/app"
	"github.com/Izumra/SKUD_OKEI/lib/config"
	"github.com/Izumra/SKUD_OKEI/lib/logger"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.New(logger.Local, os.Stdout)

	ctx := context.Background()
	//TODO: implement connect to db

	server := app.NewServer(logger)
	server.Launch(ctx, cfg.Server.Port)

	chanExit := make(chan os.Signal, 1)
	signal.Notify(chanExit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	signal := <-chanExit
	logger.Info("SKUD system was shutting down", slog.String("signal", signal.String()))
}
