echo "Flushing Redis database"
redis-cli FLUSHALL

echo "\nAdding mock auctions via POST /api/auction"
curl -d "id=1&name=bob&firstBid=0&sellerId=0&status=Mock Data 1" localhost:8080/api/auction
curl -d "id=2&name=sally&firstBid=0&sellerId=0&status=Mock Data 2" localhost:8080/api/auction
curl -d "id=3&name=cindy&firstBid=0&sellerId=0&status=Mock Data 3" localhost:8080/api/auction
curl -d "id=4&name=olivia&firstBid=0&sellerId=0&status=Mock Data 4" localhost:8080/api/auction
curl -d "id=5&name=matt&firstBid=0&sellerId=0&status=Mock Data 5" localhost:8080/api/auction
curl -d "id=6&name=tommy&firstBid=0&sellerId=0&status=Mock Data 6" localhost:8080/api/auction
curl -d "id=7&name=jack&firstBid=0&sellerId=0&status=Mock Data 7" localhost:8080/api/auction
curl -d "id=8&name=lauren&firstBid=0&sellerId=0&status=Mock Data 8" localhost:8080/api/auction
curl -d "id=9&name=abi&firstBid=0&sellerId=0&status=Mock Data 9" localhost:8080/api/auction
curl -d "id=10&name=vlad&firstBid=0&sellerId=0&status=Mock Data 10" localhost:8080/api/auction

echo "\nDone!"
