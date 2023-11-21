#!/bin/bash

echo "---running migrations on server container---"

SSL_MODE="?sslmode=disable"
DRIVER="postgres://"

if [ -z "$1" ]; then
    if [ $1 = "ci" ]; then
        cd /sql/schema
        echo "Checking GOOSE version"
        goose -version
        DB_STRING="${DRIVER}${DB_USER}${DB_PSW}@localhost:5432/${DB_NAME}${SSL_MODE}"
        echo "running migrations on '${DB_NAME}'"
        goose postgres "$DB_STRING" up
    else
        echo "invalid parameter passed, expected 'ci'"
        exit;
    fi
    # postgres://postgres:postgres@postgres:5432/
else
    cd /keenan/sql/schema
    source .env

    echo "Checking GOOSE version"
    goose -version
    DB_STRING="${DRIVER}${DB_USER}${DB_PSW}@postgres:5432/${DB_NAME}${SSL_MODE}"
    echo "running migrations on '${DB_NAME}'"
    goose postgres "$DB_STRING" up
fi