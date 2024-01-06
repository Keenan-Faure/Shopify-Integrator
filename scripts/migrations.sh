#!/bin/bash

echo "---running migrations on server container---"
goose -version
ls ./
ls /keenan

SSL_MODE="?sslmode=disable"
DRIVER="postgres://"

cd /keenan

source .env

cd sql/schema

echo "Checking GOOSE version"
goose -version
DB_STRING="${DRIVER}${DB_USER}:${DB_PSW}@postgres:5432/${DB_NAME}${SSL_MODE}"
echo ${DB_STRING}"
echo "running migrations on '${DB_NAME}'"
goose postgres "$DB_STRING" up