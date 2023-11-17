#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/reset.sh

echo "---reset containers---"

source .env

docker stop $APP_CONTAINER_NAME
docker stop $SERVER_CONTAINER_NAME
docker stop $DOCS_CONTAINER_NAME
docker stop $DB_NAME
docker stop $NGROK_CONTAINER_NAME

echo "---"
echo "removing containers"

docker rm $SERVER_CONTAINER_NAME
docker rm $APP_CONTAINER_NAME
docker rm $DOCS_CONTAINER_NAME
docker rm $DB_NAME
docker rm $NGROK_CONTAINER_NAME

#removes images

echo "---"
echo "removing images containers"

if docker image inspect $IMAGE_NAME >/dev/null 2>&1; then
  docker rmi $(docker images $IMAGE_NAME -a -q) -f
else
  echo "'$IMAGE_NAME' does not exist."
fi

if docker image inspect $APP_IMAGE_NAME >/dev/null 2>&1; then
  docker rmi $(docker images $APP_IMAGE_NAME -a -q) -f
else
  echo "'$APP_IMAGE_NAME' does not exist."
fi

if docker image inspect $DOCS_IMAGE_NAME >/dev/null 2>&1; then
  docker rmi $(docker images $DOCS_IMAGE_NAME -a -q) -f
else
  echo "'$DOCS_IMAGE_NAME' does not exist."
fi

echo "Please re-run './scripts/install.sh'"