package main

import (
	"net/http"
)

// Returns all users
func getUsers(writer http.ResponseWriter, request *http.Request) {
	// Return all users omitting password
	var res []UserInfo
	for _, user := range users {
		res = append(res, UserInfo{user.Id, user.Username})
	}

	sendResponse(res, writer)
}

// Creates a new user in the system
func createUser(writer http.ResponseWriter, request *http.Request) {
	userId := assignUserId()

	// Extract data from the request body
	_ = request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// Check if the username is already taken
	for _, user := range users {
		if username == user.Username {
			sendResponse(ApiResponse{400, "error", "Invalid username/password supplied"}, writer)
			return
		}
	}

	// Add user to the user data and update redis counter
	user := User{userId, username, password}
	users = append(users, user)
	setSyncData("users", "add", user)

	token := generateToken(user)
	sendAuthResponse(token, writer)
}

// Allows the user to login with supplied credentials
func userLogin(writer http.ResponseWriter, request *http.Request) {
	// Extract data from the request body
	_ = request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// Check if the credentials match with a user
	for _, user := range users {
		if username == user.Username {
			if password == user.Password {
				token := generateToken(user)
				sendAuthResponse(token, writer)
				return
			}
		}
	}

	sendResponse(ApiResponse{400, "error", "Invalid username/password supplied"}, writer)
}
