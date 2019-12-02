package main

import (
	"crypto/sha256"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/json"
	"strings"
)

// Return JWT header using HS256 algorithm
func generateHeader() (h string) {
	header := `{"alg": "HS256", "typ": "JWT"}`
	encoded := b64.StdEncoding.EncodeToString([]byte(header))
	return encoded
}

// Return JWT payload with user id and username
func generatePayload(userId, username string) (p string) {
	payload := `{"id": "` + userId + `", "username": "` + username + `"}`
	encoded := b64.StdEncoding.EncodeToString([]byte(payload))
	return encoded
}

// Return JWT signature using SHA256
func generateSignature(secret string) (s string) {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(secret))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// Return complete JWT token with header, payload and signature
func generateToken(user User) (t string) {
	header := generateHeader()
	payload := generatePayload(toString(user.Id), user.Username)
	signature := generateSignature(user.Password)
	return header + "." + payload + "." + signature
}

// Checks whether received token is valid
func authenticateToken(token string) (v bool) {
	components := strings.Split(token, ".")
	p, _ := b64.StdEncoding.DecodeString(components[1])
	var payload JwtToken
	_ = json.Unmarshal(p, &payload)

	// Check if the username and password match
	for _, user := range users {
		if payload.Username == user.Username {
			secret := generateSignature(user.Password)
			if secret == string(components[2]) {
				return true
			}
		}
	}
	return false
}
