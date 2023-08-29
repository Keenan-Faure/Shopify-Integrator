echo "Check goose version"
goose -version

source .env

cd ./sql/schema

goose mysql "$DSN" up

echo "---Completed Database migrations---"