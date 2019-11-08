#!/bin/bash

# Drop the database
printf "\033[1;31mFlushing Redis database\033[0m\n"
redis-cli FLUSHALL

# Create a new test user
printf "\033[1;32m\nCreating new test user via POST /api/user\033[0m\n"
curl -d "username=olafkotur&password=pi" localhost:8080/api/user && echo
printf "\033[1;36mUsername: olafkotur, Password: pi\033[0m\n"

# Add mock auctions
printf "\033[1;32m\nAdding mock auctions via POST /api/auction\033[0m\n"
curl -d "name=Mock1&firstBid=5&sellerId=1&reservePrice=42" localhost:8080/api/auction && echo
curl -d "name=Mock1&firstBid=4&sellerId=1&reservePrice=42" localhost:8080/api/auction && echo
curl -d "name=Mock1&firstBid=3&sellerId=1&reservePrice=42" localhost:8080/api/auction && echo
curl -d "name=Mock1&firstBid=2&sellerId=1&reservePrice=42" localhost:8080/api/auction && echo
curl -d "name=Mock1&firstBid=1&sellerId=1&reservePrice=42" localhost:8080/api/auction && echo

# Add mock bids
printf "\033[1;32m\nAdding mock bids via POST /api/auction/{auctionId}/bid\033[0m\n"
curl -d "bidAmount=42&bidderId=1" localhost:8080/api/auction/1/bid && echo
curl -d "bidAmount=322&bidderId=1" localhost:8080/api/auction/2/bid && echo
curl -d "bidAmount=12312&bidderId=1" localhost:8080/api/auction/3/bid && echo
curl -d "bidAmount=11&bidderId=1" localhost:8080/api/auction/4/bid && echo
curl -d "bidAmount=6&bidderId=1" localhost:8080/api/auction/5/bid && echo

printf "\033[1;32mDone!\033[0m\n"
