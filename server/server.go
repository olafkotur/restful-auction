package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}

	// Enviornment variables
	SERVER_PORT := os.Getenv("SERVER_PORT")

	// Server routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/auctions", auctions)
	router.HandleFunc("/api/auction", auction)
	router.HandleFunc("/api/user", user)
	router.HandleFunc("/api/user/login", login)

	http.ListenAndServe(":"+SERVER_PORT, router)
}

func printRequestInfo(request *http.Request) {
	fmt.Println("Method: ", request.Method)
	fmt.Println("URL: ", request.URL)
	fmt.Println("")
}

func auctions(writer http.ResponseWriter, request *http.Request) {

	type AuctionsResponse struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	res := AuctionsResponse{
		42,
		"Hello",
		"World",
	}

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(writer, err.Error(), 500)
		return
	}

	printRequestInfo(request)
	writer.Header().Set("Content-Type", "application: json")
	writer.Write(response)
}

func auction(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)
	writer.Header().Set("Content-Type", "application: json;")
	fmt.Fprintf(writer, "Auction")
}

func user(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)
	writer.Header().Set("Content-Type", "application: json;")
	fmt.Fprintf(writer, "User")
}

func login(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)
	writer.Header().Set("Content-Type", "application: json;")
	fmt.Fprintf(writer, "Login")
}
