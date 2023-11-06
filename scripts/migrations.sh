#!/bin/bash

echo "---running migrations on server container---"

cd /keenan/sql/schema
source .env

echo "Checking GOOSE version"
goose -version

SSL_MODE="?sslmode=disable"
DB_STRING="${DOCKER_DB_URL}${DB_NAME}${SSL_MODE}"

echo "running migrations on '${DB_NAME}'"

goose postgres "$DB_STRING" up