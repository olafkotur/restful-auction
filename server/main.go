package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

const DEBUG = false

var auctions []Auction
var bids []Bid
var users []User

var counter int
var client *redis.Client

func main() {
	var redisUrl string
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl = redisHost + ":" + redisPort

	// Get server information
	port := "8080"
	serverId := os.Getenv("SERVER_ID")
	url := "http://server" + serverId + ":" + port

	// DANGER: Debug use only
	if DEBUG {
		serverId = "1"
		port = toString(8080 + toInt(serverId))
		redisUrl = "localhost:6379"
		url = "http://localhost:" + port
	}

	// Create new instance of redis and ensure connection
	client = redis.NewClient(&redis.Options{Addr: redisUrl})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	// Add server information to the db
	client.Set("server"+serverId, url, 0)

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

	// Attempt to begin serving
	fmt.Printf("Listening on %s...\n\n", url)
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

// Outputs methods and url of a request
func printRequestInfo(request *http.Request) {
	if request.URL.RequestURI() == "/ping" {
		return
	}
	fmt.Println("Method: ", request.Method)
	fmt.Println("URL: ", request.URL)
}

// Sends any interface using application JSON format
func sendResponse(res interface{}, writer http.ResponseWriter) {
	response, _ := json.Marshal(res)
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write(response)
	if err != nil {
		panic(err)
	}
}

// Sends JWT auth token using application text format
func sendAuthResponse(token string, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/text")
	_, err := writer.Write([]byte(token))
	if err != nil {
		panic(err)
	}
}

// Returns requested mux variable from request URL
func getMuxVariable(target string, request *http.Request) (v string) {
	return mux.Vars(request)[target]
}

// Assigns next available unqiue auction Id
func assignAuctionId() (key int) {
	var highestKey int
	for _, auction := range auctions {
		if auction.Id > highestKey {
			highestKey = auction.Id
		}
	}
	return highestKey + 1
}

// Assigns next available unqiue bid Id
func assignBidId() (key int) {
	var highestKey int
	for _, bid := range bids {
		if bid.Id > highestKey {
			highestKey = bid.Id
		}
	}
	return highestKey + 1
}

// Assigns next available unqiue user Id
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
