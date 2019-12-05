package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func director(r *http.Request) {
	origin, _ := url.Parse(servers[serverId-1].Url)
	r.Header.Add("X-Forwarded-Host", r.Host)
	r.Header.Add("X-Origin-Host", origin.Host)
	r.URL.Scheme = "http"
	r.URL.Host = origin.Host
	fmt.Println("Redirected to:", origin.String())
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	getNextAvailable()

	// Update other servers only if data is modified
	proxy.ServeHTTP(w, r)
	if r.Method == "POST" || r.Method == "DELETE" {
		// Send update request to other servers
		for i := 0; i < maxServers; i++ {
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
}

func getPingResponse(s ServerInfo) (b bool) {
	res, err := http.Get(s.Url + "/ping")
	if err != nil {
		return false
	}

	var data PingResponse
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &data)
	res.Body.Close()

	return data.Status == "pong"
}
