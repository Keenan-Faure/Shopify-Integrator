echo "Check goose version"
goose -version

cd ./sql/schema
source .env
goose mysql "$DSN" up

echo "---Completed Database migrations---"