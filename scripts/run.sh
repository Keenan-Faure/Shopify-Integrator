#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/run.sh

echo "---run containers---"

echo "pulling latest from remote"

git pull

OS="$(uname -s)"

# Builds the go code depending of OS
if [ $OS == "Darwin" ]; then
    echo "OSX detected"
    GOOS=linux GOARCH=amd64 go build -o integrator
else
    echo "Linux detected"
    go build -o integrator
fi

docker-compose rm -f

echo "---Running Docker compose up---"

if ! docker compose up -d --force-recreate; then
    exit
else 
    source .env
    until 
        docker exec $DB_NAME pg_isready
    do 
        sleep 3; 
    done

    echo "---Running database migrations---"
    docker exec $SERVER_CONTAINER_NAME bash -c ./sql/schema/migrations.sh

    docker restart $SERVER_CONTAINER_NAME
fi