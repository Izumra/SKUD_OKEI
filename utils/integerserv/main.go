//go:build windows
// +build windows

package integrserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/Izumra/SKUD_OKEI/lib/config"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	errServiceControl = errors.New("ошибка отправки сигнала службе")
)

type integrService struct {
	logger       *slog.Logger
	manager      *mgr.Mgr
	titleService string
}

func New(logger *slog.Logger, cfg *config.Config) *integrService {
	manager, err := mgr.Connect()
	if err != nil {
		panic("Возникла ошибка при подключении к диспетчеру служб" + err.Error())
	}

	return &integrService{
		logger,
		manager,
		cfg.IntegerServer.TitleService,
	}
}

func (s *integrService) Reboot(ctx context.Context) error {

	service, err := s.manager.OpenService(s.titleService)
	if err != nil {
		return fmt.Errorf("Возникла ошибка при получении информации о службе интегратора ОРИОН: %v", err.Error())
	}
	defer service.Close()

	status, err := service.Control(svc.Stop)
	if err != nil && status.State != svc.Stopped {
		s.logger.Error("Возникла ошибка при остановке службы: %w", err)
		return errServiceControl
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer cancel()

	for status.State != svc.Stopped {
		select {
		case <-ctx.Done():
			status, err := service.Query()
			if err != nil {
				return err
			}
			return fmt.Errorf("Сервис не остановился за 10 секунд, текущее состояние %v", status)
		case <-ticker.C:
			continue
		}
	}

	log.Println("Служба IntegrServ остановлена")
	cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer ticker.Stop()
	defer cancel()

	status, err = service.Control(svc.Continue)
	if err != nil {
		s.logger.Error("Возникла ошибка при запуске службы: %w", err)
		return errServiceControl
	}

	for status.State != svc.Running {
		select {
		case <-ctx.Done():
			status, err := service.Query()
			if err != nil {
				return err
			}
			return fmt.Errorf("Сервис не запустился за 10 секунд, текущее состояние %v", status)
		case <-ticker.C:
			continue
		}
	}

	return nil
}
