package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	integrationHub "github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/http/auth"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
	"net/http"
	"strconv"
	"time"
)

type TransactionHandler struct {
	TransactionSource ports.TransactionsSource
}

func (t TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	internalID, internalIDExists := mux.Vars(r)["internal_id"]
	if integrationHub.Authorize(w, r) {
		if !internalIDExists {
			NotOKResponseWithBody(
				w,
				400,
				"Missing required parameter in request.",
				models.RequestParameters{
					InternalID: "",
				},
			)
			return
		}

		body := models.NewGetTransactionEventsByInternalIDRequestFromParameters(internalID)
		ctx := context.Background()
		transactions, err := t.TransactionSource.GetTransactionEventsByInternalID(ctx, body)
		if err != nil {
			if err == sql.ErrNoRows {
				NotOKResponseWithBody(
					w,
					404,
					"Transaction not found.",
					models.RequestParameters{
						InternalID: internalID,
					},
				)
				return
			}

			// TODO: Should have better handling for 5## errors
			// Need metrics + alerts
			NotOKResponseWithBody(
				w,
				500,
				"Error when processing request for transaction.",
				models.RequestParameters{
					InternalID: internalID,
				},
			)
			return
		}

		responseBody := models.GetTransactionEventsByInternalIDResponse{Transactions: transactions}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(responseBody)
	}

}

func (t TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	if integrationHub.Authorize(w, r) {
		values := r.URL.Query()
		storeID := values.Get("store_id")
		if storeID == "" {
			NotOKResponseWithBody(
				w,
				400,
				"Missing required parameter in request.",
				models.RequestParameters{
					StoreID: "",
				},
			)
			return
		}
		startDate := values.Get("start_date")
		if startDate == "" {
			// default to 72 hours earlier
			startDate = time.Now().UTC().Add(-time.Hour * time.Duration(72)).Format("2006-01-02T15:04:05.000Z")
		}
		endDate := values.Get("end_date")
		if endDate == "" {
			endDate = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
		}
		before := values.Get("before")
		after := values.Get("after")
		limitStr := values.Get("limit")
		if limitStr == "" {
			// default to max 50 results
			limitStr = "50"
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			// default to max 50 results if the string result from input can't be marshalled to int
			limit = 50
		}

		body := models.NewGetAllTransactionsRequestFromParameters(storeID, startDate, endDate, after, before, limit)
		transactions, prevToken, nextToken, err := t.TransactionSource.GetTransactionsByStoreID(r.Context(), body)
		if err != nil {
			if _, ok := err.(models.InvalidCursors); ok {
				NotOKResponseWithBody(
					w,
					400,
					err.Error(),
					models.RequestParameters{
						StoreID:   storeID,
						StartDate: startDate,
						EndDate:   endDate,
						Before:    before,
						After:     after,
						Limit:     limit,
					},
				)
				return
			} else if err == sql.ErrNoRows {
				NotOKResponseWithBody(
					w,
					404,
					"Transaction(s) not found.",
					models.RequestParameters{
						StoreID:   storeID,
						StartDate: startDate,
						EndDate:   endDate,
						Before:    before,
						After:     after,
						Limit:     limit,
					},
				)
				return
			}

			NotOKResponseWithBody(
				w,
				500,
				"Error when processing request for transactions.",
				models.RequestParameters{
					StoreID:   storeID,
					StartDate: startDate,
					EndDate:   endDate,
					Before:    before,
					After:     after,
					Limit:     limit,
				},
			)
			return
		}

		// TODO: figure out how to check if these forward and backward links have values - probably on cursor level
		// maybe a response at the DB level could say if there is more results forward or backwards
		previousLink := fmt.Sprintf(
			"/transactions?store_id=%s&start_date=%s&end_date=%s&before=%s&limit=%s",
			storeID,
			startDate,
			endDate,
			prevToken,
			limitStr,
		)
		nextLink := fmt.Sprintf(
			"/transactions?store_id=%s&start_date=%s&end_date=%s&after=%s&limit=%s",
			storeID,
			startDate,
			endDate,
			nextToken,
			limitStr,
		)

		// TODO: update response model for Stoplight version
		responseBody := models.GetAllTransactionsResponse{
			Transactions: transactions,
			Pagination: models.Pagination{
				Cursors: models.Cursors{
					Before: prevToken,
					After:  nextToken,
				},
				Paths: models.Paths{
					Previous: previousLink,
					Next:     nextLink,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(responseBody)
	}

}

func NotOKResponseWithBody(w http.ResponseWriter, c int, s string, r models.RequestParameters) {
	responseBody := models.HTTPResponse{
		Code:              c,
		Status:            s,
		RequestParameters: r,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseBody)
}
