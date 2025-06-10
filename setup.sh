echo "setting up env from .env file"
set -a
. .env

echo "setting up containered database"
docker compose up --build -d
sleep 5 # waiting for the db to be up

echo "migrating database schema and mock data"
migrate -path migrations/sql -database=$DB_URL up
go run main.go migrate

echo "setup finished!"