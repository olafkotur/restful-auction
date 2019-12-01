package main

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type PingResponse struct {
	Status string `json:"status"`
	Time   int64  `json:"time"`
}

type Auction struct {
	Id           int     `json:"id"`
	Status       string  `json:"status"`
	Name         string  `json:"name"`
	FirstBid     float64 `json:"firstBid"`
	SellerId     int     `json:"sellerId"`
	ReservePrice float64 `json:"reservePrice"`
}

type Bid struct {
	Id        int     `json:"id"`
	AuctionId int     `json:"auctionId"`
	BidAmount float64 `json:"bidAmount"`
	BidderId  int     `json:"bidderId"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type SyncDataInfo struct {
	Counter int    `json:"counter"`
	Type    string `json:"type"`
	Action  string `json:"action"`
}
