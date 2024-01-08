#!/bin/bash
# Note that this is only run within the scope of the `install.sh`` script
# if you wish to update your database to the latest version (migrations)
# please use `update.sh``

function migrations() {
    echo "---Running database migrations---"

    SSL_MODE="?sslmode=disable"
    DRIVER="postgres://"
    if [[ "$1" == "production" ]]; then
        cd /keenan/sql/schema
    elif [[ "$1" == "development" ]]; then
        cd sql/schema
    else
        err_msg
    fi

    echo "Checking GOOSE version"
    goose -version

   if [[ "$1" == "development" ]]; then
        DB_STRING="${DRIVER}${DB_USER}:${DB_PSW}@localhost:5432/${DB_NAME}${SSL_MODE}"
    elif [[ "$1" == "production" ]]; then
        DB_STRING="${DRIVER}${DB_USER}:${DB_PSW}@postgres:5432/${DB_NAME}${SSL_MODE}"
    else
        err_msg
    fi

    goose postgres "$DB_STRING" "$2"
}

function err_msg() {
    echo "invalid argument"
    echo "either enter:"
    echo "./scripts/migrations.sh development  -  to update your development database"
    echo "./scripts/migrations.sh production  -  to update your production (docker) database"
    exit 1;
}

if [[ ! $# -eq 0 ]] ; then
    if [[ "$1" == "development" || "$1" == "production" ]]; then
        migrations "$1" "$2"
    else
        err_msg
    fi
else
    err_msg
fi