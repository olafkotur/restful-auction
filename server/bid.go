package main

import "net/http"

func setAuctionBid(writer http.ResponseWriter, request *http.Request) {
	// auctionId := getMuxVariable("auctionId", request)

	printRequestInfo(request)
}
