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

var apiPrefix string

func main() {
	apiPrefix = os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "http://localhost:8080"
	}

	fmt.Println("Are you a buyer or a seller?")
	clientType := getUserInput("client type")

	if clientType == "seller" {
		clientSell()
	} else if clientType == "buyer" {
		clientBuy()
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
		id := getUserInput("auctionId")
		getAuction(id)
	case 4:
		id := getUserInput("auctionId")
		name := getUserInput("name")
		firstBid := getUserInput("firstBid")
		sellerId := getUserInput("sellerId")
		updateAuction(id, name, firstBid, sellerId)
	case 5:
		id := getUserInput("auctionId")
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

func clientBuy() {
	sellerOptions := []string{
		"List all active auctions",
		"Get auction by id",
		"Place a bid or an auction",
		"List all bids for a particular auction",
		"Create a new user",
		"Login as existing user",
	}

	fmt.Printf("------Available buyer options-------\n")
	for i, opt := range sellerOptions {
		fmt.Printf("%d: %s\n", i+1, opt)
	}
	fmt.Println()

	option, _ := strconv.Atoi(getUserInput("option"))

	switch option {
	case 1:
		getAuctions()
	case 2:
		id := getUserInput("auctionId")
		getAuction(id)
	case 3:
		id := getUserInput("auctionId")
		bidAmount := getUserInput("bidAmount")
		bidderId := getUserInput("bidderId")
		addAuctionBid(id, bidAmount, bidderId)
	case 4:
		id := getUserInput("auctionId")
		getBidsByAuctionId(id)
	case 5:
		username := getUserInput("username")
		password := getUserInput("password")
		createUser(username, password)
	case 6:
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
