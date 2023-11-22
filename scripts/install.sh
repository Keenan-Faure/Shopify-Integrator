#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/install.sh

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
cd ../Shopify-Integrator
echo "pulling latest 'Shopify-Integrator'"
git pull

OS="$(uname -s)"

if ! command -v go &> /dev/null
then
    echo "Golang required but it's not installed."
    echo "Please visit https://go.dev/dl/"
    exit;
fi

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
    chmod +x ./sql/schema/migrations.sh
    docker exec $SERVER_CONTAINER_NAME bash -c ./sql/schema/migrations.sh

    docker restart $SERVER_CONTAINER_NAME
fi