package main

import (
	"log"
	"net/http"
	"net/url"
)

func getAuctions() {
	res, err := http.Get(apiPrefix + "/api/auctions")
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printArrayResponseBody(res)
}

func addAuction(name, firstBid, sellerId, reservePrice string) {
	values := url.Values{
		"name":         {name},
		"firstBid":     {firstBid},
		"sellerId":     {sellerId},
		"reservePrice": {reservePrice},
	}

	res, err := http.PostForm(apiPrefix+"/api/auction", values)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
}

func getAuction(auctionId string) {
	res, err := http.Get(apiPrefix + "/api/auction/" + auctionId)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printResponseBody(res)
}

func updateAuction(auctionId, name, firstBid, sellerId string) {
	values := url.Values{
		"name":     {name},
		"firstBid": {firstBid},
		"sellerId": {sellerId},
	}

	res, err := http.PostForm(apiPrefix+"/api/auction/"+auctionId, values)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printResponseBody(res)
}

func deleteAuction(auctionId string) {
	req, err := http.NewRequest("DELETE", apiPrefix+"/api/auction/"+auctionId, nil)
	if err != nil {
		log.Println(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printResponseBody(res)
}

func addAuctionBid(auctionId, bidAmount, bidderId string) {
	values := url.Values{
		"bidAmount": {bidAmount},
		"bidderId":  {bidderId},
	}

	res, err := http.PostForm(apiPrefix+"/api/auction", values)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printResponseBody(res)
}

func getBidsByAuctionId(auctionId string) {
	res, err := http.Get(apiPrefix + "/api/auction/" + auctionId + "/bids")
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printArrayResponseBody(res)
}

func createUser(username, password string) {
	values := url.Values{
		"username": {username},
		"password": {password},
	}

	res, err := http.PostForm(apiPrefix+"/api/user", values)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printResponseBody(res)
}

func userLogin(username, password string) {
	values := url.Values{
		"username": {username},
		"password": {password},
	}

	res, err := http.PostForm(apiPrefix+"/api/user/login", values)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	printResponseBody(res)
}
