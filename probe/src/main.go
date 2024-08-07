package main

import (
	"os"
    "io"
	"encoding/json"
	"log"
	"fmt"
	"sync"
	"net/http"
	"time"
	"strings"
	"math/rand"

	"github.com/mgutz/ansi"
)

type Log int64

const (
    logError Log = iota
    logInfo
    logStatus
    logInput
	logSuccess
	logSection
	logSubSection
)

// Function to print logs
func printLog(log Log, text string) {
	switch log {
	case logError:
		fmt.Printf("[%s] %s %s\n", ansi.ColorFunc("red")("!"), ansi.ColorFunc("red")("ERROR:"), ansi.ColorFunc("cyan")(text))
	case logInfo:
		fmt.Printf("[%s] %s\n", ansi.ColorFunc("blue")("*"), text)
	case logStatus:
		fmt.Printf("[*] %s\n", text)
	case logInput:
		fmt.Printf("[%s] %s", ansi.ColorFunc("yellow")("?"), text)
	case logSuccess:
		fmt.Printf("[%s] %s\n", ansi.ColorFunc("green")("+"), text)
	case logSection:
		fmt.Printf("\t[%s] %s\n", ansi.ColorFunc("yellow")("-"), text)
	case logSubSection:
		fmt.Printf("\t\t[%s] %s\n", ansi.ColorFunc("magenta")(">"), text)
	}
}

type Task struct {
	BadBytes string  `json:"badBytes"`
	Result string    `json:"result"`
}

// TaskStore is a simple in-memory database of tasks; TaskStore methods are
// safe to call concurrently.
type TaskStore struct {
	sync.Mutex

	tasks  map[string]Task
	serverIP string
}

func NewTaskStore() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[string]Task)
	ts.serverIP = ""
	return ts
}

// DeleteTask deletes the task with the given id. If no such id exists, an error
// is returned.
func (ts *TaskStore) DeleteTask(id string) error {
	ts.Lock()
	defer ts.Unlock()

	if _, ok := ts.tasks[id]; !ok {
		return fmt.Errorf("task with id=%d not found", id)
	}

	delete(ts.tasks, id)
	return nil
}

// GetTaskStatus retrieves a task from the store, by id. If no such id exists, 
// then its created
func (ts *TaskStore) GetTaskStatus(id string) Task {
	ts.Lock()
	task, exists := ts.tasks[id]
	if exists {
		ts.Unlock()
		printLog(logInfo, fmt.Sprintf("[%s] Deliverying task status", ansi.ColorFunc("default+hb")(id)))
		return task
	}

	printLog(logInfo, fmt.Sprintf("[%s] Creating task", ansi.ColorFunc("default+hb")(id)))

	// Creates a new task in the store.
	task = Task{
		BadBytes: 	"",
		Result:  	"Scanning",
	}

	ts.tasks[id] = task
	ts.Unlock()

	go ts.backgroundScan(id)

	return task
}

// Function that scans the sample
func scanSample(id string) (string,string) {
    printLog(logInfo, fmt.Sprintf("[%s] Scanning sample...", ansi.ColorFunc("default+hb")(id)))
    
	// TODO
    // Simulate scanning
    time.Sleep(15 * time.Second)
    
    // In a real scenario, you would implement actual scanning logic here
    // For now, we'll use a simple random result
    if rand.Float32() < 0.7 { // 70% chance of detection
        result := "Detected"
        badBytes := "MDAwNDhlM2QKMDAwMDAwMDAgIDY1IDc0IDVmIDYxIDY0IDY0IDY5IDc0ICA2OSA2ZiA2ZSA2MSA2YyA1ZiA3NCA2OSAgfGV0X2FkZGl0aW9uYWxfdGl8CjAwMDAwMDEwICA2MyA2YiA2NSA3NCA3MyAwMCA2NyA2NSAgNzQgNWYgNzQgNjkgNjMgNmIgNjUgNzQgIHxja2V0cy5nZXRfdGlja2V0fAowMDAwMDAyMCAgNzMgMDAgNzMgNjUgNzQgNWYgNzQgNjkgIDYzIDZiIDY1IDc0IDczIDAwIDUzIDc5ICB8cy5zZXRfdGlja2V0cy5TeXwKMDAwMDAwMzAgIDczIDc0IDY1IDZkIDJlIDRlIDY1IDc0ICAyZSA1MyA2ZiA2MyA6YiA2NSA3NCA3MyAgfHN0ZW0uTmV0LlNvY2tldHN8"
        return result, badBytes
    } else {
        return "Undetected", ""
    }
}

// Function that handles sample download and scans
func (ts *TaskStore) backgroundScan(id string) {
	printLog(logInfo, fmt.Sprintf("[%s] Downloading sample from Server", ansi.ColorFunc("default+hb")(id)))
	
    err := GetSample(ts.serverIP, id)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Failed to download sample: %v", ansi.ColorFunc("default+hb")(id), err))
        return
    }

	// Scan sample
    scanResult, badBytes := scanSample(id)

	if scanResult == "Detected" {
		printLog(logSuccess, fmt.Sprintf("[%s] Scan completed: %s", ansi.ColorFunc("default+hb")(id), ansi.ColorFunc("red")("Detected")))
	} else {
		printLog(logSuccess, fmt.Sprintf("[%s] Scan completed: %s", ansi.ColorFunc("default+hb")(id), ansi.ColorFunc("green")("Undetected")))
	}

	// Update status
    ts.Lock()
    defer ts.Unlock()
    
    if task, exists := ts.tasks[id]; exists {
        task.Result = scanResult
        task.BadBytes = badBytes
        ts.tasks[id] = task
    }

	// Delete sample
	printLog(logInfo, fmt.Sprintf("[%s] Deleting sample...", ansi.ColorFunc("default+hb")(id)))
	err = DelSample(id)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Failed to delete sample: %v", ansi.ColorFunc("default+hb")(id), err))
        return
    }
}

// Function to delete the sample file
func DelSample(id string) error {
    filename := fmt.Sprintf("sample_%s.bin", id)
    err := os.Remove(filename)
    if err != nil {
        if os.IsNotExist(err) {
            printLog(logError, fmt.Sprintf("[%s] Sample file not found: %s", ansi.ColorFunc("default+hb")(id), filename))
            return nil
        }
        return fmt.Errorf("failed to delete sample file: %v", err)
    }
    printLog(logSuccess, fmt.Sprintf("[%s] Sample file deleted successfully", ansi.ColorFunc("default+hb")(id)))
    return nil
}

// Function to get sample content from the server
func GetSample(ip string, id string) error {
    // Create a new HTTP request
    req, err := http.NewRequest("GET", "http://"+ip+":8000/api/v1/sample/download/"+id, nil)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Error creating request: %v", ansi.ColorFunc("default+hb")(id), err))
        return err
    }

    // Set the headers for the request
    req.Header.Set("Accept", "application/octet-stream")
    req.Header.Set("Content-Type", "application/octet-stream")

    // Create a new HTTP client
    client := &http.Client{}

    // Send the request
    resp, err := client.Do(req)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Error downloading sample from server: %v", ansi.ColorFunc("default+hb")(id), err))
        return err
    }
    defer resp.Body.Close()

    // Check if the response status is OK
    if resp.StatusCode != http.StatusOK {
        printLog(logError, fmt.Sprintf("[%s] Unexpected status code: %d", ansi.ColorFunc("default+hb")(id), resp.StatusCode))
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    // Save Sample
    err = saveSample(id, resp.Body)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Error saving sample: %v", ansi.ColorFunc("default+hb")(id), err))
        return err
    }

    printLog(logSuccess, fmt.Sprintf("[%s] Sample downloaded and saved successfully", ansi.ColorFunc("default+hb")(id)))
    return nil
}

// Function to save the sample to a file
func saveSample(id string, data io.Reader) error {
    filename := fmt.Sprintf("sample_%s.bin", id)
    out, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer out.Close()

    _, err = io.Copy(out, data)
    if err != nil {
        return fmt.Errorf("failed to save sample: %v", err)
    }
    return nil
}

// Handler for the Tasks Status GET request
func (ts *TaskStore) getTaskStatusHandler(w http.ResponseWriter, req *http.Request) {	
	// Check if server address is set
	if ts.serverIP == "" {
		ts.serverIP = strings.Split(req.RemoteAddr, ":")[0]
		printLog(logInfo, fmt.Sprintf("%s %s", ansi.ColorFunc("default+hb")("Server IP: "), ansi.ColorFunc("cyan")(ts.serverIP)))
	}
  
	id := req.PathValue("id")

	task := ts.GetTaskStatus(id)

	js, err := json.Marshal(task)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusInternalServerError)
	  return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	mux := http.NewServeMux()
	server := NewTaskStore()

	mux.HandleFunc("GET /task/{id}", server.getTaskStatusHandler)

	log.Fatal(http.ListenAndServe("localhost:9090", mux))
}