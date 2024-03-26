package persons

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
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

func (s *Service) BindKey(ctx context.Context, sessionId string) (*integrserv.KeyData, error) {
	return nil, nil
}
func (s *Service) GetKeyData(ctx context.Context, sessionId string) (*integrserv.KeyData, error) {
	return nil, nil
}
func (s *Service) UpdataKeyData(ctx context.Context, sessionId string) (*integrserv.KeyData, error) {
	return nil, nil
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
