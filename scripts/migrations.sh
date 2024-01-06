#!/bin/bash

echo "---running migrations on server container---"
source .env

cd sql/schema
goose -version

DB_STRING="postgres://${DB_USER}:${DB_PSW}@postgres:5432/${DB_NAME}?sslmode=disable"
goose postgres "$DB_STRING" up