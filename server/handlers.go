package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getAuctions(writer http.ResponseWriter, request *http.Request) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()

	type ResponseItem struct {
		Id       int     `json:"id"`
		Name     string  `json:"name"`
		FirstBid float32 `json:"firstbid"`
		SellerId int     `json:"sellerId"`
		Status   string  `json:"status"`
	}

	type Response []ResponseItem

	rows, _ := db.Query("SELECT * FROM auctions")
	var id, sellerId int
	var firstBid float32
	var name, status string
	var res Response

	// Traverse each row, adds data to the res array
	for rows.Next() {
		rows.Scan(&id, &name, &firstBid, &sellerId, &status)
		object := ResponseItem{id, name, firstBid, sellerId, status}
		res = append(res, object)
	}
	rows.Close()
	db.Close()

	response, _ := json.Marshal(res)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}

func getAuction(writer http.ResponseWriter, request *http.Request) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()

	uri := request.URL.String()
	rId := strings.Split(uri, "/api/auction/")[1]

	type Response struct {
		Id       int     `json:"id"`
		Name     string  `json:"name"`
		FirstBid float32 `json:"firstbid"`
		SellerId int     `json:"sellerId"`
		Status   string  `json:"status"`
	}

	rows, _ := db.Query("SELECT * FROM auctions WHERE id=" + rId)
	var id, sellerId int
	var firstBid float32
	var name, status string
	var res Response

	// Traverse each row, and set the response
	for rows.Next() {
		rows.Scan(&id, &name, &firstBid, &sellerId, &status)
		res = Response{id, name, firstBid, sellerId, status}
	}
	rows.Close()
	db.Close()

	response, _ := json.Marshal(res)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}

func deleteAuction(writer http.ResponseWriter, request *http.Request) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()

	uri := request.URL.String()
	rId := strings.Split(uri, "/api/auction/")[1]

	type Response struct {
		Status string `json:"status"`
	}

	var res Response
	statement, _ = db.Prepare("DELETE FROM auctions WHERE id=?")
	r, _ := statement.Exec(rId)

	// Check if any rows have been deleted
	rowsAffected, _ := r.RowsAffected()
	if rowsAffected > 0 {
		res = Response{"success"}
	} else {
		res = Response{"invalid auctionId"}
	}

	response, _ := json.Marshal(res)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}

func updateAuction(writer http.ResponseWriter, request *http.Request) {

}
