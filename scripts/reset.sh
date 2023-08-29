echo "Restarting Docker containers"

echo "---Reset database migrations---"

source .env
cd ./sql/schema

goose mysql "$DSN" reset