package persons

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/resp"
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
	eventsService  controllers.EventsService
	sessStore      auth.SessionStorage
	integrServAddr string
}

func NewService(
	logger *slog.Logger,
	eventsService controllers.EventsService,
	sessStore auth.SessionStorage,
	integrServAddr string,
) *Service {
	return &Service{
		logger,
		eventsService,
		sessStore,
		fmt.Sprintf("%s/soap/IOrionPro", integrServAddr),
	}
}

// TODO: Figute out how to pass filter array of the params
func (s *Service) GetPersons(
	ctx context.Context,
	sessionId string,
	offset int64,
	count int64,
	filterParams []string,
) ([]*integrserv.PersonData, error) {
	op := "internal/services/persons.Service.GetPersons"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type FilterItem struct {
		XMLName xml.Name `xml:"Filter"`
		Value   string   `xml:"Value"`
	}
	type Data struct {
		XMLName      xml.Name
		WithoutPhoto bool
		Offset       int64
		Count        int64
		Filter       []FilterItem
	}

	filter := make([]FilterItem, len(filterParams))
	for i, v := range filterParams {
		filter[i] = FilterItem{
			Value: v,
		}
	}
	reqData := Data{
		WithoutPhoto: true,
		XMLName: xml.Name{
			Local: "GetPersons",
		},
		Offset: offset,
		Count:  count,
		Filter: filter,
	}

	var expBody []*integrserv.PersonData
	respBody := &integrserv.OperationResultPersons{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &expBody,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "GetPersons", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while getting the list of the users", err)
		return nil, err
	}

	return expBody, nil
}

func (s *Service) GetPersonsCount(
	ctx context.Context,
	sessionId string,
) (int64, error) {
	op := "internal/services/persons.Service.GetPersonsCount"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return -1, err
	}

	type Data struct {
		XMLName xml.Name
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "GetPersonsCount",
		},
	}

	type Count struct {
		OperationResult int64
	}
	var resp Count
	respBody := &integrserv.OperationResultInt{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &resp,
	}

	err = req.PreparedReqToXMLIntegerServ(ctx, "GetPersonsCount", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while counts the quantity of the users", err)
		return -1, err
	}

	return resp.OperationResult, nil
}

func (s *Service) GetPersonById(
	ctx context.Context,
	sessionId string,
	id int64,
) (*integrserv.PersonData, error) {
	op := "internal/services/persons.Service.GetPersonById"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		XMLName xml.Name
		Id      int64
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "GetPersonById",
		},
		Id: id,
	}

	var resp integrserv.PersonData
	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &resp,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "GetPersonById", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while finding the user by id", err)
		return nil, err
	}

	return &resp, nil
}

func (s *Service) AddPerson(
	ctx context.Context,
	sessionId string,
	data integrserv.PersonData,
) (*integrserv.PersonData, error) {
	op := "internal/services/persons.Service.AddPerson"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type ReqData struct {
		XMLName    xml.Name
		PersonData integrserv.PersonData
	}
	reqData := ReqData{
		XMLName: xml.Name{
			Local: "AddPerson",
		},
		PersonData: data,
	}

	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &data,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "AddPerson", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while additing the new person", err)
		return nil, err
	}

	return &data, nil
}

func (s *Service) UpdatePerson(
	ctx context.Context,
	sessionId string,
	data integrserv.PersonData,
) (*integrserv.PersonData, error) {
	op := "internal/services/persons.Service.UpdatePerson"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		XMLName    xml.Name
		PersonData integrserv.PersonData
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "UpdatePerson",
		},
		PersonData: data,
	}

	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &data,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "UpdatePerson", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while updating person data", err)
		return nil, err
	}

	return &data, nil
}

func (s *Service) DeletePerson(
	ctx context.Context,
	sessionId string,
	data integrserv.PersonData,
) (*integrserv.PersonData, error) {
	op := "internal/services/persons.Service.DeletePerson"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		XMLName    xml.Name
		PersonData integrserv.PersonData
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "DeletePerson",
		},
		PersonData: data,
	}

	respBody := &integrserv.OperationResult{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &data,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "DeletePerson", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while deleting the person", err)
		return nil, err
	}

	return &data, nil
}

func (s *Service) GetDepartments(
	ctx context.Context,
	sessionId string,
) ([]*integrserv.Department, error) {
	op := "internal/services/persons.Service.GetDepartments"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		XMLName xml.Name
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "GetDepartments",
		},
	}

	var departments []*integrserv.Department
	respBody := &integrserv.OperationResultDepartments{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &departments,
	}
	err = req.PreparedReqToXMLIntegerServ(ctx, "GetDepartments", s.integrServAddr, reqData, respBody)
	if err != nil {
		logger.Info("Occured the error while deleting the person", err)
		return nil, err
	}

	return departments, nil
}

func (s *Service) BindKey(ctx context.Context, sessionId string) (*integrserv.KeyData, error) {
	return nil, nil
}
func (s *Service) GetKeyData(ctx context.Context, sessionId string) (*integrserv.KeyData, error) {
	return nil, nil
}
func (s *Service) UpdataKeyData(ctx context.Context, sessionId string) (*integrserv.KeyData, error) {
	return nil, nil
}

func (s *Service) GetDaylyUserStats(ctx context.Context, sessionId string, id int64, date time.Time) ([]*resp.Activity, error) {
	op := "internal/services/persons.Service.GetDaylyUserStats"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	var stats sync.WaitGroup
	chanErr := make(chan error)

	response := make([]*resp.Activity, 17)
	for i := 0; i < 17; i++ {
		stats.Add(1)
		hour := i
		go func() {
			defer stats.Done()

			beginTime := time.Date(date.Year(), date.Month(), date.Day(), -(1 + hour), 0, 0, 0, date.Location())
			endTime := time.Date(date.Year(), date.Month(), date.Day(), -(2 + hour), 0, 0, 0, date.Location())

			filter := integrserv.EventFilter{
				XMLName: xml.Name{
					Local: "GetEvents",
				},
				BeginTime: beginTime,
				EndTime:   endTime,
				Persons: integrserv.Persons{
					PersonData: []*integrserv.PersonData{
						{
							Id: id,
						},
					},
				},
			}

			eventsComing, err := s.eventsService.GetEvents(ctx, &filter, 0, 0)
			if err != nil {
				chanErr <- err
				return
			}

			var countComing int
			var countLeaving int
			for _, e := range eventsComing {
				if e.PassMode == 1 {
					countComing++
				} else if e.PassMode == 2 {
					countLeaving++
				}
			}

			response[hour] = &resp.Activity{
				Time:    endTime,
				Coming:  countComing,
				Leaving: countLeaving,
			}
		}()
	}

	go func() {
		stats.Wait()
		slices.Reverse(response)
		close(chanErr)
	}()

	if err := <-chanErr; err != nil {
		logger.Error("Occured the error while requesting for the day stats", err)
		return nil, err
	}

	return response, nil
}

func (s *Service) GetMonthlyUserStats(ctx context.Context, sessionId string, id int64, month time.Time) ([]*resp.Activity, error) {
	op := "internal/services/persons.Service.GetMonthlyUserStats"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	var stats sync.WaitGroup
	chanErr := make(chan error)

	count := time.Date(month.Year(), month.Month(), 0, 0, 0, 0, 0, time.UTC).Day()
	response := make([]*resp.Activity, count)
	for i := 0; i < count; i++ {
		stats.Add(1)
		day := i
		go func() {
			defer stats.Done()

			beginTime := time.Date(month.Year(), month.Month(), -(day + 1), 0, 0, 0, 0, month.Location())
			endTime := time.Date(month.Year(), month.Month(), -(day + 1), 24, 0, 0, 0, month.Location())

			filter := integrserv.EventFilter{
				XMLName: xml.Name{
					Local: "GetEvents",
				},
				BeginTime: beginTime,
				EndTime:   endTime,
				Persons: integrserv.Persons{
					PersonData: []*integrserv.PersonData{
						{
							Id: id,
						},
					},
				},
			}

			eventsComing, err := s.eventsService.GetEvents(ctx, &filter, 0, 0)
			if err != nil {
				chanErr <- err
				return
			}

			var countComing int
			var countLeaving int
			for _, e := range eventsComing {
				if e.PassMode == 1 {
					countComing++
				} else if e.PassMode == 2 {
					countLeaving++
				}
			}

			response[day] = &resp.Activity{
				Time:    endTime,
				Coming:  countComing,
				Leaving: countLeaving,
			}
		}()
	}

	go func() {
		stats.Wait()
		slices.Reverse(response)
		close(chanErr)
	}()

	if err := <-chanErr; err != nil {
		logger.Error("Occured the error while requesting for the day stats", err)
		return nil, err
	}

	return response, nil
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
