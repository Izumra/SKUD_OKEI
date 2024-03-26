package req

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
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
		return err
	}
	defer resp.Body.Close()

	// data, _ := xml.Marshal(expBody)
	// log.Println(string(data))

	// stream, _ := io.ReadAll(resp.Body)
	// log.Println(string(stream))

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(expBody)
	if err != nil {
		return err
	}

	return nil
}
