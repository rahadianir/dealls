echo "cleaning up data in database"
echo y | migrate -path migrations/sql -database=$DB_URL down 

echo "cleaning up database instance"
docker compose down --volumes