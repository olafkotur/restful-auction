package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Outputs methods and url of a request
func printRequestInfo(request *http.Request) {
	if request.URL.RequestURI() == "/ping" {
		return
	}
	log.Println("Method: ", request.Method)
	log.Println("URL: ", request.URL)
	fmt.Println()
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
