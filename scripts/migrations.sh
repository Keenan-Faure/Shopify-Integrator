#!/bin/bash

echo "---running migrations on server container---"

cd /keenan/sql/schema
source .env

echo "Checking GOOSE version"
goose -version

SSL_MODE="?sslmode=disable"
DRIVER="postgres://"

DB_STRING="${DRIVER}${DB_USER}:${DB_PSW}@postgres:5432/${DB_NAME}${SSL_MODE}"
echo ${DB_STRING}
echo "running migrations on '${DB_NAME}'"
goose postgres "$DB_STRING" up
