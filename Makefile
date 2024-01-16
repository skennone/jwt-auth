## run/auth: run the cmd/auth application
.PHONY: run/auth
run/auth:
	go run ./cmd/auth/ -db-dsn=${Auth_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${Auth_DB_DSN}