#!/bin/bash
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
    chmod +x ./sql/schema/migrations.github.sh
    docker exec $SERVER_CONTAINER_NAME bash -c ./sql/schema/migrations.github.sh

    docker restart $SERVER_CONTAINER_NAME
fi