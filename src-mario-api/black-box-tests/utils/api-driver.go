package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
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

func NewAPIClient(baseURL string, logger APIClientLogger) *APIClient {
	return &APIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 5 * time.Second},
		logger:     logger,
	}
}

const (
	LocalURL = "http://localhost:8080"
)

func GetBaseURL(t *testing.T) string {
	url := os.Getenv("BASE_URL")
	if url == "" {
		url = LocalURL
		startWebserver(t)
	}

	return url
}

func startWebserver(t *testing.T) {
	t.Helper()
	compose := testcontainers.NewLocalDockerCompose(
		[]string{"../../docker-compose.yaml"},
		strings.ToLower(uuid.New().String()),
	)
	webContainer := compose.WithCommand([]string{"up", "-d", "web"})
	invokeErr := webContainer.Invoke()

	if invokeErr.Error != nil {
		t.Fatal(invokeErr)
	}

	t.Cleanup(func() {
		compose.Down()
	})
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

	return fmt.Errorf("given up checking health_check after %dms, %v", time.Since(start).Milliseconds(), err)
}

func (a *APIClient) GetTransaction(internalID string) (models.GetTransactionEventsByInternalIDResponse, error) {
	payload := models.GetTransactionEventsByInternalIDRequest{InternalID: internalID}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(payload)
	if err != nil {
		return models.GetTransactionEventsByInternalIDResponse{}, fmt.Errorf("failed to encode struct to json")
	}

	url := a.baseURL + fmt.Sprintf("/transactions/%s", internalID)
	a.logger.Log("GET", url)

	req, err := http.NewRequest(http.MethodGet, url, &body)
	if err != nil {
		return models.GetTransactionEventsByInternalIDResponse{}, fmt.Errorf("error creating http request: %w", err)
	}

	res, err := a.httpClient.Do(req)
	if err != nil {
		return models.GetTransactionEventsByInternalIDResponse{}, fmt.Errorf("problem reaching %s, %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return models.GetTransactionEventsByInternalIDResponse{}, fmt.Errorf("unexpected status %d from GET %q", res.StatusCode, url)
	}

	var response models.GetTransactionEventsByInternalIDResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return models.GetTransactionEventsByInternalIDResponse{}, fmt.Errorf("unexpected status %d from GET %q", res.StatusCode, url)
	}
	return response, err
}

func (a *APIClient) GetTransactions(after, before, endDate, startDate, storeID string, limit int) (models.GetAllTransactionsResponse, error) {
	bulkRequest := models.GetAllTransactionsRequest{
		After:     after,
		Before:    before,
		EndDate:   endDate,
		StartDate: startDate,
		StoreID:   storeID,
		Limit:     limit,
	}
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(bulkRequest)
	if err != nil {
		return models.GetAllTransactionsResponse{}, fmt.Errorf("failed to encode struct to json")
	}

	url := fmt.Sprintf("%s/transactions?store_id=%s&start_date=%s&end_date=%s&limit=%d", a.baseURL, storeID, startDate, endDate, limit)
	a.logger.Log("GET", url)

	req, err := http.NewRequest(http.MethodGet, url, &body)
	if err != nil {
		return models.GetAllTransactionsResponse{}, fmt.Errorf("error creating http request: %w", err)
	}

	res, err := a.httpClient.Do(req)
	if err != nil {
		return models.GetAllTransactionsResponse{}, fmt.Errorf("problem reaching %s, %w", url, err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return models.GetAllTransactionsResponse{}, fmt.Errorf("unexpected status %d from GET %q", res.StatusCode, url)
	}

	var response models.GetAllTransactionsResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return models.GetAllTransactionsResponse{}, fmt.Errorf("unexpected status %d from GET %q", res.StatusCode, url)
	}
	return response, err
}
