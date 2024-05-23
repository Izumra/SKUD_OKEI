package req

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
)

func PreparedReqToXMLIntegerServ(
	ctx context.Context,
	serverMethod string,
	serverAddres string,
	data any,
	respBody any,
) error {
	headers := map[string]string{
		"Content-Type": "text/xml; charset=utf-8",
		"SOAPAction":   fmt.Sprintf("urn:OrionProIntf-IOrionPro#%s", serverMethod),
	}

	preparedBody := integrserv.EnvelopeReq{
		XmlsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		BodySoap: integrserv.BodyReq{
			Data: data,
		},
	}

	body, err := xml.Marshal(preparedBody)
	if err != nil {
		return err
	}

	envelope := integrserv.EnvelopeResp{
		XmlsSoapEnv:     "http://schemas.xmlsoap.org/soap/envelope/",
		XmlsXsd:         "http://www.w3.org/2001/XMLSchema",
		XmlsXsi:         "http://www.w3.org/2001/XMLSchema-instance",
		XmlsSoapEnc:     "http://schemas.xmlsoap.org/soap/encoding/",
		OperationResult: respBody,
	}

	err = ReqToXMLIntegerServ(
		ctx,
		"POST",
		serverAddres,
		headers,
		body,
		&envelope,
	)
	if m, ok := respBody.(map[string]interface{}); ok {
		log.Println(m)
		if v, ok := m["InnerExceptionMessage"].(string); ok {
			return fmt.Errorf("Ошибка при запросе к серверу xml-rpc - %s", v)
		}
	}

	return err
}
