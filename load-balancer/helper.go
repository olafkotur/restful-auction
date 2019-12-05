package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const ignoredTimeout = 60 * time.Second

func updateServerInfo() {
	// Set max value the available servers
	keys, _ := client.Keys("server:*").Result()
	maxServers = len(keys)

	// Update only if there is a change detected or a new server is discovered
	for i := 0; i < maxServers; i++ {
		url, _ := client.Get("server:" + toString(i+1)).Result()
		if !serverExists(url) {
			fmt.Println("Discovered new server, adding:", url)
			servers = append(servers, ServerInfo{i + 1, url})
		}
	}
	fmt.Println()
}

// Changes next request target to the best available server using round robin
func changeTargetServer() {
	if serverId >= maxServers {
		serverId = 1
	} else {
		serverId++
	}

	// Keep going until an active server is found
	if isIgnored() {
		fmt.Printf("Target server %d has been marked as ignored: %s\n", serverId, servers[serverId-1].Url)
		changeTargetServer()
	} else {
		fmt.Printf("Set next target server to server: %d\n\n", serverId)
	}
}

// Checks if the requested server is alive
func getNextAvailable() (b bool) {
	if !getPingResponse(servers[serverId-1]) {
		log.Printf("No response from server %s, added to ignore list", servers[serverId-1].Url)
		ignored = append(ignored, servers[serverId-1])

		// Try again
		time.Sleep(500 * time.Millisecond)
		changeTargetServer()
		getNextAvailable()
	}
	return true
}

// Attempt to allow a server to join the network again
func checkIgnoredStatus() {
	for {
		time.Sleep(ignoredTimeout)
		for _, server := range ignored {
			if getPingResponse(server) {
				index := server.Id - 1
				ignored = append(ignored[:index], ignored[index+1:]...)
				log.Printf("Previously ignored server %s is now back online, removed from ignored\n\n", server.Url)
			}
		}
	}
}

// Checks whether the server in question should be ignored
func isIgnored() (b bool) {
	for _, server := range ignored {
		if servers[serverId-1].Id == server.Id {
			return true
		}
	}
	return false
}

// Checks whether the server exists in the list of known servers
func serverExists(url string) (ex bool) {
	for _, server := range servers {
		if url == server.Url {
			return true
		}
	}
	return false
}

func toString(i int) (s string) {
	return strconv.Itoa(i)
}
