package req

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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

	reader, err := io.ReadAll(resp.Body)
	if err == io.EOF {
		resp.Body.Close()
	}

	if resp.StatusCode >= 400 {
		return errors.New(string(reader))
	}
	if err := json.Unmarshal(reader, &expBody); err != nil {
		return err
	}

	return nil
}
