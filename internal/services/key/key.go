package key

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	valueobject "github.com/Izumra/SKUD_OKEI/domain/value-object"
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers"
	"github.com/Izumra/SKUD_OKEI/internal/lib/req"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/Izumra/SKUD_OKEI/internal/storage/cache"
)

var (
	ErrSessionTokenInvalid = errors.New("сессия пользователя не действительна")
	ErrGettingStats        = errors.New("неожиданная ошибка при загрузке статистики пользователя")
	ErrAccessDenied        = errors.New("вам отказано в доступе")
)

type Service struct {
	logger         *slog.Logger
	sessStore      auth.SessionStorage
	eventService   controllers.EventsService
	integrServAddr string
}

func NewService(
	logger *slog.Logger,
	sessStore auth.SessionStorage,
	eventService controllers.EventsService,
	integrServAddr string,
) *Service {
	return &Service{
		logger,
		sessStore,
		eventService,
		fmt.Sprintf("%s/soap/IOrionPro", integrServAddr),
	}
}

func (s *Service) GetKeys(ctx context.Context, sessionId string, offset int64, count int64) ([]*integrserv.KeyData, error) {
	op := "internal/services/key/Service.GetKeys"
	logger := s.logger.With(slog.String("op", op))

	type Data struct {
		XMLName xml.Name
		Offset  int64
		Count   int64
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "GetKeys",
		},
		Offset: offset,
		Count:  count,
	}

	var expBody []*integrserv.KeyData
	respBody := integrserv.OperationResultEvents{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &expBody,
	}

	err := req.PreparedReqToXMLIntegerServ(ctx, "GetKeys", s.integrServAddr, reqData, &respBody)
	if err != nil {
		logger.Info("Occured the error while taking events by filter", err)
		return nil, err
	}

	return expBody, nil
}

func (s *Service) GetKeyData(ctx context.Context, sessionId string, card string) (*integrserv.KeyData, error) {
	op := "internal/services/key/Service.GetKeyData"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		XMLName xml.Name
		CardNo  string
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "GetKeyData",
		},
		CardNo: card,
	}

	var resp integrserv.KeyData
	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &resp,
	}

	err = req.PreparedReqToXMLIntegerServ(ctx, "GetKeyData", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while finding the user by id", err)
		return nil, err
	}

	return &resp, nil
}

func (s *Service) UpdateKeyData(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error) {
	op := "internal/services/key/Service.UpdateKeyData"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type ReqData struct {
		XMLName xml.Name
		KeyData *integrserv.KeyData
	}
	reqData := ReqData{
		XMLName: xml.Name{
			Local: "UpdateKeyData",
		},
		KeyData: keyData,
	}

	var respData integrserv.KeyData
	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &respData,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "UpdateKeyData", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while updating person data", err)
		return nil, err
	}

	return &respData, nil
}

func (s *Service) AddKey(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error) {
	op := "internal/services/key/Service.AddKey"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	keyData.CodeType = 4
	keyData.AccessLevelId = 3
	keyData.IsBlocked = false
	keyData.IsStoreInDevice = true
	keyData.IsStoreInS2000 = false
	keyData.IsInStopList = false

	keyData.StartDate = time.Now()
	keyData.EndDate = time.Date(
		keyData.StartDate.Year()+4,
		keyData.StartDate.Month(),
		keyData.StartDate.Day(),
		keyData.StartDate.Hour(),
		keyData.StartDate.Minute(),
		keyData.StartDate.Second(),
		keyData.StartDate.Nanosecond(),
		keyData.StartDate.Location(),
	)

	type ReqData struct {
		XMLName xml.Name
		KeyData *integrserv.KeyData
	}
	reqData := ReqData{
		XMLName: xml.Name{
			Local: "AddKey",
		},
		KeyData: keyData,
	}

	var respData integrserv.KeyData
	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &respData,
	}

	err = req.PreparedReqToXMLIntegerServ(ctx, "AddKey", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while updating person data", err)
		return nil, err
	}

	return &respData, nil
}

func (s *Service) ReadKeyCode(ctx context.Context, sessionId string, idReader int) (string, error) {
	op := "internal/services/key/Service.ReadKeyCode"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return "", err
	}

	type Reader struct {
		AccessPointID int
		Passmode      int
	}
	readers := []Reader{
		{
			AccessPointID: 1,
			Passmode:      1,
		},
		{
			AccessPointID: 1,
			Passmode:      2,
		},
		{
			AccessPointID: 2,
			Passmode:      1,
		},
		{
			AccessPointID: 2,
			Passmode:      2,
		},
	}

	if idReader < 1 || idReader > 5 {
		return "", fmt.Errorf("Неверный номер считывателя")
	}

	selectedReader := readers[idReader-1]

	timeSurvey := time.Now()

	filter := integrserv.EventFilter{
		XMLName: xml.Name{
			Local: "GetEvents",
		},
		BeginTime: time.Date(timeSurvey.Year(), timeSurvey.Month(), timeSurvey.Day(), timeSurvey.Hour(), timeSurvey.Minute()-5, 0, 0, timeSurvey.Location()),
		EndTime:   timeSurvey,
	}

	events, err := s.eventService.GetEvents(ctx, &filter)
	if err != nil {
		logger.Info("Occured the error while reading the card", err)
		return "", err
	}

	var keys []string
	regExp := regexp.MustCompile(`.,  (.*) Считыватель$`)
	for i := range events {
		if regExp.MatchString(events[i].Description) {
			if events[i].PassMode == selectedReader.Passmode && events[i].AccessPointId == selectedReader.AccessPointID {
				submatches := regExp.FindStringSubmatch(events[i].Description)
				keys = append(keys, submatches[1])
			}
		}
	}

	if len(keys) != 0 {
		return keys[len(keys)-1], nil
	}

	return "", fmt.Errorf("Считать ключ не удалось, попробуйте еще раз")
}

func (s *Service) ConvertWiegandToTouchMemory(ctx context.Context, sessionId string, code int, codeSize int) (string, error) {
	op := "internal/services/key/Service.ConvertWiegandToTouchMemory"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return "", err
	}

	type ReqData struct {
		XMLName  xml.Name
		Code     int
		CodeSize int
	}
	reqData := ReqData{
		XMLName: xml.Name{
			Local: "ConvertWiegandToTouchMemory",
		},
		Code:     code,
		CodeSize: codeSize,
	}

	var respData integrserv.Result
	respBody := &integrserv.OperationResultString{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &respData,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "ConvertWiegandToTouchMemory", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while updating person data", err)
		return "", err
	}

	return respData.OperationResult, nil
}
func (s *Service) ConvertPinToTouchMemory(ctx context.Context, sessionId string, pin string) (string, error) {
	op := "internal/services/key/Service.ConvertPinToTouchMemory"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return "", err
	}

	type ReqData struct {
		XMLName xml.Name
		Pin     string
	}
	reqData := ReqData{
		XMLName: xml.Name{
			Local: "ConvertPinToTouchMemory",
		},
		Pin: pin,
	}

	var respData integrserv.Result
	respBody := &integrserv.OperationResultString{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &respData,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "ConvertPinToTouchMemory", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while updating person data", err)
		return "", err
	}

	return respData.OperationResult, nil
}

func (s *Service) accessGuardian(ctx context.Context, sessionId string) error {
	user, err := s.sessStore.GetByID(ctx, sessionId)
	if err != nil {
		if errors.Is(err, cache.ErrSessionNotFound) {
			return ErrSessionTokenInvalid
		}
		return err
	}

	if user.Role == valueobject.StudentRole {
		return ErrAccessDenied
	}

	return nil
}
