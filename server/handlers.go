package main

import (
	"encoding/json"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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

	var res Response

	// Extract data from the POST
	request.ParseForm()
	id := toInt(request.Form.Get("id"))
	name := request.Form.Get("name")
	firstBid := toFloat(request.Form.Get("firstBid"))
	sellerId := toInt(request.Form.Get("sellerId"))
	status := request.Form.Get("status")

	item, _ := json.Marshal(Auction{id, name, firstBid, sellerId, status})
	err := client.Set("auction:"+toString(id), item, 0)
	if err != nil {
		panic(err)
	}

	sendResponse(res, writer)
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
	type Response struct {
		Status string `json:"status"`
	}

	var res Response

	// Get auctionId from path
	uri := request.URL.String()
	auctionId := strings.Split(uri, "/api/auction/")[1]

	// Delete auction by id from the db
	db := getDatabase()
	statement, _ := db.Prepare("DELETE FROM auctions WHERE id=?")
	r, _ := statement.Exec(auctionId)
	db.Close()

	// Check if any rows have been deleted and set the response
	rowsAffected, _ := r.RowsAffected()
	if rowsAffected > 0 {
		res = Response{"success"}
	} else {
		res = Response{"invalid auctionId"}
	}

	sendResponse(res, writer)
	printRequestInfo(request)
}

// Updates a single auction that matches a given auctionId: {}
func updateAuction(writer http.ResponseWriter, request *http.Request) {

}
