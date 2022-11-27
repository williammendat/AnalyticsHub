package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type (
	HttpClient struct {
		httpClient *http.Client
		ctx        context.Context
	}
)

func NewClient() *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{
			Timeout: time.Second * 0,
		},
		ctx: context.Background(),
	}
}

func Get[T any](client *HttpClient, url string, headers map[string]string) (T, error) {
	var m T

	r, err := http.NewRequestWithContext(client.ctx, "GET", url, nil)
	if err != nil {
		return m, err
	}

	// Important to set
	r.Header.Add("Content-Type", "application/json")

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	res, err := client.httpClient.Do(r)
	if err != nil {
		return m, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m, err
	}

	return parseJSON[T](body)
}

func Post[T any](client *HttpClient, url string, data any, headers map[string]string) (T, error) {
	var m T

	b, err := toJSON(data)
	if err != nil {
		return m, err
	}

	byteReader := bytes.NewReader(b)
	r, err := http.NewRequestWithContext(client.ctx, "POST", url, byteReader)
	if err != nil {
		return m, err
	}

	// Important to set
	r.Header.Add("Content-Type", "application/json")

	for key, value := range headers {
		r.Header.Set(key, value)
	}

	res, err := client.httpClient.Do(r)
	if err != nil {
		return m, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return m, err
	}

	return parseJSON[T](body)
}

func parseJSON[T any](s []byte) (T, error) {
	var r T
	if err := json.Unmarshal(s, &r); err != nil {
		return r, err
	}
	return r, nil
}

func toJSON(T any) ([]byte, error) {
	return json.Marshal(T)
}
