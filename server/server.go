package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// curl -d "id=1&name=bob&firstBid=12.2&sellerId=32&status=Testing1" localhost:8080/api/auction
// curl -d "id=2&name=phil&firstBid=7329&sellerId=12&status=Testing2" localhost:8080/api/auction
// curl -d "id=3&name=sandra&firstBid=12312&sellerId=2&status=Testing3" localhost:8080/api/auction
// curl -d "name=bobby&status=newmessage" localhost:8080/api/auction

var client *redis.Client

func main() {
	// Env variables
	var SERVER_PORT string
	err := godotenv.Load()
	if err != nil {
		SERVER_PORT = "8080"
	} else {
		SERVER_PORT = os.Getenv("SERVER_PORT")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Ensure that the redis db is connected
	_, err = client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	// Auction
	router.HandleFunc("/api/auctions", getAuctions).Methods("GET")
	router.HandleFunc("/api/auction", addAuction).Methods("POST")
	router.HandleFunc("/api/auction/{auctionId}", getAuction).Methods("GET")
	router.HandleFunc("/api/auction/{auctionId}", updateAuction).Methods("POST")
	router.HandleFunc("/api/auction/{auctionId}", deleteAuction).Methods("DELETE")

	// Bid
	// router.HandleFunc("/api/auction/{id}")

	fmt.Printf("Listening on port %s...\n\n", SERVER_PORT)
	http.ListenAndServe(":"+SERVER_PORT, router)
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

func sendSuccessResponse(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("Success\n"))
}

func sendFailedResponse(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("Failed\n"))
}

func getMuxVariable(target string, request *http.Request) (v string) {
	return mux.Vars(request)[target]
}

func toString(i int) (s string) {
	return strconv.Itoa(i)
}

func toInt(s string) (i int) {
	res, _ := strconv.Atoi(s)
	return res
}

func toFloat(s string) (f float64) {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}
