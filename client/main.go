package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
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
	apiPrefix = os.Getenv("API_PREFIX")
	if apiPrefix == "" {
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

	fmt.Printf("------Available seller options-------\n")
	for i, opt := range sellerOptions {
		fmt.Printf("%d: %s\n", i+1, opt)
	}
	fmt.Println()

	option, _ := strconv.Atoi(getUserInput("option"))

	switch option {
	case 1:
		getAuctions()
	case 2:
		name := getUserInput("name")
		firstBid := getUserInput("firstBid")
		sellerId := getUserInput("sellerId")
		reservePrice := getUserInput("reservePrice")
		addAuction(name, firstBid, sellerId, reservePrice)
	case 3:
		id := getUserInput("id")
		getAuction(id)
	case 4:
		id := getUserInput("id")
		name := getUserInput("name")
		firstBid := getUserInput("firstBid")
		sellerId := getUserInput("sellerId")
		updateAuction(id, name, firstBid, sellerId)
	case 5:
		id := getUserInput("id")
		deleteAuction(id)
	case 6:
		username := getUserInput("username")
		password := getUserInput("password")
		createUser(username, password)
	case 7:
		username := getUserInput("username")
		password := getUserInput("password")
		userLogin(username, password)
	}

	fmt.Println("\nWould you like to request again? y/n")
	again := getUserInput("option")
	if again == "y" || again == "Y" {
		clientSell()
	} else {
		return
	}
}

func getUserInput(inputType string) (i string) {
	fmt.Println("Please enter the " + inputType)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.Replace(input, "\r\n", "", -1)
	input = strings.Replace(input, "\n", "", -1)
	return input
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
