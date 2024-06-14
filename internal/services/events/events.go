package events

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/internal/lib/req"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
)

type Service struct {
	logger         *slog.Logger
	sessStore      auth.SessionStorage
	integrServAddr string
}

func NewService(
	logger *slog.Logger,
	sessStore auth.SessionStorage,
	integrServAddr string,
) *Service {
	return &Service{
		logger,
		sessStore,
		fmt.Sprintf("%s/soap/IOrionPro", integrServAddr),
	}
}

func (s *Service) GetEvents(ctx context.Context, eventsFilter *integrserv.EventFilter) ([]integrserv.Event, error) {
	op := "internal/services/events.Service.GetEvents"
	logger := s.logger.With(slog.String("op", op))

	var expBody []integrserv.Event
	respBody := integrserv.OperationResultEvents{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &expBody,
	}
	err := req.PreparedReqToXMLIntegerServ(ctx, "GetEvents", s.integrServAddr, eventsFilter, &respBody)
	if err != nil {
		logger.Info("Occured the error while taking events by filter", err)
		return nil, err
	}
	return expBody, nil
}

func (s *Service) GetEventsCount(ctx context.Context, eventsFilter *integrserv.EventCountFilter) (int64, error) {
	op := "internal/services/events.Service.GetEventsCount"
	logger := s.logger.With(slog.String("op", op))

	type Count struct {
		OperationResult int64
	}
	resp := Count{}
	respBody := integrserv.OperationResultInt{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &resp,
	}
	err := req.PreparedReqToXMLIntegerServ(ctx, "GetEventsCount", s.integrServAddr, eventsFilter, &respBody)
	if err != nil {
		logger.Info("Occured the error while taking the count of events by filter", err)
		return -1, err
	}
	return resp.OperationResult, nil
}
