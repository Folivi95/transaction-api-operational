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

	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/logger/zapctx"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/adapters/postgres/migrations"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/models"
	"github.com/saltpay/transaction-api-operational/src-luigi-transformer/internal/application/ports"
)

type DBHandler struct {
	db            *sql.DB
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

	err = checkDBVersion(db, 3)
	if err != nil {
		return DBHandler{}, err
	}

	return DBHandler{db: db, metricsClient: metricsClient}, nil
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

func checkDBVersion(db *sql.DB, schemaVersion uint) error {
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

func (h DBHandler) Store(ctx context.Context, transactionJSON []byte) error {
	// Check DB Health
	err := h.db.Ping()
	if err != nil {
		zapctx.From(ctx).Warn("[DBHandler] db not healthy", zap.Error(err))
		return err
	}
	_, err = h.db.ExecContext(ctx, "INSERT INTO transactions_api(transaction) VALUES($1)", transactionJSON)
	if err != nil {
		err = fmt.Errorf("[DBHandler] unable to insert to database, err: %v", err)
		return err
	}

	h.metricsClient.Count("db_write_counter", 1, []string{"tapi"})
	return nil
}

func (h DBHandler) StoreAuxTable(ctx context.Context, transaction models.W4IncomingTransaction, insert bool) error {
	if insert {
		_, err := h.db.ExecContext(ctx, "INSERT INTO aux_transactions_api(id, salt_token_id, expiry_date) VALUES($1, $2, $3) ON CONFLICT (id) DO UPDATE SET id = $1, salt_token_id = $2, expiry_date = $3;", transaction.After.ID, transaction.After.SaltTokenID, transaction.After.CardExpire)
		if err != nil {
			err = fmt.Errorf("[DBHandler] unable to insert to aux database, err: %v", err)
			return err
		}

		h.metricsClient.Count("db_write_counter", 1, []string{"aux_table"})
		return nil
	}

	// if it has not Token data, update the ID to the latest document in chain
	_, err := h.db.ExecContext(ctx, "UPDATE aux_transactions_api SET id = $1 WHERE id = $2", transaction.After.ID, transaction.After.DocPrevID)
	if err != nil {
		err = fmt.Errorf("[DBHandler] unable to insert to aux database, err: %v", err)
		return err
	}

	h.metricsClient.Count("db_write_counter", 1, []string{"aux_table"})
	return nil
}

func (h DBHandler) Upsert(ctx context.Context, entry map[string]interface{}, tableName string) error {
	query, err := QueryBuilder(tableName, "id", entry)
	if err != nil {
		zapctx.From(ctx).Warn("[DBHandler] unable to build query", zap.Error(err))
		return nil
	}
	_, err = h.db.ExecContext(ctx, query)
	if err != nil {
		err = fmt.Errorf("[DBHandler] unable to upsert to database, err: %v", err)
		return err
	}

	h.metricsClient.Count("db_write_counter", 1, []string{tableName})
	return nil
}

func (h DBHandler) Get(ctx context.Context, query string, target []interface{}) error {
	row := h.db.QueryRowContext(ctx, query)

	if err := row.Scan(target...); err != nil {
		if err == sql.ErrNoRows {
			zapctx.From(ctx).Error("[DBHandler] Could not find entries in db", zap.Error(err))
		} else {
			zapctx.From(ctx).Error("[DBHandler] Error when querying", zap.Error(err))
		}
		return err
	}

	h.metricsClient.Count("db_read_counter", 1, []string{})
	return nil
}
