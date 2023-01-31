package tapi_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIClientLogger interface {
	Log(...interface{})
	Logf(string, ...interface{})
}

type APIClient struct {
	baseURL    string
	httpClient *http.Client
	logger     APIClientLogger
}

const (
	timeout  = time.Minute * 20
	waitTime = time.Second * 15
)

func NewAPIClient(baseURL string, logger APIClientLogger) *APIClient {
	client := &http.Client{Timeout: 15 * time.Second}
	return &APIClient{
		baseURL:    baseURL,
		httpClient: client,
		logger:     logger,
	}
}

func (a *APIClient) CheckIfHealthy() error {
	url := a.baseURL + "/internal/health_check"
	a.logger.Log("GET", url)

	res, err := a.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("problem reaching %s, %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d from POST %q", res.StatusCode, url)
	}

	return nil
}

func (a *APIClient) WaitForAPIToBeHealthy(retries int) error {
	var (
		err   error
		start = time.Now()
	)

	for retries > 0 {
		if err = a.CheckIfHealthy(); err != nil {
			retries -= 1
			time.Sleep(1 * time.Second)
		} else {
			return nil
		}
	}

	return fmt.Errorf("given up checking healthcheck after %dms, %v", time.Since(start).Milliseconds(), err)
}

func (a *APIClient) GetTransaction(singleTransaction SingleTransactionRequest) ([]byte, error) {
	singleTransactionJSON, err := json.Marshal(singleTransaction)
	if err != nil {
		return []byte{}, fmt.Errorf("json marshal error: %w", err)
	}

	body := bytes.NewReader(singleTransactionJSON)
	endTime := time.Now().Add(timeout)
	for time.Until(endTime) > 0 {
		req, err := http.NewRequest(http.MethodGet, a.baseURL+"/transactions/single", body)
		if err != nil {
			return []byte{}, fmt.Errorf("error creating http request: %w", err)
		}

		res, err := a.httpClient.Do(req)
		if err != nil {
			return []byte{}, fmt.Errorf("error during request: %w", err)
		}

		if res.StatusCode != http.StatusOK {
			time.Sleep(waitTime)
			fmt.Sprintf("expected %v, but got %v", http.StatusOK, res.StatusCode)
			continue
		}

		_ = res.Body.Close()
		bodyJSON, _ := io.ReadAll(res.Body)
		return bodyJSON, nil
	}
	return []byte{}, nil
}
