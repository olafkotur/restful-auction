package main

import (
	"database/sql"
	"encoding/json"
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
