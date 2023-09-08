#!/bin/bash
cd /web-app/sql/schema
source .env

echo "Checking GOOSE version"
goose -version

SSL_MODE="?sslmode=disable"
DB_STRING="${DOCKER_DB_URL}${DATABASE}${SSL_MODE}"

echo "running migrations on '${DATABASE}'"

goose postgres "$DB_STRING" up