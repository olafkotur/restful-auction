package main

type Auction struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	FirstBid float64 `json:"firstbid"`
	SellerId int     `json:"sellerId"`
	Status   string  `json:"status"`
}
