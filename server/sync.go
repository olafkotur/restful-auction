package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func setSyncData(typ, action string, d interface{}) {
	counter++
	info, _ := json.Marshal(SyncDataInfo{counter, typ, action})
	data, _ := json.Marshal(d)
	client.Set("syncInfo", info, 0)
	client.Set("lastRequest", data, 0)
}

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
		handleAuction(syncInfo.Action, data)
		counter = syncInfo.Counter

	} else if syncInfo.Type == "bids" {
		var data Bid
		_ = json.Unmarshal([]byte(raw), &data)
		handleBid(syncInfo.Action, data)
		counter = syncInfo.Counter

	} else if syncInfo.Type == "users" {
		var data User
		_ = json.Unmarshal([]byte(raw), &data)
		handleUser(syncInfo.Action, data)
		counter = syncInfo.Counter
	}
}

func handleAuction(action string, data Auction) {
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

func handleBid(action string, data Bid) {
	fmt.Println("Handling auction data with action:", action)
	if action == "add" {
		bids = append(bids, Bid{data.Id, data.AuctionId, data.BidAmount, data.BidderId})
	}
}

func handleUser(action string, data User) {
	fmt.Println("Handling auction data with action:", action)
	if action == "add" {
		users = append(users, User{data.Id, data.Username, data.Password})
	}
}
