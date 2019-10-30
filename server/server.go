package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Env variables
	var SERVER_PORT string
	err := godotenv.Load()
	if err != nil {
		SERVER_PORT = "8080"
	} else {
		SERVER_PORT = os.Getenv("SERVER_PORT")
	}

	// Server routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/auctions", getAuctions).Methods("GET")
	router.HandleFunc("/api/auction", addAuction).Methods("POST")
	router.HandleFunc("/api/auction/{id}", getAuction).Methods("GET")
	router.HandleFunc("/api/auction/{id}", updateAuction).Methods("POST")
	router.HandleFunc("/api/auction/{id}", deleteAuction).Methods("DELETE")

	fmt.Printf("Listening on port %s...\n\n", SERVER_PORT)
	http.ListenAndServe(":"+SERVER_PORT, router)
}

func getDatabase() (d *sql.DB) {
	db, _ := sql.Open("sqlite3", "./auction.db")
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS auctions (id INTEGER PRIMARY KEY, name TEXT, firstBid REAL, sellerId INTEGER, status TEXT)")
	statement.Exec()
	return db
}

func printRequestInfo(request *http.Request) {
	fmt.Println("Method: ", request.Method)
	fmt.Println("URL: ", request.URL)
	fmt.Println("")
}

func sendResponse(res interface{}, writer http.ResponseWriter) {
	response, _ := json.Marshal(res)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
}

func toInt(s string) (i int) {
	str, _ := strconv.Atoi(s)
	return str
}
