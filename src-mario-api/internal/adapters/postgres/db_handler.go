package postgres

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	goBindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	zapctx "github.com/saltpay/go-zap-ctx"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/adapters/postgres/migrations"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-mario-api/internal/application/ports"
)

type DBHandler struct {
	DB            *sql.DB
	metricsClient ports.MetricsClient
}

const (
	nPingsDefault     = 10
	backoffTime       = 1 * time.Second
	queryLimitDefault = 50
)

func NewDBHandler(connectionString string, metricsClient ports.MetricsClient, nPings int) (DBHandler, error) {
	pgxConfig, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return DBHandler{}, fmt.Errorf("parsing postgres URI: %w", err)
	}

	db := stdlib.OpenDB(*pgxConfig)

	if nPings == 0 {
		nPings = nPingsDefault
	}

	err = pingUntilAvailable(db, nPings)
	if err != nil {
		return DBHandler{}, err
	}

	err = CheckDBVersion(db, 3)
	if err != nil {
		return DBHandler{}, err
	}

	return DBHandler{DB: db, metricsClient: metricsClient}, nil
}

func pingUntilAvailable(db *sql.DB, nPings int) error {
	var err error
	for i := 0; i < nPings; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(backoffTime)
	}
	return err
}

func CheckDBVersion(db *sql.DB, schemaVersion uint) error {
	sourceInstance, err := goBindata.WithInstance(goBindata.Resource(migrations.AssetNames(), migrations.Asset))
	if err != nil {
		return fmt.Errorf("[DBHandler] Error when creating sourceInstance interface: %w", err)
	}

	postgresConfig := new(postgres.Config)
	targetInstance, err := postgres.WithInstance(db, postgresConfig)
	if err != nil {
		return fmt.Errorf("[DBHandler] Error when creating targetInstance interface: %w", err)
	}

	migrateInstance, err := migrate.NewWithInstance("go-bindata", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return fmt.Errorf("[DBHandler] Error when creating targetInstance interface: %w", err)
	}

	err = migrateInstance.Migrate(schemaVersion)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[DBHandler] Migrate failed: %w", err)
	}

	return sourceInstance.Close()
}

// TODO: change object interface to the correct type and maybe have the query here.
func (h DBHandler) Store(ctx context.Context, query string, object []interface{}) error {
	_, err := h.DB.ExecContext(ctx, query, object...)
	if err != nil {
		err := fmt.Errorf("[DBHandler] unable to insert to database, err: %v", err)
		return err
	}
	return nil
}

// TODO: change object interface to the correct type and maybe have the query here.
func (h DBHandler) Get(ctx context.Context, query string, object []interface{}) ([]byte, error) {
	row := h.DB.QueryRowContext(ctx, query, object...)
	var body []byte

	if err := row.Scan(&body); err != nil {
		if err == sql.ErrNoRows {
			zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Could not find %v in DB:", object), zap.Error(err))
		} else {
			h.metricsClient.Count("db_read_counter", 1, []string{"failed"})
			zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Error when querying %v", object), zap.Error(err))
		}
		return []byte{}, err
	}
	h.metricsClient.Count("db_read_counter", 1, []string{"success"})
	zapctx.Debug(ctx, fmt.Sprintf("[DBHandler] Successfully read %v from DB", object))
	return body, nil
}

func (h DBHandler) SaveTransaction(ctx context.Context, transaction models.Transaction) error {
	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		zapctx.Error(ctx, "[DBHandler] Unable to unmarshal message to transaction", zap.Error(err))
	}

	err = pingUntilAvailable(h.DB, nPingsDefault)
	if err != nil {
		zapctx.Warn(ctx, "[DBHandler] db not healthy", zap.Error(err))
		return err
	}
	_, err = h.DB.ExecContext(ctx, "INSERT INTO transactions_api(transaction) VALUES($1)", transactionJSON)
	if err != nil {
		err = fmt.Errorf("[DBHandler] unable to insert to database, err: %v", err)
		return err
	}

	h.metricsClient.Count("db_write_counter", 1, []string{"tapi"})
	return nil
}

func (h DBHandler) GetTransactionEventsByInternalID(ctx context.Context, request models.GetTransactionEventsByInternalIDRequest) ([]models.Transaction, error) {
	row := h.DB.QueryRowContext(ctx, "SELECT transaction FROM transactions_api WHERE transaction->>'id' = $1", request.InternalID)
	var body []byte

	if err := row.Scan(&body); err != nil {
		if err == sql.ErrNoRows {
			zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Could not find %v in DB:", request.InternalID), zap.Error(err))
		} else {
			zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Error when querying transaction %v", request.InternalID), zap.Error(err))
			h.metricsClient.Count("db_transacation_counter_request", 1, []string{"failed"})
		}
		return []models.Transaction{}, err
	}

	transaction, err := models.NewTransactionFromJSON(body)
	if err != nil {
		zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Could not find parse object with id %s:", request.InternalID), zap.Error(err))
		return []models.Transaction{}, err
	}
	h.metricsClient.Count("db_transacation_counter_request", 1, []string{"success"})
	zapctx.Debug(ctx, fmt.Sprintf("[DBHandler] Successfully read %v from DB", request.InternalID))
	return []models.Transaction{transaction}, nil // TODO: come back and refactor to retrieve all transactions related to this id
}

func (h DBHandler) GetTransactionsByStoreID(ctx context.Context, request models.GetAllTransactionsRequest) ([]models.Transaction, string, string, error) {
	if request.Limit < 1 {
		request.Limit = queryLimitDefault
	}

	var row *sql.Rows
	var err error

	if len(request.After) > 0 && len(request.Before) > 0 {
		return []models.Transaction{}, "", "", models.InvalidCursors{}
	} else if len(request.After) < 1 && len(request.Before) < 1 {
		row, err = h.DB.QueryContext(ctx, `SELECT id, transaction FROM transactions_api WHERE (transaction->'card_acceptor'->>'id' = $1) ORDER BY id DESC LIMIT $2`, request.StoreID, request.Limit)
	} else if len(request.After) > 1 {
		pointer, _ := base64.StdEncoding.DecodeString(request.After)
		row, err = h.DB.QueryContext(ctx, `SELECT id, transaction FROM transactions_api WHERE (transaction->'card_acceptor'->>'id' = $1 AND id < $2) ORDER BY id DESC LIMIT $3`, request.StoreID, string(pointer), request.Limit)
	} else {
		pointer, _ := base64.StdEncoding.DecodeString(request.Before)
		row, err = h.DB.QueryContext(ctx, `SELECT id, transaction FROM transactions_api WHERE (transaction->'card_acceptor'->>'id' = $1 AND id > $2) ORDER BY id ASC LIMIT $3`, request.StoreID, string(pointer), request.Limit)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Could not find transactions for %s in DB:", request.StoreID), zap.Error(err))
		} else {
			h.metricsClient.Count("db_read_counter", 1, []string{"failed"})
			zapctx.Error(ctx, fmt.Sprintf("[DBHandler] Error when querying transactions for %v", request.StoreID), zap.Error(err))
		}
		return []models.Transaction{}, "", "", err
	}

	var transactions []models.Transaction
	var finalID string
	prevID := ""
	for row.Next() {
		var transactionRaw []byte
		err = row.Scan(&finalID, &transactionRaw)
		if err != nil {
			h.metricsClient.Count("db_transaction_counter_request", 1, []string{"failed"})
			zapctx.Warn(ctx, "[DBHandler] error when reading DBs row:", zap.Error(err))
			continue
		}

		if prevID == "" {
			prevID = finalID
		}

		transaction, err := models.NewTransactionFromJSON(transactionRaw)
		if err != nil {
			return []models.Transaction{}, "", "", err
		}
		h.metricsClient.Count("db_transaction_counter_request", 1, []string{"success"})
		transactions = append(transactions, transaction)
	}

	finalIDToken := base64.StdEncoding.EncodeToString([]byte(finalID))
	prevIDToken := base64.StdEncoding.EncodeToString([]byte(prevID))
	zapctx.Debug(ctx, fmt.Sprintf("[DBHandler] Successfully read %v from DB", request.StoreID))
	return transactions, prevIDToken, finalIDToken, nil
}
