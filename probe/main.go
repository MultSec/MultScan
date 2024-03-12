package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"sync"
	"flag"
)

type Response struct {
	BadBytes string `json:"badBytes"`
	Result   string `json:"result"`
}

type Payload struct {
	Payload []byte
	Scanning bool
	sync.Mutex
}

// Function to write the payload to a file
func WritePayload(contents []byte) string {
	tempFile, err := os.CreateTemp("", "payload")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// Write the payload to the temporary file
	tempFile.Write(contents)

	// Return the path of the temporary file
	return tempFile.Name()
}

// Function to check if a file exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

// Funtion that che
func Quarantined(contents []byte) (string, bool) {
	// Write the payload to a file
	path := WritePayload(contents)

	// Sleep for 3 seconds
	time.Sleep(3 * time.Second)

	// Check if the file exists
	if Exists(path) {
		// Check if the file is quarantined
		file, err := os.Open(path)
		if err != nil {
			return path, true
		}
		defer file.Close()
		return path, false
	}
	return path, true
}

// Function to get payload content from a download link
func getPayload(ip string) []byte {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "http://"+ip+":8080/api/v1/payload/download", nil)
	if err != nil {
		fmt.Println(err)
	}

	// Set the headers for the request
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("Content-Type", "application/octet-stream")

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Return the payload
	return body
}

// Function to scan for bad bytes
func scanBadBytes(payload []byte) string {
	// Create a new bad bytes string
	badBytes := ""

	// Return the bad bytes
	return badBytes
}

// Function to run the payload on disk
func runPayload(payload []byte, path string) {
	// Write the payload to a file
	filePath := WritePayload(payload)

	// Run the with fork and run
	

func (p *Payload) Scan(ip string, getBadBytes *bool) Response {
	// Create a new response
	result := Response{BadBytes: "", Result: "Undetected"}

	// Get the payload content
	payload := getPayload(ip)

	// Lock the payload
	p.Lock()
	defer p.Unlock()

	// Set the payload content
	p.Payload = payload

	// Check if the payload is quarantined
	if path, quarantined := Quarantined(payload); quarantined {
		result.Result = "Detected"
		// If the bad bytes are requested scan for them
		if *getBadBytes {
			result.BadBytes = scanBadBytes(payload)
		}

		// Return the result
		return result
	} else {
		// Check for dynamic 
		runPayload(payload, path)
	}
}

func reqHandler(w http.ResponseWriter, r *http.Request, getBadBytes *bool) {
	// Extract the IP address from the request without the port number
	ip := strings.Split(string(r.RemoteAddr), ":")[0]

	// Scan the payload
	result := scan(ip, getBadBytes)

	// Send the response
	json.NewEncoder(w).Encode(result)
}

func main() {
	getBadBytes := flag.Bool("getBadBytes", false, "Get the bad bytes")
	flag.Parse()

	http.HandleFunc("/", reqHandler)
	http.ListenAndServe(":9090", nil)
}