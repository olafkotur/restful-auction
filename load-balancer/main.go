package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/go-redis/redis/v7"
)

const DEBUG = true
const MAX_SERVERS = 2

var client *redis.Client
var servers []ServerInfo
var serverId int

func main() {
	// Get redis url from environment variables
	var redisUrl string
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl = redisHost + ":" + redisPort

	// Testing only
	if DEBUG {
		redisUrl = "localhost:6379"
	}

	// Create new instance of redis
	client = redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	// Get active server info
	serverId = 1

	// Handle reverse proxy
	proxy := &httputil.ReverseProxy{Director: director}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		servers = getServerInfo()

		// Ensure that target server is active before sending a request
		alive := getPingResponse(servers[serverId-1])
		if !alive {
			fmt.Printf("Target server with id %d is inactive, switching servers\n", serverId)
			changeTargetServer()
		}
		proxy.ServeHTTP(w, r)
	})

	fmt.Printf("Listening on port 8080...\n\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func director(r *http.Request) {
	origin, _ := url.Parse(servers[serverId-1].Url)
	fmt.Println(origin)
	r.Header.Add("X-Forwarded-Host", r.Host)
	r.Header.Add("X-Origin-Host", origin.Host)
	r.URL.Scheme = "http"
	r.URL.Host = origin.Host
	fmt.Println("Redirected to:", origin.String())
	changeTargetServer()
}

func changeTargetServer() {
	if serverId >= MAX_SERVERS {
		serverId = 1
	} else {
		serverId++
	}

	fmt.Printf("Set next target server to server: %d\n\n", serverId)
}

func getServerInfo() (in []ServerInfo) {
	var info []ServerInfo
	s1, _ := client.Get("server1").Result()
	s2, _ := client.Get("server2").Result()
	info = append(info, ServerInfo{1, s1})
	info = append(info, ServerInfo{2, s2})
	return info
}

func getPingResponse(s ServerInfo) (alive bool) {
	res, err := http.Get(s.Url + "/api/ping")
	if err != nil {
		fmt.Println("No response from server: ", s.Url)
		return false
	}

	var data PingResponse
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Bad response from server: ", s.Url)
		return false
	}
	_ = json.Unmarshal(body, &data)

	if data.Status != "pong" {
		fmt.Println("Invalid response from server: ", s.Url)
		return false
	}

	return true
}
