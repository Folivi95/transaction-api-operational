package handlers

import (
	"encoding/json"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	responseBody := models.HTTPResponse{
		Code:   200,
		Status: "Internal health status - Healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseBody)
	return
}
