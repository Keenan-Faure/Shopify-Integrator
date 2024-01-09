#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/reset.sh

source .env
list="$APP_CONTAINER_NAME $SERVER_CONTAINER_NAME $DOCS_CONTAINER_NAME $DB_NAME $NGROK_CONTAINER_NAME"
image_list="$IMAGE_NAME $DOCS_IMAGE_NAME $APP_IMAGE_NAME"

function docker_stop() {
  echo "--stopping active containers--"
  for item in $list
  do
    docker stop $item
  done
}

function docker_rm() {
  echo "--removing containers--"
  for item in $list
  do
    docker rm $item
  done
}

function docker_rmi() {
  echo "--removing container images--"
  for item in $image_list
  do
    if docker image inspect $item >/dev/null 2>&1; then
      docker rmi $(docker images $item -a -q) -f
    else
      echo "'$item' does not exist."
    fi
  done
}

docker_stop
docker_rm

if [[ ! $# -eq 0 ]] ; then
    if [[ "$1" == "rmi" ]]; then
      docker_rmi
    fi
fi

echo "+-------------------------------------+"
echo "|Please re-run './scripts/install.sh' |"
echo "+-------------------------------------+"