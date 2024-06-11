#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/install.sh

function install_hosts() {
    chmod +x ./scripts/hosts.sh
    ./scripts/hosts.sh
}

function check_prerequisites() {
    if [ ! -f '.env' ]
    then
        echo "error: .env not configured, please setup."
        exit;
    fi

    if [ ! -f './ngrok/ngrok.yml' ]
    then
        echo "error: ngrok.yml not configured, please setup."
        exit;
    fi
}

function create_workspace() {
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
        git pull
        cd ../
    else
        mkdir app && cd app
        git clone "https://github.com/MrKkyle/Shopify-Integrator-App" .
        cd ../
    fi
}

function check_go() {
    if ! command -v go &> /dev/null
    then
        echo "error: Golang required but it's not installed."
        echo "error: Please visit https://go.dev/dl/"
        exit;
    fi
}

function check_docker() {
    if ! command -v docker &> /dev/null
    then
        echo "error: Docker required but it's not installed or running"
        echo "error: Please visit https://www.docker.com/"
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
        echo "error: could not run docker compose"
        exit
    else
        source .env
        until
            docker exec $DB_NAME pg_isready
        do 
            sleep 3; 
        done
        # waits for the docker container to be running
        until
            [ "$( docker container inspect -f '{{.State.Status}}' $SERVER_CONTAINER_NAME )" = "running" ]
        do
            echo "waiting for $SERVER_CONTAINER_NAME"
            sleep 3;
        done
        docker exec $SERVER_CONTAINER_NAME bash -c "/keenan/scripts/migrations.sh 'production' 'up'"
        docker restart $SERVER_CONTAINER_NAME
    fi
}

check_prerequisites
create_workspace
check_docker
#check_go
#build_go
install_app

echo "+-------------------------------------+"
echo "|If you wish to install hosts         |"
echo "|run 'sudo ./scripts/hosts.sh'        |"
echo "|[password prompt will appear]        |"
echo "+-------------------------------------+"