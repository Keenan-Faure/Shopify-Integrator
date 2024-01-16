#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/update_app.sh

function err_msg() {
    echo "error: invalid argument"
    echo "error: either enter:"
    echo "error: ./scripts/update.sh development -  to update your development database"
    echo "error: ./scripts/update.sh production  -  to update your production (docker) database"
    exit 1;
}

function production_update() {
    source .env
    # make conditions here based on which migrations they want to do
    docker exec $SERVER_CONTAINER_NAME bash -c "cd /keenan/sql/schema && /keenan/scripts/migrations.sh '$1' '$2'"
    docker restart $SERVER_CONTAINER_NAME
}

function dev_update() {
    ./scripts/migrations.sh "$1" "$2"
}

if [[ ! $# -eq 0 ]] ; then
    if [[ "$1" == "development" ]]; then
        if [ -z "$2" ]; then
            echo "error: argument 2 must be either 'up' 'down' 'reset'"
            echo "error: e.g ./scripts/update.sh development up"
            exit;
        fi
        dev_update "$1" "$2"
    elif [[ "$1" == "production" ]]; then
        if [ -z "$2" ]; then
            echo "error: argument 2 must be either 'up' 'down' 'reset'"
            echo "error: e.g ./scripts/update.sh development up"
            exit;
        fi
        production_update "$1" "$2"
    else
        err_msg
    fi
else
    err_msg
fi