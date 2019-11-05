package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Auction struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	FirstBid float64 `json:"firstbid"`
	SellerId int     `json:"sellerId"`
	Status   string  `json:"status"`
}

// Fetches all auctions from the db: [{}, {}...]
func getAuctions(writer http.ResponseWriter, request *http.Request) {
	var res []Auction

	// Fetch all keys from the database
	keys := client.Keys("auction:*").Val()
	for _, key := range keys {
		data := Auction{}
		auction, _ := client.Get(key).Result()
		json.Unmarshal([]byte(auction), &data)
		res = append(res, Auction{data.Id, data.Name, data.FirstBid, data.SellerId, data.Status})
	}

	sendResponse(res, writer)
	printRequestInfo(request)
}

// Adds an auction to the db with given auction details
func addAuction(writer http.ResponseWriter, request *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}

	// Extract data from the POST
	request.ParseForm()
	id := toInt(request.Form.Get("id"))
	name := request.Form.Get("name")
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))
	status := request.Form.Get("status")

	// Check if exists and add to the db
	previousData := Auction{}
	previous, _ := client.Get("auction:" + toString(id)).Result()
	json.Unmarshal([]byte(previous), &previousData)

	if previousData.Id != id {
		item, _ := json.Marshal(Auction{id, name, firstBid, sellerId, status})
		client.Set("auction:"+toString(id), item, 0)
		sendSuccessResponse(writer)
	} else {
		sendFailedResponse(writer)
	}

	printRequestInfo(request)
}

// Fetches a single auction that matches given auctionId: {}
func getAuction(writer http.ResponseWriter, request *http.Request) {

	// Get auctionId from path
	uri := request.URL.String()
	auctionId := strings.Split(uri, "/api/auction/")[1]

	// Fetch auction from the db
	data := Auction{}
	auction, _ := client.Get("auction:" + auctionId).Result()
	json.Unmarshal([]byte(auction), &data)
	res := Auction{data.Id, data.Name, data.FirstBid, data.SellerId, data.Status}

	sendResponse(res, writer)
	printRequestInfo(request)
}

// Deletes a signle auction that matches given auctionId: {}
func deleteAuction(writer http.ResponseWriter, request *http.Request) {

	// Get auctionId from path
	uri := request.URL.String()
	auctionId := strings.Split(uri, "/api/auction/")[1]

	res, _ := client.Del("auction:" + auctionId).Result()
	if res != 0 {
		sendSuccessResponse(writer)
	} else {
		sendFailedResponse(writer)
	}
	printRequestInfo(request)
}

// Updates a single auction that matches a given auctionId: {}
func updateAuction(writer http.ResponseWriter, request *http.Request) {
	// Get auctionId from path
	uri := request.URL.String()
	auctionId := strings.Split(uri, "/api/auction/")[1]

	auctionId = auctionId
}
