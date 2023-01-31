# Transaction "Mario" API

Last Updated: 25.03.22

## Postgres Migrations

Database migrations are located at `internal/adapters/postgres/migrations/` and need to be synced with `bindata.go` -
which is in the same folder. To create a new migration, run:

```
$ brew install golang-migrate
$ migrate create -ext sql -dir ./internal/adapters/postgres/migrations -seq name_of_migration
```

To generate Go code and update bindata, run:

```
$ go install github.com/kevinburke/go-bindata/...@latest
$ go-bindata -pkg migrations -ignore bindata -nometadata -prefix internal/adapters/postgres/migrations/ -o ./internal/adapters/postgres/migrations/bindata.go ./internal/adapters/postgres/migrations
```

When ensuring the version of the DB, the system will run these migrations. To update a migration, make the appropriate
changes to the sql files and update the bindata file.

## Endpoints

**BASE URL for platform-dev** : https://transaction-api-operational.platform-dev.eu-west-1.salt

Currently, TAPI has two endpoints: one that returns a single transaction and one that returns a bulk of transactions.

### Single Transaction

`GET /transactions/{acquiring_host}/{transaction_id}`

**Request example**

`GET /transactions/solar/123`

Transaction ID refers to the unique ID given by a transaction by a unique acquiring host. Acquiring host refers to the
host of that transaction (W4, Solar, etc.)

**Success Response**

**Code** : `200 OK`

A successful response returns a json with one transaction. The schema for this TRANSACTION is
defined [here](https://salt.stoplight.io/docs/acquiring-admin-api/branches/feat%2Ftransaction-model/c2NoOjMwOTU0NzQ5-transaction)

```yaml
{
  "transaction": TRANSACTION
}
```

**Code** : `404 NOT FOUND`

```yaml
{
  "code": "TRANSACTION_NOT_FOUND",
  "description": "Transaction ID not found",
}
```

### Bulk Transaction

`GET /transactions/{acquiring_host}`

** Parameters **

```yaml
  "card_acceptor_id": "STRING", // required
  "after": "STRING",            // optional
  "before": "STRING",           // optional
  "limit": INT,                 // optional - default 50
```

**Request example**

`/transactions/solar?card_acceptor_id=some_id&after=1&limit=1`

Card Acceptor ID refers to a unique ID that identifies a unique merchant in Saltpay's system (commonly referred as MID
or contract number). As in singles, acquiring host refers to the host of these transactions (W4, Solar, etc.).
After/Before are unique, encoded, identifiers for our paging mechanism. It is obtained in the http response of this
endpoint and can be used to get a next set of transactions. However, users should use only one of them per request.
Limit refers to the limit of transactions returned by the endpoint (default: 50)

**Success Response**

**Code** : `200 OK`

A successful response returns a json with N transactions. The schema for a TRANSACTION is
defined [here](https://salt.stoplight.io/docs/acquiring-admin-api/branches/feat%2Ftransaction-model/c2NoOjMwOTU0NzQ5-transaction)

```yaml
{
  "transactions": [ ]TRANSACTION,
  "pagination": {
    "cursors": {
      "before": STRING,
      "after": STRING
    },
    "paths": {
      "previous": STRING,
      "next": STRING
    }
  }
}
```

`after` refers to the token of the latest transaction in the response. On the other hand, `before` refers to the first
transaction in the response. These should be used when deciding the next forward/backward request. Links are the path
and parameters for an eventual next request (using either of the cursors)

**Code** : `404 NOT FOUND`

```yaml
{
  "code": "MERCHANT_NOT_FOUND",
  "description": "Merchant ID not found",
}
```
