package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// -------- Auctions --------
func GetAuctionsResponse(db *sql.DB) (r []byte) {
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

	response, _ := json.Marshal(res)
	return response
}

func auctions(writer http.ResponseWriter, request *http.Request) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()

	var response []byte
	uri := request.URL.String()

	if uri == "/api/auctions" {
		response = GetAuctionsResponse(db)
	}

	db.Close()

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}

func auction(writer http.ResponseWriter, request *http.Request) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()

	var response []byte

	response = []byte("yeet")

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}

func user(writer http.ResponseWriter, request *http.Request) {
	var response []byte

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}
