#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/install.sh

function create_worspace() {
    ## Shopify Integrator Docs
    cd ../
    echo "---creating workspace---"

    if [ -d "Shopify-Integrator-docs" ]
    then
        cd Shopify-Integrator-docs
        echo "pulling latest 'Shopify-Integrator-docs'"
        git pull
    else
        git clone "https://github.com/Keenan-Faure/Shopify-Integrator-docs"
        cd Shopify-Integrator-docs
    fi

    ## Shopify Integrator
    cd ../Shopify-Integrator
    echo "pulling latest 'Shopify-Integrator'"
    git pull

    ## Shopify Integrator App
    echo "pulling latest 'Shopify-Integrator App'"
    if [ -d "app" ]
    then
        cd app
        echo "pulling latest 'Shopify-Integrator-docs'"
        git pull
    else
        mkdir app && cd app
        git clone "https://github.com/MrKkyle/Shopify-Integrator-App" .
        cd ../
    fi
}

function check_go() {
    if ! command -v go &> /dev/null
    then
        echo "Golang required but it's not installed."
        echo "Please visit https://go.dev/dl/"
        exit;
    fi
}

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

function install_app() {
    # removes stopped service containers 
    docker-compose rm -f

    echo "---Running Docker compose up---"

    # tells docker to recreate all containers regardless of whether
    # the images have been changed or not 
    if ! docker compose up -d --force-recreate; then
        exit
    else
        source .env
        until
            docker exec $DB_NAME pg_isready
        do 
            sleep 3; 
        done
        docker exec $SERVER_CONTAINER_NAME bash -c "/keenan/scripts/migrations.sh 'production' 'up'"
        docker restart $SERVER_CONTAINER_NAME
    fi
}

create_worspace
check_go
build_go
install_app