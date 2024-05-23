//go:build windows
// +build windows

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Izumra/SKUD_OKEI/internal/app"
	"github.com/Izumra/SKUD_OKEI/internal/lib/req"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/Izumra/SKUD_OKEI/internal/services/events"
	"github.com/Izumra/SKUD_OKEI/internal/services/key"
	"github.com/Izumra/SKUD_OKEI/internal/services/persons"
	"github.com/Izumra/SKUD_OKEI/internal/storage/cache/embedded"
	"github.com/Izumra/SKUD_OKEI/internal/storage/main/sqlite"
	"github.com/Izumra/SKUD_OKEI/lib/config"
	"github.com/Izumra/SKUD_OKEI/lib/logger"
	integrServUtil "github.com/Izumra/SKUD_OKEI/utils/integerserv"
)

func main() {
	//cfg := config.MustLoad()
	cfg := config.MustLoadByPath("./config/local.yaml")

	logger := logger.New(logger.Local, os.Stdout)

	ctx := context.Background()

	db := sqlite.NewConnetion(cfg)

	integrServiceUtil := integrServUtil.New(logger, cfg)
	go func() {
		for {
			<-req.IntegrServiceUtilExitERRChan
			integrServiceUtil.Reboot(context.Background())
		}
	}()

	sessStore := embedded.NewSessStore()

	authService := auth.NewService(logger, sessStore, db, db)
	eventsService := events.NewService(logger, sessStore, cfg.Server.IntegerServAddr)
	cardService := key.NewService(logger, sessStore, eventsService, cfg.Server.IntegerServAddr)
	personsService := persons.NewService(logger, eventsService, sessStore, cfg.Server.IntegerServAddr)

	services := app.Services{
		AuthService:    authService,
		EventsService:  eventsService,
		PersonsService: personsService,
		CardService:    cardService,
	}

	server := app.NewServer(logger, sessStore, &services)
	server.Launch(ctx, cfg.Server.Port)

	chanExit := make(chan os.Signal, 1)
	signal.Notify(chanExit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	signal := <-chanExit
	logger.Info("SKUD system was shutting down", slog.String("signal", signal.String()))
}
