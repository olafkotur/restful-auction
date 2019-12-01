package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

const DEBUG = true

var auctions []Auction
var bids []Bid
var users []User

var client *redis.Client

func main() {
	// Get redis url from environment variables
	var redisUrl string
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl = redisHost + ":" + redisPort

	// Get server information
	serverId := os.Getenv("SERVER_ID")
	port := toString(8080 + toInt(serverId))
	url := "http://server" + serverId + ":8080"

	// Testing only
	if DEBUG {
		serverId = "2"
		port = "8082"
		redisUrl = "localhost:6379"
		url = "http://localhost:" + port
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

	// Add server information to the db
	_ = client.Set("server"+serverId, url, 0)

	router := mux.NewRouter().StrictSlash(true)

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

	// Misc
	router.HandleFunc("/api/docs", getDocumentation).Methods("GET")
	router.HandleFunc("/api/ping", getPingResponse).Methods("GET")

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

func getPingResponse(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)

	res := PingResponse{"pong", time.Now().Unix()}
	sendResponse(res, writer)
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
