#!/bin/bash

OS="$(uname -s)"

# Builds the go code depending of OS
if [ $OS == "Darwin" ]; then
    echo "OSX detected"
    GOOS=linux GOARCH=amd64 go build -o integrator
else
    echo "Linux detected"
    go build -o integrator
fi

source .env
docker-compose rm -f
if ! docker compose up -d --force-recreate --no-deps server postgres; then
    exit
else 
    until
        docker exec $DB_NAME pg_isready
    do 
        sleep 3;
    done
    docker restart $SERVER_CONTAINER_NAME
    sleep 4;
    docker exec $SERVER_CONTAINER_NAME bash -c /keenan/scripts/migrations.sh
    if [ $? -eq 0 ]; then
        echo SUCCEEDED
    else
        echo FAILED
        exit;
    fi

    docker restart $SERVER_CONTAINER_NAME
fi