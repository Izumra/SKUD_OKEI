package req

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

func WithRawBody(ctx context.Context, method string, url string, headers map[string]string, body []byte, expBody any) error {
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

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&expBody)
	if err != nil {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&expBody)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
