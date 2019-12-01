package main

import (
	"net/http"
)

// Returns all auctions
func getAuctions(writer http.ResponseWriter, request *http.Request) {
	res := auctions
	sendResponse(res, writer)
}

// Adds a new auction
func addAuction(writer http.ResponseWriter, request *http.Request) {
	// Extract data from the request body
	_ = request.ParseForm()
	name := request.Form.Get("name")
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))
	reservePriceS := request.Form.Get("reservePrice")

	// Safety in case no reserve price is given
	var reservePrice float64
	if reservePriceS != "" {
		reservePrice = toFloat(reservePriceS)
	} else {
		reservePrice = 0
	}

	auctionId := assignAuctionId()
	status := "available"

	// Check if auction already exists
	for _, auction := range auctions {
		if auctionId == auction.Id {
			sendResponse(ApiResponse{405, "error", "Invalid input"}, writer)
			return
		}
	}

	auction := Auction{auctionId, status, name, firstBid, sellerId, reservePrice}
	auctions = append(auctions, auction)
	setSyncData("auctions", "add", auction)

	res := ApiResponse{200, "success", "Successful operation"}
	sendResponse(res, writer)
}

// Returns a specific auction by id
func getAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := toInt(getMuxVariable("auctionId", request))

	// Fetch auction by id
	var res Auction
	for _, auction := range auctions {
		if auctionId == auction.Id {
			res = auction
		}
	}

	sendResponse(res, writer)
}

// Updates a specific auction by id
func updateAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := toInt(getMuxVariable("auctionId", request))

	// Extract data from the request body
	_ = request.ParseForm()
	name := request.Form.Get("name")
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))

	// Update auction only if it exists
	for i, auction := range auctions {
		if auctionId == auction.Id {
			auctions[i] = Auction{auction.Id, auction.Status, name, firstBid, sellerId, auction.ReservePrice}
			setSyncData("auctions", "update", auctions[i])
			sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
			return
		}
	}

	sendResponse(ApiResponse{404, "error", "Auction not found"}, writer)
}

// Removes a specific auction by ic
func deleteAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := toInt(getMuxVariable("auctionId", request))

	// Delete auction only if it exists
	for i, auction := range auctions {
		if auctionId == auction.Id {
			auctions = append(auctions[:i], auctions[i+1:]...)
			setSyncData("auctions", "remove", auction)
			sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
			return
		}
	}

	sendResponse(ApiResponse{404, "error", "Auction not found"}, writer)
}
