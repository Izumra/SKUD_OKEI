package admin

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/reqs"
	valueobject "github.com/Izumra/SKUD_OKEI/domain/value-object"
	"github.com/Izumra/SKUD_OKEI/internal/lib/req"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/Izumra/SKUD_OKEI/internal/storage/cache"
)

var (
	ErrSessionTokenInvalid = errors.New("сессия пользователя не действительна")
	ErrAccessDenied        = errors.New("вам отказано в доступе")
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

func (s *Service) GetEvents(ctx context.Context, sessionId string, eventsFilter *reqs.EventFilter) ([]*integrserv.Event, error) {
	op := "internal/services/monitor.Service.GetEvents"
	logger := s.logger.With(slog.String("op", op))

	user, err := s.sessStore.GetByID(ctx, sessionId)
	if err != nil {
		if errors.Is(err, cache.ErrSessionNotFound) {
			return nil, ErrSessionTokenInvalid
		}
		return nil, err
	}

	if user.Role == valueobject.StudentRole {
		return nil, ErrAccessDenied
	}

	//TODO: making the req structure for the event filter
	var expBody []*integrserv.Event
	respBody := &integrserv.OperationResultEvents{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &expBody,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "GetEvents", s.integrServAddr, eventsFilter, respBody)
	if err != nil {
		logger.Info("Occured the error while taking events by filter", err)
		return nil, err
	}
	return expBody, nil
}
func (s *Service) GetEventsCount(ctx context.Context, sessionId string, eventsFilter *reqs.EventFilter) (int64, error) {
	op := "internal/services/monitor.Service.GetEventsCount"
	logger := s.logger.With(slog.String("op", op))

	user, err := s.sessStore.GetByID(ctx, sessionId)
	if err != nil {
		if errors.Is(err, cache.ErrSessionNotFound) {
			return -1, ErrSessionTokenInvalid
		}
		return -1, err
	}

	if user.Role == valueobject.StudentRole {
		return -1, ErrAccessDenied
	}

	type RespData struct {
		OperationResult int64
	}
	var resp RespData
	respBody := &integrserv.OperationResultInt{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &resp,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "GetEventsCount", s.integrServAddr, eventsFilter, respBody)
	if err != nil {
		logger.Info("Occured the error while taking the count of events by filter", err)
		return -1, err
	}
	return resp.OperationResult, nil
}
