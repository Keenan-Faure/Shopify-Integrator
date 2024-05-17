#!/bin/bash

function build_go() {
    OS="$(uname -s)"

    # Builds the go code depending of OS
    if [ $OS == "Darwin" ]; then
        echo "OSX detected"
        GOOS=linux GOARCH=amd64 go build -o integrator
    else
        echo "Linux detected"
        go build -o integrator
    fi
}

function create_env() {
    if [ ! -f '.env' ]; then
        if [ -f '.example.env' ]; then
            cp .example.env .env
        fi
    fi
}

function is_development() {
    export GIN_MODE=debug
    export DOCKER_FILE="Dockerfile.github"
}

function install_app_gh() {
    source .env
    is_development
    docker-compose rm -f
    if ! docker compose up -d --force-recreate --no-deps postgres server; then
        exit
    else
        until
            docker exec $DB_NAME pg_isready
        do
            sleep 3;
        done
        docker exec $SERVER_CONTAINER_NAME bash -c "/keenan/scripts/migrations.sh 'production' 'up'"
        docker restart $SERVER_CONTAINER_NAME
    fi
}

build_go
create_env
install_app_gh
