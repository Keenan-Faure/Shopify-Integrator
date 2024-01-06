#!/bin/bash

echo "---running migrations on server container---"
echo "Checking GOOSE version"
goose -version

SSL_MODE="?sslmode=disable"
DRIVER="postgres://"

cd /keenan
source .env
cd sql/schema

DB_STRING="${DRIVER}${DB_USER}:${DB_PSW}@postgres:5432/${DB_NAME}${SSL_MODE}"
echo "$DB_STRING"
goose postgres "$DB_STRING" up