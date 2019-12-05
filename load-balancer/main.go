package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/go-redis/redis/v7"
)

var proxy *httputil.ReverseProxy
var client *redis.Client
var servers []ServerInfo
var ignored []ServerInfo
var maxServers int
var serverId = 1

func main() {
	redisUrl := "redis:6379"
	port := "8080"

	// Create new instance of redis and ensure connection
	client = redis.NewClient(&redis.Options{Addr: redisUrl})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	// Attempt to add old servers into the network
	go checkIgnoredStatus()

	// Handle reverse proxy
	proxy = &httputil.ReverseProxy{Director: director}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		updateServerInfo()

		// Dont attempt to process query if no servers are able to respond
		if maxServers < 1 {
			fmt.Println("There are no current active servers available")
			return
		}

		handleRedirect(w, r)
	})

	fmt.Printf("Listening on port %s...\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
