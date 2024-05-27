package req

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
)

var IntegrServiceUtilExitERRChan = make(chan error)

var (
	ErrOrionConnect = errors.New("Орион отвалился")
)

func ReqToXMLIntegerServ(ctx context.Context, method string, url string, headers map[string]string, body []byte, expBody *integrserv.EnvelopeResp) error {
	buffer := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(ctx, method, url, buffer)
	if err != nil {
		return err
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errDescription := err.Error()
		if strings.HasSuffix(errDescription, ": EOF") || strings.HasSuffix(err.Error(), "No connection could be made because the target machine actively refused it.") {
			IntegrServiceUtilExitERRChan <- ErrOrionConnect
			return ErrOrionConnect
		}
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var errOrion integrserv.Error
	envelopeResp := integrserv.ErrEnvelopeResp{
		OperationResult: integrserv.OperationResultServiceError{
			Result: &errOrion,
		},
	}
	checkData := make([]byte, len(data))
	checkData = append(checkData, data...)
	if err := xml.Unmarshal(checkData, &envelopeResp); err == nil {
		return fmt.Errorf(errOrion.InnerExceptionMessage)
	}

	err = xml.Unmarshal(data, expBody)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return err
	}

	return nil
}
