package main

import (
	"net/http"
)

// Returns all users
func getUsers(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)

	// Return all users omitting password
	var res []UserInfo
	for _, user := range users {
		res = append(res, UserInfo{user.Id, user.Username})
	}

	sendResponse(res, writer)
}

// Creates a new user in the system
func createUser(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)
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

	users = append(users, User{userId, username, password})
	sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
}

// Allows the user to login with supplied credentials
func userLogin(writer http.ResponseWriter, request *http.Request) {
	printRequestInfo(request)

	// Extract data from the request body
	_ = request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")

	// Check if the credentials match with a user
	for _, user := range users {
		if username == user.Username {
			if password == user.Password {
				sendResponse(ApiResponse{200, "success", "Successful operation"}, writer)
				return
			}
		}
	}

	sendResponse(ApiResponse{400, "error", "Invalid username/password supplied"}, writer)
}
