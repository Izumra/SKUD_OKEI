package key

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/internal/lib/req"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
)

//TODO: implement the key service

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

func (s *Service) GetKeys(ctx context.Context, sessionId string, offset int64, count int64) ([]*integrserv.KeyData, error) {
	return nil, nil
}

func (s *Service) GetKeyData(ctx context.Context, sessionId string, card string) (*integrserv.KeyData, error) {
	op := "internal/services/key.Service.GetKeyData"
	logger := s.logger.With(slog.String("op", op))

	type Data struct {
		XMLName xml.Name
		CardNo  string
	}
	reqData := Data{
		XMLName: xml.Name{
			Local: "GetKeyData",
		},
	}

	var expBody []*integrserv.Event
	respBody := integrserv.OperationResultEvents{
		SoapEnvEncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		XmlnsNS1:             "urn:OrionProIntf-IOrionPro",
		XmlnsNS2:             "urn:OrionProIntf",

		Result: &expBody,
	}

	err := req.PreparedReqToXMLIntegerServ(ctx, "GetKeyData", s.integrServAddr, reqData, &respBody)
	if err != nil {
		logger.Info("Occured the error while taking events by filter", err)
		return nil, err
	}

	return nil, nil
}

func (s *Service) UpdateKeyData(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error) {
	return nil, nil
}

func (s *Service) AddKey(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error) {
	return nil, nil
}
