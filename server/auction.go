package main

import (
	"encoding/json"
	"net/http"
)

// Fetches all auctions from the db
func getAuctions(writer http.ResponseWriter, request *http.Request) {
	var res []Auction

	// Fetch all keys from the database and return all auctions
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
	auctionId := assignKeyId("auction")
	status := "available"

	// Extract data from the request body
	_ = request.ParseForm()
	name := request.Form.Get("name")
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))

	// Check if auction already exists, add to database if not
	previousData := Auction{}
	previous, _ := client.Get("auction:" + toString(auctionId)).Result()
	_ = json.Unmarshal([]byte(previous), &previousData)
	if previousData.Id != auctionId {
		item, _ := json.Marshal(Auction{auctionId, status, name, firstBid, sellerId})
		client.Set("auction:"+toString(auctionId), item, 0)
		sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
	} else {
		sendResponse(ApiResponse{405, "error", "Invalid input"}, writer)
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
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))

	// Get exisiting data from db
	previousData := Auction{}
	previous, _ := client.Get("auction:" + auctionId).Result()
	_ = json.Unmarshal([]byte(previous), &previousData)

	// Update record if it exists
	if previousData.Id == toInt(auctionId) {
		item, _ := json.Marshal(Auction{previousData.Id, previousData.Status, name, firstBid, sellerId})
		client.Set("auction:"+auctionId, item, 0)
		sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
	} else {
		sendResponse(ApiResponse{404, "error", "Auction not found"}, writer)
	}
}

// Deletes a single auction that matches given auctionId
func deleteAuction(writer http.ResponseWriter, request *http.Request) {
	auctionId := getMuxVariable("auctionId", request)

	// Delete auction from database
	res, _ := client.Del("auction:" + auctionId).Result()
	if res != 0 {
		sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
	} else {
		sendResponse(ApiResponse{404, "error", "Auction not found"}, writer)
	}
	printRequestInfo(request)
}
