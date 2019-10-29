package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func auctions(writer http.ResponseWriter, request *http.Request) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()

	var response []byte
	method := request.Method
	uri := request.URL.String()

	if method == "GET" {
		if uri == "/api/auctions" {
			response = GetAuctionsResponse(db)
		}
	}

	db.Close()

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
	printRequestInfo(request)
}

func auction(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application: json;")
	fmt.Fprintf(writer, "Auction")
	printRequestInfo(request)
}

func user(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application: json;")
	fmt.Fprintf(writer, "User")
	printRequestInfo(request)
}
