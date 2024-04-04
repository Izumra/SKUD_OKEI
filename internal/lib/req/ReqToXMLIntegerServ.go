package req

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"log"
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

	//log.Println(string(body))

	data, _ := io.ReadAll(resp.Body)
	// stream, _ := io.ReadAll(resp.Body)
	// log.Println(string(stream))

	//decoder := xml.NewDecoder(data)
	err = xml.Unmarshal(data, expBody)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		log.Println(string(data))
		return err
	}

	return nil
}
