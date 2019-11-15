package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getUsers(writer http.ResponseWriter, request *http.Request) {
	users := []UserInfo{}

	// Fetch all users and add them to the result array
	keys := client.Keys("user:*").Val()
	for _, key := range keys {
		userData := UserInfo{}
		user, _ := client.Get(key).Result()
		_ = json.Unmarshal([]byte(user), &userData)
		users = append(users, userData)
	}

	sendResponse(users, writer)
	printRequestInfo(request)
}

func createUser(writer http.ResponseWriter, request *http.Request) {
	userId := assignKeyId("user")

	// Extract data from the request body
	_ = request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// Check if user already exists with that name
	keys := client.Keys("user:*").Val()
	for _, key := range keys {
		userData := User{}
		user, _ := client.Get(key).Result()
		_ = json.Unmarshal([]byte(user), &userData)
		if userData.Username == username {
			sendResponse(ApiResponse{400, "error", "Invalid username/password supplied"}, writer)
			return
		}
	}

	// Create user in the database
	item, _ := json.Marshal(User{userId, username, password})
	client.Set("user:"+toString(userId), item, 0)

	sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
	printRequestInfo(request)
}

func userLogin(writer http.ResponseWriter, request *http.Request) {
	// Extract data from the request body
	_ = request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// Check if the credentials match with a user
	keys := client.Keys("user:*").Val()
	isLoggedIn := false
	for _, key := range keys {
		userData := User{}
		user, _ := client.Get(key).Result()
		_ = json.Unmarshal([]byte(user), &userData)
		if userData.Username == username && userData.Password == password {
			fmt.Println(userData.Username, username)
			isLoggedIn = true
		}
	}

	if isLoggedIn {
		sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
	} else {
		sendResponse(ApiResponse{400, "error", "Invalid username/password supplied"}, writer)
	}

	printRequestInfo(request)
}
