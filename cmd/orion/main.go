//go:build windows
// +build windows

package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Izumra/SKUD_OKEI/internal/lib/req"
	"github.com/Izumra/SKUD_OKEI/lib/config"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

var (
	errServiceControl = errors.New("ошибка отправки сигнала службе")
)

func main() {
	cfg := config.MustLoadByPath("./config/local.yaml")

	service := newOrionService(cfg)

	var err error
	for {
		time.Sleep(3 * time.Second)
		headers := map[string]string{
			"Content-Type": "text/xml; charset=utf-8",
			"SOAPAction":   "urn:OrionProIntf-IOrionPro#GetReplService",
		}

		err = req.ReqToXMLIntegerServ(
			context.Background(),
			"POST",
			cfg.IntegerServer.Addr,
			headers,
			nil,
			nil,
		)
		if err != nil && errors.Is(err, req.ErrOrionConnect) {
			service.Reboot()
		}
	}
}

type orionService struct {
	service *mgr.Service
}

func newOrionService(cfg *config.Config) *orionService {
	manager, err := mgr.Connect()
	if err != nil {
		panic("Возникла ошибка при подключении к диспетчеру служб" + err.Error())
	}

	service, err := manager.OpenService(cfg.IntegerServer.TitleService)
	if err != nil {
		panic("Возникла ошибка при получении информации о службе интегратора ОРИОН" + err.Error())
	}

	return &orionService{
		service,
	}
}

func (s orionService) Reboot() error {
	_, err := s.service.Control(svc.Stop)
	if err != nil {
		log.Println("Возникла ошибка при остановке службы: %w", err)
		return errServiceControl
	}

	_, err = s.service.Control(svc.Continue)
	if err != nil {
		log.Println("Возникла ошибка при запуске службы: %w", err)
		return errServiceControl
	}

	return nil
}
