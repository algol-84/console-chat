#!/bin/bash
source .env

export MIGRATION_DSN="host=${DB_HOST:-pg_auth} port=5432 dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v