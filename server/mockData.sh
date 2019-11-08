echo "Flushing Redis database"
redis-cli FLUSHALL

echo "\nAdding mock auctions via POST /api/auction"
curl -d "id=1&name=Mock1&firstBid=5&sellerId=0" localhost:8080/api/auction
curl -d "id=2&name=Mock1&firstBid=4&sellerId=0" localhost:8080/api/auction
curl -d "id=3&name=Mock1&firstBid=3&sellerId=0" localhost:8080/api/auction
curl -d "id=4&name=Mock1&firstBid=2&sellerId=0" localhost:8080/api/auction
curl -d "id=5&name=Mock1&firstBid=1&sellerId=0" localhost:8080/api/auction

echo "\nDone!"
