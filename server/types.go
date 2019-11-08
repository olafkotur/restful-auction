package main

type ApiResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Auction struct {
	Id       int     `json:"id"`
	Status   string  `json:"status"`
	Name     string  `json:"name"`
	FirstBid float64 `json:"firstBid"`
	SellerId int     `json:"sellerId"`
}

type Bid struct {
	Id        int     `json:"id"`
	AuctionId int     `json:"auctionId"`
	BidAmount float64 `json:"bidAmount"`
	BidderId  int     `json:"bidderId"`
}

// type User struct {
// 	Id       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
