//go:build windows
// +build windows

package integrserv

import (
	"errors"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/lib/config"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	errServiceControl = errors.New("ошибка отправки сигнала службе")
)

type integrService struct {
	logger  *slog.Logger
	service *mgr.Service
}

func New(logger *slog.Logger, cfg *config.Config) *integrService {
	manager, err := mgr.Connect()
	if err != nil {
		panic("Возникла ошибка при подключении к диспетчеру служб" + err.Error())
	}

	service, err := manager.OpenService(cfg.IntegerServer.TitleService)
	if err != nil {
		panic("Возникла ошибка при получении информации о службе интегратора ОРИОН" + err.Error())
	}

	return &integrService{
		logger,
		service,
	}
}

func (s *integrService) Reboot() error {
	op := "utils/IntegrServ/orionService.Reboot"
	logger := s.logger.With(slog.String("op", op))

	_, err := s.service.Control(svc.Stop)
	if err != nil {
		logger.Error("Возникла ошибка при остановке службы: %w", err)
		return errServiceControl
	}

	_, err = s.service.Control(svc.Continue)
	if err != nil {
		logger.Error("Возникла ошибка при запуске службы: %w", err)
		return errServiceControl
	}

	return nil
}
