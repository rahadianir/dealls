set -a
source .env
docker compose up --build -d
sleep 5 # waiting for the db to be up
migrate -path migrations/sql -database=$DB_URL up
go run main.go
# sleep 15
# migrate -path migrations/sql -database=$DB_URL down
# docker compose down --volumes