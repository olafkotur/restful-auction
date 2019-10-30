package main

import (
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Fetches all auctions from the db: [{}, {}...]
func getAuctions(writer http.ResponseWriter, request *http.Request) {
	type ResponseItem struct {
		Id       int     `json:"id"`
		Name     string  `json:"name"`
		FirstBid float32 `json:"firstbid"`
		SellerId int     `json:"sellerId"`
		Status   string  `json:"status"`
	}
	type Response []ResponseItem

	var res Response
	var id, sellerId int
	var firstBid float32
	var name, status string

	// Fetch db and traverse each row adding to the response array
	rows, _ := db().Query("SELECT * FROM auctions")
	for rows.Next() {
		rows.Scan(&id, &name, &firstBid, &sellerId, &status)
		object := ResponseItem{id, name, firstBid, sellerId, status}
		res = append(res, object)
	}
	db().Close()
	rows.Close()

	sendResponse(res, writer)
	printRequestInfo(request)
}

// Adds an auction to the db with given auction details
func addAuction(writer http.ResponseWriter, request *http.Request) {
	type Response struct {
		Id       int     `json:"id"`
		Name     string  `json:"name"`
		FirstBid float32 `json:"firstbid"`
		SellerId int     `json:"sellerId"`
		Status   string  `json:"status"`
	}

	// Extract data from the POST
	request.ParseForm()
	id := request.Form.Get("id")
	name := request.Form.Get("name")
	firstBid := request.Form.Get("firstBid")
	sellerId := request.Form.Get("sellerId")
	status := request.Form.Get("status")

	// Add a new auction to the db with given schema
	statement, _ := db().Prepare("INSERT INTO auctions(id, name, firstBid, sellerId, status) values(?, ?, ?, ?, ?)")
	statement.Exec(id, name, firstBid, sellerId, status)

	var res Response
	sendResponse(res, writer)
	printRequestInfo(request)
}

// Fetches a single auction that matches given auctionId: {}
func getAuction(writer http.ResponseWriter, request *http.Request) {
	type Response struct {
		Id       int     `json:"id"`
		Name     string  `json:"name"`
		FirstBid float32 `json:"firstbid"`
		SellerId int     `json:"sellerId"`
		Status   string  `json:"status"`
	}

	var res Response
	var id, sellerId int
	var firstBid float32
	var name, status string

	// Get auctionId from path
	uri := request.URL.String()
	auctionId := strings.Split(uri, "/api/auction/")[1]

	// Fetch db and traverse each row setting the response
	rows, _ := db().Query("SELECT * FROM auctions WHERE id=" + auctionId)
	for rows.Next() {
		rows.Scan(&id, &name, &firstBid, &sellerId, &status)
		res = Response{id, name, firstBid, sellerId, status}
	}
	db().Close()
	rows.Close()

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
	statement, _ := db().Prepare("DELETE FROM auctions WHERE id=?")
	r, _ := statement.Exec(auctionId)

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
