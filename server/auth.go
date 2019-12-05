package main

import (
	"crypto/sha256"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

const salt = "N3DxzwcD4zFE"
const expThreshold = 1000 * time.Minute

// Return JWT header using HS256 algorithm
func generateHeader() (h string) {
	header := `{"alg": "HS256", "typ": "JWT"}`
	encoded := b64.StdEncoding.EncodeToString([]byte(header))
	return strings.TrimRight(encoded, "=")
}

// Return JWT payload with user id and username
func generatePayload(username string) (p string) {
	exp := int(time.Now().Add(expThreshold).Unix())
	payload := `{"exp": "` + toString(exp) + `", "username": "` + username + `"}`
	encoded := b64.StdEncoding.EncodeToString([]byte(payload))
	return strings.TrimRight(encoded, "=")
}

// Return JWT signature using SHA256
func generateSignature(secret string) (s string) {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(secret + salt))
	hash := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return strings.TrimRight(hash, "=")
}

// Return complete JWT token with header, payload and signature
func generateToken(user User) (t string) {
	header := generateHeader()
	payload := generatePayload(user.Username)
	return header + "." + payload + "." + user.Password
}

// Checks whether received token is valid
func authenticateToken(token string) (v bool) {
	// Add padding for comparison
	components := strings.Split(token, ".")
	if i := len(components[1]) % 4; i != 0 {
		components[1] += strings.Repeat("=", 4-i)
	}

	// Decode payload data to get user info
	p, _ := b64.StdEncoding.DecodeString(components[1])
	var payload JwtToken
	_ = json.Unmarshal(p, &payload)

	// Check if the session has expired
	if int64(toInt(payload.Exp)) < time.Now().Unix() {
		return false
	}

	// Check if the username and password match
	for _, user := range users {
		if payload.Username == user.Username {
			if user.Password == string(components[2]) {
				return true
			}
		}
	}
	return false
}
