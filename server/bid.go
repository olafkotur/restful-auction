package main

import (
	"encoding/json"
	"net/http"
)

func addAuctionBid(writer http.ResponseWriter, request *http.Request) {
	bidId := assignKeyId("bid")
	auctionId := getMuxVariable("auctionId", request)

	// Extract data from the request body
	_ = request.ParseForm()
	bidAmount := toFloat(request.Form.Get("bidAmount"))
	bidderId := toInt(request.Form.Get("bidderId"))

	// Get auction details
	auctionData := Auction{}
	auction, _ := client.Get("auction:" + auctionId).Result()
	_ = json.Unmarshal([]byte(auction), &auctionData)

	// Check if the requested bid is higher than previous bids
	keys := client.Keys("bid:*").Val()
	for _, key := range keys {
		bidData := Bid{}
		bid, _ := client.Get(key).Result()
		_ = json.Unmarshal([]byte(bid), &bidData)
		if bidData.BidAmount >= bidAmount {
			sendResponse(ApiResponse{2, "error", "The bid amount must be greater than the highest bid"}, writer)
			return
		}
	}

	// Add auction first bid if not yet set
	if auctionData.FirstBid <= 0 {
		item, _ := json.Marshal(Auction{auctionData.Id, auctionData.Status, auctionData.Name, bidAmount, auctionData.SellerId})
		client.Set("auction:"+auctionId, item, 0)
	}

	// Add bid to database
	item, _ := json.Marshal(Bid{bidId, toInt(auctionId), bidAmount, bidderId})
	client.Set("bid:"+toString(bidId), item, 0)

	printRequestInfo(request)
	sendResponse(ApiResponse{1, "Success", "The bid has been placed"}, writer)
}

func getBidsByAuctionId(writer http.ResponseWriter, request *http.Request) {
	res := []Bid{}
	auctionId := toInt(getMuxVariable("auctionId", request))

	// Get bids
	keys := client.Keys("bid:*").Val()
	for _, key := range keys {
		bidData := Bid{}
		bid, _ := client.Get(key).Result()
		_ = json.Unmarshal([]byte(bid), &bidData)
		if bidData.AuctionId == auctionId {
			res = append(res, bidData)
		}
	}

	printRequestInfo(request)
	sendResponse(res, writer)
}
