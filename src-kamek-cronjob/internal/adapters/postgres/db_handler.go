package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	goBindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"

	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/adapters/postgres/migrations/migrations"
	"github.com/saltpay/transaction-api-operational/src-kamek-cronjob/internal/application/ports"
)

type DBHandler struct {
	DB            *sql.DB
	metricsClient ports.MetricsClient
}

const (
	nPingsDefault = 10
	backoffTime   = 1 * time.Second
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

func (h DBHandler) DeletePastXTransactions(ctx context.Context, interval int) error {
	query := fmt.Sprintf("DELETE FROM transactions_api WHERE created_at < CURRENT_DATE - INTERVAL '%d days'", interval)
	_, err := h.DB.ExecContext(ctx, query)
	if err != nil {
		zapctx.From(ctx).Error(fmt.Sprintf("[DBHandler] Delete past %d days of transaction data failed:", interval), zap.Error(err))
		return err
	}

	zapctx.From(ctx).Info(fmt.Sprintf("[DBHandler] Successfully deleted transaction data that are older than %d days from DB.", interval))
	return nil
}

func (h DBHandler) GetTransaction(ctx context.Context, transactionID string) error {
	row := h.DB.QueryRowContext(ctx, "SELECT transaction FROM transactions_api WHERE transaction->>'id' = $1", transactionID)
	var body []byte

	if err := row.Scan(&body); err != nil {
		if err == sql.ErrNoRows {
			zapctx.From(ctx).Error(fmt.Sprintf("[DBHandler] Could not find %v in DB:", transactionID), zap.Error(err))
		} else {
			zapctx.From(ctx).Error(fmt.Sprintf("[DBHandler] Error when querying transaction %v", transactionID), zap.Error(err))
		}
		return err
	}

	zapctx.From(ctx).Info(fmt.Sprintf("[DBHandler] Successfully read %v from DB", transactionID))
	return nil
}

func (h DBHandler) Store(ctx context.Context, query string, object []interface{}) error {
	_, err := h.DB.ExecContext(ctx, query, object...)
	if err != nil {
		err := fmt.Errorf("[DBHandler] unable to insert to database, err: %v", err)
		return err
	}

	return nil
}
