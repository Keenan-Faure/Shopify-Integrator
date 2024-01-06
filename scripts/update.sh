#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/update_app.sh

source .env

docker exec $SERVER_CONTAINER_NAME bash -c "cd /keenan/sql/schema && /keenan/scripts/migrations.sh"