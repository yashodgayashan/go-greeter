/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Constants for common strings
const (
	DefaultName = "Stranger"
	AppVersion  = "1.0.0"
)

// UserInfo represents user information structure
type UserInfo struct {
	Name     string `json:"name"`
	Age      int    `json:"age,omitempty"`
	Location string `json:"location,omitempty"`
	Email    string `json:"email,omitempty"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

func main() {
	serverMux := http.NewServeMux()

	// Existing endpoint
	serverMux.HandleFunc("/greeter/greet", greet)

	// New endpoints
	serverMux.HandleFunc("/greeter/farewell", farewell)
	serverMux.HandleFunc("/greeter/health", healthCheck)
	serverMux.HandleFunc("/greeter/time-greet", timeBasedGreet)
	serverMux.HandleFunc("/greeter/user-info", userInfoHandler)
	serverMux.HandleFunc("/greeter/bulk-greet", bulkGreet)

	serverPort := 9090
	server := http.Server{
		Addr:              fmt.Sprintf(":%d", serverPort),
		Handler:           serverMux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	go func() {
		log.Printf("Starting HTTP Greeter on port %d\n", serverPort)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP ListenAndServe error: %v", err)
		}
		log.Println("HTTP server stopped serving new requests.")
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh // Wait for shutdown signal

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown error: %v", err)
		return
	}
	log.Println("Shutdown complete.")
}

func greet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = DefaultName
	}
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

// farewell handles goodbye messages
func farewell(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = DefaultName
	}
	fmt.Fprintf(w, "Goodbye, %s! Have a great day!\n", name)
}

// healthCheck provides service health status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	health := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   AppVersion,
	}

	if err := json.NewEncoder(w).Encode(health); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// timeBasedGreet provides time-appropriate greetings
func timeBasedGreet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = DefaultName
	}

	hour := time.Now().Hour()
	var greeting string

	switch {
	case hour < 12:
		greeting = "Good morning"
	case hour < 17:
		greeting = "Good afternoon"
	default:
		greeting = "Good evening"
	}

	fmt.Fprintf(w, "%s, %s!\n", greeting, name)
}

// userInfoHandler handles user information (GET and POST)
func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUserInfo(w, r)
	case http.MethodPost:
		createUserInfo(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed\n", r.Method)
	}
}

// getUserInfo returns user information from query parameters
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "name parameter is required"}); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	user := UserInfo{Name: name}

	if ageStr := r.URL.Query().Get("age"); ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil && age > 0 {
			user.Age = age
		}
	}

	if location := r.URL.Query().Get("location"); location != "" {
		user.Location = location
	}

	if email := r.URL.Query().Get("email"); email != "" {
		user.Email = email
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// createUserInfo creates user information from JSON payload
func createUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user UserInfo
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if encErr := json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"}); encErr != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	if user.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "name field is required"}); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Simulate user creation success
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": fmt.Sprintf("User %s created successfully", user.Name),
		"user":    user,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// bulkGreet handles multiple names at once
func bulkGreet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	namesParam := r.URL.Query().Get("names")
	if namesParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "names parameter is required"}); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	names := strings.Split(namesParam, ",")
	greetings := make([]string, 0, len(names))

	for _, name := range names {
		name = strings.TrimSpace(name)
		if name != "" {
			greetings = append(greetings, fmt.Sprintf("Hello, %s!", name))
		}
	}

	if len(greetings) == 0 {
		greetings = append(greetings, fmt.Sprintf("Hello, %s!", DefaultName))
	}

	response := map[string][]string{
		"greetings": greetings,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
