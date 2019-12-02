package main

import (
	"net/http"
	"strings"
)

func addAuctionBid(writer http.ResponseWriter, request *http.Request) {
	// Authenticate user
	token := strings.Split(request.Header.Get("Authorization"), "Bearer ")[1]
	isVerified := authenticateToken(token)
	if !isVerified {
		sendResponse(ApiResponse{401, "error", "Unauthorized"}, writer)
		return
	}

	bidId := assignBidId()
	auctionId := toInt(getMuxVariable("auctionId", request))

	// Extract data from the request body
	_ = request.ParseForm()
	bidAmount := toFloat(request.Form.Get("bidAmount"))
	bidderId := toInt(request.Form.Get("bidderId"))

	// Check if bid is higher or equal to than the first bid
	exists := false
	sellerId := 0
	for _, auction := range auctions {
		if auctionId == auction.Id {
			exists = true
			sellerId = auction.SellerId
			if bidAmount < auction.FirstBid {
				sendResponse(ApiResponse{404, "error", "Invalid input"}, writer)
				return
			}
		}
	}

	// Check if the auction exists or is self bid
	if !exists || sellerId == bidderId {
		sendResponse(ApiResponse{404, "error", "Invalid input"}, writer)
		return
	}

	// Check if bid is higher than previous bids
	for _, bid := range bids {
		if auctionId == bid.AuctionId {
			if bidAmount <= bid.BidAmount {
				sendResponse(ApiResponse{404, "error", "Invalid input"}, writer)
				return
			}
		}
	}

	// Add bid to bid data and update redis counter
	bid := Bid{bidId, auctionId, bidAmount, bidderId}
	bids = append(bids, bid)
	setSyncData("bids", "add", bid)

	sendResponse(ApiResponse{200, "Success", "Successful operation"}, writer)
}

func getBidsByAuctionId(writer http.ResponseWriter, request *http.Request) {
	var res []Bid
	auctionId := toInt(getMuxVariable("auctionId", request))

	// Return only bids for the specified auction
	for _, bid := range bids {
		if auctionId == bid.AuctionId {
			res = append(res, bid)
		}
	}

	sendResponse(res, writer)
}
