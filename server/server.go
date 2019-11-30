package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var auctions []Auction
var bids []Bid
var users []User

// var client *redis.Client

func main() {
	// Env variables
	// var REDIS_URL string
	// err := godotenv.Load()
	// REDIS_URL = os.Getenv("REDIS_URL")
	// if err != nil {
	// 	REDIS_URL = "localhost:6379"
	// }

	port := "8080"

	// client = redis.NewClient(&redis.Options{
	// 	Addr:     REDIS_URL,
	// 	Password: "",
	// 	DB:       0,
	// })

	// Ensure that the redis db is connected
	// _, err = client.Ping().Result()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	router := mux.NewRouter().StrictSlash(true)

	// Documentaion
	router.HandleFunc("/api/docs", getDocumentation).Methods("GET")

	// Auction
	router.HandleFunc("/api/auctions", getAuctions).Methods("GET")
	router.HandleFunc("/api/auction", addAuction).Methods("POST")
	router.HandleFunc("/api/auction/{auctionId}", getAuction).Methods("GET")
	router.HandleFunc("/api/auction/{auctionId}", updateAuction).Methods("POST")
	router.HandleFunc("/api/auction/{auctionId}", deleteAuction).Methods("DELETE")

	// Bid
	router.HandleFunc("/api/auction/{auctionId}/bid", addAuctionBid).Methods("POST")
	router.HandleFunc("/api/auction/{auctionId}/bids", getBidsByAuctionId).Methods("GET")

	// User
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/user", createUser).Methods("POST")
	router.HandleFunc("/api/user/login", userLogin).Methods("POST")

	fmt.Printf("Listening on port %s...\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getDocumentation(writer http.ResponseWriter, request *http.Request) {
	// Read docs html file
	bytes, err := ioutil.ReadFile("./docs.html")
	if err != nil {
		panic(err)
	}

	// Send response as html
	writer.Header().Set("Content-Type", "text/html")
	_, err = writer.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func printRequestInfo(request *http.Request) {
	fmt.Println("Method: ", request.Method)
	fmt.Println("URL: ", request.URL)
	fmt.Println("")
}

func sendResponse(res interface{}, writer http.ResponseWriter) {
	response, _ := json.Marshal(res)
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write(response)
	if err != nil {
		panic(err)
	}
}

func getMuxVariable(target string, request *http.Request) (v string) {
	return mux.Vars(request)[target]
}

func assignAuctionId() (key int) {
	var highestKey int
	for _, auction := range auctions {
		if auction.Id > highestKey {
			highestKey = auction.Id
		}
	}
	return highestKey + 1
}

func assignBidId() (key int) {
	var highestKey int
	for _, bid := range bids {
		if bid.Id > highestKey {
			highestKey = bid.Id
		}
	}
	return highestKey + 1
}

func assignUserId() (key int) {
	var highestKey int
	for _, user := range users {
		if user.Id > highestKey {
			highestKey = user.Id
		}
	}
	return highestKey + 1
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
