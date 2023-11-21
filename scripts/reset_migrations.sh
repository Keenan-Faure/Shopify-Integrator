#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/reset_migrations.sh

echo "---reset database migrations---"

cd ./sql/schema

SSL_MODE="?sslmode=disable"
DB_STRING="${DRIVER}${DB_USER}${DB_PSW}@localhost:5432/${DB_NAME}${SSL_MODE}"

if ! goose postgres "$DB_STRING" reset ; then

else
    echo "re-run migrations"
    goose postgres "$DB_STRING" up
fi