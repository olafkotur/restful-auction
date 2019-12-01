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
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

const DEBUG = false

var MAX_SERVERS int
var client *redis.Client
var servers []ServerInfo
var serverId int

func main() {
	// Get environment variables
	var redisUrl string
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisUrl = redisHost + ":" + redisPort
	MAX_SERVERS = toInt(os.Getenv("MAX_SERVERS"))

	// Testing only
	if DEBUG {
		redisUrl = "localhost:6379"
		MAX_SERVERS = 2
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
		updateServerInfo()

		// Ensure that target server is active before sending a request
		alive := false
		for !alive {
			alive = getPingResponse(servers[serverId-1])
			if !alive {
				fmt.Printf("Target server with id %d is inactive, switching servers\n", serverId)
				changeTargetServer()
				time.Sleep(100 * time.Millisecond)
			}
		}

		// Update other servers only if data is modified
		proxy.ServeHTTP(w, r)
		if r.Method == "POST" || r.Method == "DELETE" {
			// Send update request to other servers
			for i := 0; i < MAX_SERVERS; i++ {
				if serverId == i+1 {
					fmt.Println("Skipping updating server:", serverId)
					continue
				}
				fmt.Println("Updating server", servers[i].Url+"/sync")
				res, err := http.Get(servers[i].Url + "/sync")
				if err != nil {
					continue
				}
				res.Body.Close()
			}
		}

		changeTargetServer()
	})

	fmt.Printf("Listening on port 8080...\n\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func director(r *http.Request) {
	origin, _ := url.Parse(servers[serverId-1].Url)
	r.Header.Add("X-Forwarded-Host", r.Host)
	r.Header.Add("X-Origin-Host", origin.Host)
	r.URL.Scheme = "http"
	r.URL.Host = origin.Host
	fmt.Println("Redirected to:", origin.String())
}

func changeTargetServer() {
	if serverId >= MAX_SERVERS {
		serverId = 1
	} else {
		serverId++
	}
	fmt.Printf("Set next target server to server: %d\n\n", serverId)
}

func updateServerInfo() {
	for i := 0; i < MAX_SERVERS; i++ {
		url, _ := client.Get("server" + toString(i+1)).Result()

		// Update only if there is a change detected or a new server is discovered
		if len(servers) >= i+1 && url != servers[i].Url {
			fmt.Println("Detected change in server, updating:", url)
			servers[i] = ServerInfo{i + 1, url}
		} else {
			fmt.Println("Discovered new server, adding:", url)
			servers = append(servers, ServerInfo{i + 1, url})
		}
	}
}

func getPingResponse(s ServerInfo) (alive bool) {
	fmt.Println("Pinging:", s.Url)
	res, err := http.Get(s.Url + "/ping")
	if err != nil {
		fmt.Println("No response from server: ", s.Url)
		return false
	}

	var data PingResponse
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &data)
	res.Body.Close()

	if data.Status != "pong" {
		fmt.Println("Invalid response from server: ", s.Url)
		return false
	}

	return true
}

func toInt(s string) (i int) {
	res, _ := strconv.Atoi(s)
	return res
}

func toString(i int) (s string) {
	return strconv.Itoa(i)
}
