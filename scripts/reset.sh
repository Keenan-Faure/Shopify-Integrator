echo "Restarting Docker containers"

source .env

docker stop $CONTAINER_NAME

docker rm $CONTAINER_NAME

#removes images
if docker image inspect $IMAGE_NAME >/dev/null 2>&1; then
  docker rmi $(docker images $IMAGE_NAME -a -q) -f
else
  echo "'$IMAGE_NAME' does not exist."
fi

echo "---Redo database migrations---"

source .env
cd ./sql/schema

goose mysql "$DSN" reset