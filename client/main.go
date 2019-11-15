package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

/*
	1. Get user input on which end point they would like to target
	2. Get user input for each data point they need to enter
	3. Make the request
	4. Print result
*/

/* Seller endpoints
1. POST /api/auction
2. UPDATE /api/auction/{auctionId}
3. DELETE /api/auction/{auctionId}
4. POST /api/user
5. POST /api/user/login
*/

var apiPrefix string

func main() {
	env := os.Getenv("ENV")
	if env == "production" {
		apiPrefix = "TODO"
	} else {
		apiPrefix = "http://localhost:8080"
	}

	// clientType := os.Getenv("CLIENT_TYPE")
	clientType := "seller"

	if clientType == "seller" {
		clientSell()
	} else if clientType == "buyer" {

	} else {
		fmt.Println("You must specify whether the client is a 'seller' or a 'buyer'")
	}
}

func clientSell() {
	sellerOptions := []string{
		"List all active auctions",
		"Create a new auction",
		"Get auction by id",
		"Update an exisiting auction",
		"Delete an exisiting auction",
		"Create a new user",
		"Login as existing user",
	}

	fmt.Printf("------Available seller options-------\n\n")
	for i, opt := range sellerOptions {
		fmt.Printf("%d: %s\n", i+1, opt)
	}
	fmt.Println()

	// TODO: Get user input
	option := 9
	id := "3"
	name, firstBid, sellerId, reservePrice := "bob", "1", "2", "0"
	bidAmount, bidderId := "52", "1"
	username, password := "olafciu", "lapis"

	switch option {
	case 1:
		getAuctions()
	case 2:
		addAuction(name, firstBid, sellerId, reservePrice)
	case 3:
		getAuction(id)
	case 4:
		updateAuction(id, name, firstBid, sellerId)
	case 5:
		deleteAuction(id)
	case 6:
		addAuctionBid(id, bidAmount, bidderId)
	case 7:
		getBidsByAuctionId(id)
	case 8:
		createUser(username, password)
	case 9:
		userLogin(username, password)
	}
}

func printArrayResponseBody(res *http.Response) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
	}

	for _, d := range data {
		keys := reflect.ValueOf(d).MapKeys()
		for _, key := range keys {
			fmt.Print(key.String(), ": ", d[key.String()], " ")
		}
		fmt.Println()
	}
}

func printResponseBody(res *http.Response) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
	}

	keys := reflect.ValueOf(data).MapKeys()
	for _, key := range keys {
		fmt.Print(key.String(), ": ", data[key.String()], " ")
	}
	fmt.Println()
}
