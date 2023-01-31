/* When changing a SQL migration script, run
go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/adapters/postgres/migrations/
   -o ./internal/adapters/postgres/migrations/bindata.go ./internal/adapters/postgres/migrations
to update bindata.go
*/

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS transactions_api
(
    id          SERIAL      NOT NULL PRIMARY KEY,
    transaction jsonb       NOT NULL DEFAULT '{}'::jsonb,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_timestamp
    BEFORE UPDATE
    ON transactions_api
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();