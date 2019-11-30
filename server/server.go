package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

var auctions []Auction
var bids []Bid
var users []User

var client *redis.Client

func main() {
	// Get redis url from environment variables
	var redisUrl string
	externalPort := os.Getenv("EXTERNAL_PORT")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl = redisHost + ":" + redisPort
	if redisUrl == ":" {
		redisUrl = "redis:6379"
	}

	// Create new instance of redis
	client = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	// Ensure that the redis db is connected
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

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

	fmt.Printf("Listening on port %s...\n\n", externalPort)
	log.Fatal(http.ListenAndServe(":8080", router))
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
