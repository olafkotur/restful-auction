package main

import (
	"encoding/json"
	"net/http"
)

// Fetches all auctions from the db
func getAuctions(writer http.ResponseWriter, request *http.Request) {
	var res []Auction

	// Fetch all keys from the database
	keys := client.Keys("auction:*").Val()
	for _, key := range keys {
		data := Auction{}
		auction, _ := client.Get(key).Result()
		_ = json.Unmarshal([]byte(auction), &data)
		res = append(res, Auction{data.Id, data.Status, data.Name, data.FirstBid, data.SellerId})
	}

	sendResponse(res, writer)
	printRequestInfo(request)
}

// Adds an auction to the db with given auction details
func addAuction(writer http.ResponseWriter, request *http.Request) {
	// Extract data from the request body
	_ = request.ParseForm()
	id := toInt(request.Form.Get("id"))
	name := request.Form.Get("name")
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))
	status := request.Form.Get("status")

	// Check if exists and add to the db
	previousData := Auction{}
	previous, _ := client.Get("auction:" + toString(id)).Result()
	_ = json.Unmarshal([]byte(previous), &previousData)

	if previousData.Id != id {
		item, _ := json.Marshal(Auction{id, status, name, firstBid, sellerId})
		client.Set("auction:"+toString(id), item, 0)
		sendResponse(ApiResponse{1, "success", "The auction has been added"}, writer)
	} else {
		sendResponse(ApiResponse{2, "error", "An auction with this id already exists"}, writer)
	}

	printRequestInfo(request)
}

// Fetches a single auction that matches given auctionId
func getAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := getMuxVariable("auctionId", request)

	// Fetch auction from the db
	data := Auction{}
	auction, _ := client.Get("auction:" + auctionId).Result()
	_ = json.Unmarshal([]byte(auction), &data)
	res := Auction{data.Id, data.Status, data.Status, data.FirstBid, data.SellerId}

	sendResponse(res, writer)
	printRequestInfo(request)
}

// Updates a single auction that matches a given auctionId
func updateAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := getMuxVariable("auctionId", request)

	// Extract data from the POST
	_ = request.ParseForm()
	name := request.Form.Get("name")
	status := request.Form.Get("status")

	// Get exisiting data from db
	previousData := Auction{}
	previous, _ := client.Get("auction:" + auctionId).Result()
	_ = json.Unmarshal([]byte(previous), &previousData)

	// Update record if it exists
	if previousData.Id == toInt(auctionId) {
		item, _ := json.Marshal(Auction{previousData.Id, status, name, previousData.FirstBid, previousData.SellerId})
		client.Set("auction:"+auctionId, item, 0)
		sendResponse(ApiResponse{1, "success", "Auction has been updated"}, writer)
	} else {
		sendResponse(ApiResponse{2, "error", "Auction with specified id does not exist"}, writer)
	}
}

// Deletes a single auction that matches given auctionId
func deleteAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := getMuxVariable("auctionId", request)

	// Delete auction
	res, _ := client.Del("auction:" + auctionId).Result()
	if res != 0 {
		sendResponse(ApiResponse{1, "success", "Auction has been deleted"}, writer)
	} else {
		sendResponse(ApiResponse{2, "error", "Could not delete auction"}, writer)
	}
	printRequestInfo(request)
}
