package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Updates redis DB with most recent updated data and increases counter
func setSyncData(typ, action string, d interface{}) {
	// Get previous sync info
	raw, _ := client.Get("syncInfo").Result()
	var syncInfo SyncDataInfo
	_ = json.Unmarshal([]byte(raw), &syncInfo)

	// Set new sync info
	counter++
	fmt.Println("Setting:", counter, typ, action, serverId, syncInfo.PrimaryServerId)
	info, _ := json.Marshal(SyncDataInfo{counter, typ, action, serverId, syncInfo.PrimaryServerId})
	data, _ := json.Marshal(d)
	client.Set("syncInfo", info, 0)
	client.Set("lastRequest", data, 0)
}

// Attempts to read most recent data from redis and update self counter
func getSyncData(writer http.ResponseWriter, request *http.Request) {
	raw, _ := client.Get("syncInfo").Result()
	var syncInfo SyncDataInfo
	_ = json.Unmarshal([]byte(raw), &syncInfo)

	// Check if server is up to date
	fmt.Println("Saved counter:", counter, "Actual counter:", syncInfo.Counter)
	if counter >= syncInfo.Counter {
		fmt.Println("Data is up to date, skipping sync")
		return
	}
	// Set the data update to the correct destination
	raw, _ = client.Get("lastRequest").Result()
	if syncInfo.Type == "auctions" {
		var data Auction
		_ = json.Unmarshal([]byte(raw), &data)
		handleSyncAuction(syncInfo.Action, data)
		counter = syncInfo.Counter

	} else if syncInfo.Type == "bids" {
		var data Bid
		_ = json.Unmarshal([]byte(raw), &data)
		handleSyncBid(syncInfo.Action, data)
		counter = syncInfo.Counter

	} else if syncInfo.Type == "users" {
		var data User
		_ = json.Unmarshal([]byte(raw), &data)
		handleSyncUser(syncInfo.Action, data)
		counter = syncInfo.Counter
	}
}

// Updates auction data based on most recent data in redis
func handleSyncAuction(action string, data Auction) {
	fmt.Println("Handling auction data with action:", action)
	if action == "add" {
		auctions = append(auctions, Auction{data.Id, data.Status, data.Name, data.FirstBid, data.SellerId, data.ReservePrice})
	} else if action == "remove" {
		for i, auction := range auctions {
			if data.Id == auction.Id {
				auctions = append(auctions[:i], auctions[i+1:]...)
			}
		}
	} else if action == "update" {
		for i, auction := range auctions {
			if data.Id == auction.Id {
				auctions[i] = Auction{auction.Id, auction.Status, data.Name, data.FirstBid, data.SellerId, auction.ReservePrice}
			}
		}
	}
}

// Updates bid data based on most recent data in redis
func handleSyncBid(action string, data Bid) {
	fmt.Println("Handling auction data with action:", action)
	if action == "add" {
		bids = append(bids, Bid{data.Id, data.AuctionId, data.BidAmount, data.BidderId})
	}
}

// Updates user data based on most recent data in redis
func handleSyncUser(action string, data User) {
	fmt.Println("Handling auction data with action:", action)
	if action == "add" {
		users = append(users, User{data.Id, data.Username, data.Password})
	}
}

// Sends copy of server data
func getRecoveryData(writer http.ResponseWriter, request *http.Request) {
	res := RecoveryData{auctions, bids, users}
	sendResponse(res, writer)
}

func attemptDataRecovery() {
	raw, _ := client.Get("syncInfo").Result()
	var syncInfo SyncDataInfo
	_ = json.Unmarshal([]byte(raw), &syncInfo)

	// Check if server is up to date
	if counter >= syncInfo.Counter {
		return
	}
	fmt.Println("Data loss detected, attempting to recover")

	// Use second best server if the primary is current
	recoverId := syncInfo.PrimaryServerId
	fmt.Println(recoverId, serverId)
	if recoverId == serverId {
		recoverId = syncInfo.SecondaryServerId
		fmt.Println("Using second best recovery server")
	}

	// Attempt to get data from the most up to date server
	url, _ := client.Get("server:" + recoverId).Result()
	res, err := http.Get(url + "/recover")
	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not recover from data loss, aborting")
	}

	// Save all data
	var data RecoveryData
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &data)
	auctions = data.Auctions
	bids = data.Bids
	users = data.Users

	fmt.Println("Successfully recovered from data loss, qualified to serve")
	res.Body.Close()
}
