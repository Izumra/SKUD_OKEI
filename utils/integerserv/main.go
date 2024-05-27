//go:build windows
// +build windows

package integrserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	errServiceControl = errors.New("ошибка отправки сигнала службе")
)

type Reboot func(ctx context.Context) error

func RebootManager(titleService string, delayTry time.Duration) Reboot {
	var lock sync.Mutex

	return func(ctx context.Context) error {
		lock.Lock()
		defer lock.Unlock()

		manager, err := mgr.Connect()
		if err != nil {
			return fmt.Errorf("Возникла ошибка при подключении к диспетчеру служб" + err.Error())
		}
		defer manager.Disconnect()

		service, err := manager.OpenService(titleService)
		if err != nil {
			return fmt.Errorf("Возникла ошибка при получении информации о службе интегратора ОРИОН: %v", err.Error())
		}
		defer service.Close()

		status, err := service.Control(svc.Stop)
		if err != nil && status.State != svc.Stopped {
			return fmt.Errorf("Ошибка отправки сигнала службе: " + err.Error())
		}

		ctx, cancel := context.WithTimeout(ctx, delayTry)
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

		ctx, cancel = context.WithTimeout(ctx, delayTry)
		defer ticker.Stop()
		defer cancel()

		if err = service.Start(); err != nil {
			return err
		}

		status, err = service.Query()
		if err != nil {
			return fmt.Errorf("Возникла ошибка при запуске службы: %w", err)
		}

		for status.State != svc.Running {
			select {
			case <-ctx.Done():
				status, err := service.Query()
				if err != nil {
					return err
				}
				return fmt.Errorf("Сервис не запустился за установленное время, текущее состояние %v", status)
			case <-ticker.C:
				continue
			}
		}

		return nil
	}
}
