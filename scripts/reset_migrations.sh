#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/reset_migrations.sh

echo "---Reset database migrations---"

cd ./sql/schema

SSL_MODE="?sslmode=disable"
DB_STRING="${DOCKER_DB_URL}${DATABASE}${SSL_MODE}"

goose postgres "$DB_STRING" reset