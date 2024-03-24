package persons

import (
	"bytes"
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
		integrServAddr,
	}
}

// TODO: Figute out how to pass filter array of the params
func (s *Service) GetPersons(
	ctx context.Context,
	sessionId string,
	offset int64,
	count int64,
) (*integrserv.OperationResult, error) {
	op := "internal/services/persons.Service.GetPersons"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		Offset int64
		Count  int64
		Filter []string
	}
	body := Data{
		Offset: offset,
		Count:  count,
		Filter: nil,
	}
	result, err := s.makeReqToIntegerServ(ctx, "GetPersons", &body)
	if err != nil {
		logger.Info("Occured the error while getting the list of the users", err)
		return nil, err
	}

	return result, nil
}

func (s *Service) GetPersonsCount(
	ctx context.Context,
	sessionId string,
) (*integrserv.OperationResult, error) {
	op := "internal/services/persons.Service.GetPersonsCount"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct{}
	body := Data{}
	result, err := s.makeReqToIntegerServ(ctx, "GetPersonsCount", &body)
	if err != nil {
		logger.Info("Occured the error while counts the quantity of the users", err)
		return nil, err
	}

	return result, nil
}

func (s *Service) GetPersonById(
	ctx context.Context,
	sessionId string,
	id int64,
) (*integrserv.OperationResult, error) {
	op := "internal/services/persons.Service.GetPersonById"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		Id int64
	}
	body := Data{
		Id: id,
	}
	result, err := s.makeReqToIntegerServ(ctx, "GetPersonById", &body)
	if err != nil {
		logger.Info("Occured the error while finding the user by id", err)
		return nil, err
	}

	return result, nil
}

func (s *Service) AddPerson(
	ctx context.Context,
	sessionId string,
	data integrserv.PersonData,
) (*integrserv.OperationResult, error) {
	op := "internal/services/persons.Service.AddPerson"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		PersonData integrserv.PersonData
	}
	body := Data{
		PersonData: data,
	}
	result, err := s.makeReqToIntegerServ(ctx, "AddPerson", &body)
	if err != nil {
		logger.Info("Occured the error while additing the new person", err)
		return nil, err
	}

	return result, nil
}

func (s *Service) UpdatePerson(
	ctx context.Context,
	sessionId string,
	data integrserv.PersonData,
) (*integrserv.OperationResult, error) {
	op := "internal/services/persons.Service.UpdatePerson"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		PersonData integrserv.PersonData
	}
	body := Data{
		PersonData: data,
	}
	result, err := s.makeReqToIntegerServ(ctx, "UpdatePerson", &body)
	if err != nil {
		logger.Info("Occured the error while updating person data", err)
		return nil, err
	}

	return result, nil
}

func (s *Service) DeletePerson(
	ctx context.Context,
	sessionId string,
	data integrserv.PersonData,
) (*integrserv.OperationResult, error) {
	op := "internal/services/persons.Service.DeletePerson"
	logger := s.logger.With(slog.String("op", op))

	err := s.accessGuardian(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	type Data struct {
		PersonData integrserv.PersonData
	}
	body := Data{
		PersonData: data,
	}
	result, err := s.makeReqToIntegerServ(ctx, "DeletePerson", &body)
	if err != nil {
		logger.Info("Occured the error while deleting the person", err)
		return nil, err
	}

	return result, nil
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

func (s *Service) makeReqToIntegerServ(
	ctx context.Context,
	serverMethod string,
	data any,
) (*integrserv.OperationResult, error) {
	headers := map[string]string{
		"Content-Type": "text/xml; charset=utf-8",
		"SOAPAction":   fmt.Sprintf("urn:OrionProIntf-IOrionPro#%s", serverMethod),
	}

	var body bytes.Buffer
	encoder := xml.NewEncoder(&body)
	encoder.Encode(&data)

	var expectedBody integrserv.OperationResult

	err := req.WithRawBody(
		ctx,
		"POST",
		fmt.Sprintf("%s/soap/IOrionPro", s.integrServAddr),
		headers,
		body.Bytes(),
		&expectedBody,
	)

	return &expectedBody, err
}
