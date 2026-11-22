Education go lang REST-API project

GO 1.27.4
Postrge 16

## Fetures

- `go mod tidy` required after import
- `go build ./... && go vet ./...` build + test

## Scripts

### Run

```
set -a; source .env; set +a

go run ./cmd/api
```

### Migration:

```
sudo service postgresql status || sudo service postgresql start

set -a; source .env; set +a

psql "$DB_URL" -v ON_ERROR_STOP=1 --single-transaction -f migrations/000001_init.up.sql
```