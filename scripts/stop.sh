#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./stop.sh

echo "---stopping containers---"

source .env

docker stop $APP_CONTAINER_NAME
docker stop $SERVER_CONTAINER_NAME
docker stop $DOCS_CONTAINER_NAME
docker stop $DB_NAME