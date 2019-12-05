package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

const DEBUG = false

var client *redis.Client
var counter int
var serverId = 1

var auctions []Auction
var bids []Bid
var users []User

func main() {
	redisUrl := "redis:6379"
	host, _ := os.Hostname()
	port := "8080"

	// Create new instance of redis and ensure connection
	client = redis.NewClient(&redis.Options{Addr: redisUrl})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	// Assign server id
	keys, _ := client.Keys("server:*").Result()
	for _, key := range keys {
		id := toInt(strings.Split(key, ":")[1])
		if id >= serverId {
			serverId = id + 1
		}
	}

	if DEBUG { // DEBUG: Testing only
		port = toString(8080 + serverId)
	}

	// Add port if it does not exist
	if !strings.Contains(host, ":") {
		host += ":" + port
	}

	// Set server information
	client.Set("server:"+toString(serverId), "http://"+host, 0)
	attemptDataRecovery()

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
	router.HandleFunc("/ping", getPingResponse).Methods("GET")
	router.HandleFunc("/sync", getSyncData).Methods("GET")
	router.HandleFunc("/recover", getRecoveryData).Methods("GET")

	// Attempt to begin serving
	fmt.Printf("Listening at: %s\n", host)
	_ = http.ListenAndServe(":"+port, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		printRequestInfo(request)
		router.ServeHTTP(writer, request)
	}))
}

// Redirect user to generated postman documentation
func getDocumentation(writer http.ResponseWriter, request *http.Request) {
	url := "https://documenter.getpostman.com/view/8555555/SWDzgMiw?version=latest"
	http.Redirect(writer, request, url, http.StatusSeeOther)
}

// Respond to a ping request
func getPingResponse(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)

	res := PingResponse{"pong", time.Now().Unix()}
	sendResponse(res, writer)
}
