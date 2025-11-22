package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Data structure for the payload
type Payload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

const serverURL = "http://localhost:8888"
const testKey = "test-key-123"
const testValue = "This is the value for our test data."

func main() {
	// 1. Test POST Request (Write to DB)
	postPayload := Payload{Key: testKey, Value: testValue}
	log.Println("--- Starting POST Request ---")
	if err := postData(postPayload); err != nil {
		log.Fatalf("POST failed: %v", err)
	}

	// Give the database a moment to write (optional, but good for testing)
	time.Sleep(50 * time.Millisecond)

	// 2. Test GET Request (Read from DB)
	log.Println("--- Starting GET Request ---")
	if err := getData(testKey); err != nil {
		log.Fatalf("GET failed: %v", err)
	}
}

// Sends a POST request to the server to write data.
func postData(payload Payload) error {
	jsonPayload, _ := json.Marshal(payload)
	
	start := time.Now()
	resp, err := http.Post(serverURL+"/api/data", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("POST request error: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	log.Printf("POST response status: %s (Duration: %v)", resp.Status, duration)

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("POST failed with status %s. Body: %s", resp.Status, body)
	}
	
	return nil
}

// Sends a GET request to the server to read data.
func getData(key string) error {
	start := time.Now()
	resp, err := http.Get(serverURL + "/api/data/" + key)
	if err != nil {
		return fmt.Errorf("GET request error: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	log.Printf("GET response status: %s (Duration: %v)", resp.Status, duration)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET failed with status %s", resp.Status)
	}

	var result Payload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("Error decoding GET response: %w", err)
	}

	log.Printf("Retrieved Key: %s, Value: %s", result.Key, result.Value)
	return nil
}
