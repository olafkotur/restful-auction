package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Env variables
	var SERVER_PORT string
	err := godotenv.Load()
	if err != nil {
		SERVER_PORT = "8080"
	} else {
		SERVER_PORT = os.Getenv("SERVER_PORT")
	}

	// Server routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/auctions", auctions)
	router.HandleFunc("/api/auction", auction)
	router.HandleFunc("/api/user", user)

	fmt.Printf("Listening on port %s...\n\n", SERVER_PORT)
	http.ListenAndServe(":"+SERVER_PORT, router)
}

func printRequestInfo(request *http.Request) {
	fmt.Println("Method: ", request.Method)
	fmt.Println("URL: ", request.URL)
	fmt.Println("")
}
