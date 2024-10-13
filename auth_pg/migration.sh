#!/bin/bash
source .env
@echo ">  Run migration script...."
export MIGRATION_DSN="host=pg_auth port=5432 dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

sleep 4 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v