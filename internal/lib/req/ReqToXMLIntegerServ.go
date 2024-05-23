package req

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
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
		if strings.HasSuffix(err.Error(), ": EOF") {
			IntegrServiceUtilExitERRChan <- ErrOrionConnect
			return ErrOrionConnect
		}
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	err = xml.Unmarshal(data, expBody)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return err
	}

	return nil
}
