package main

type ServerInfo struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type PingResponse struct {
	Status string `json:"status"`
	Time   int64  `json:"time"`
}
