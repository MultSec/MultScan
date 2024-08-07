package main

import (
	"encoding/json"
	"log"
	"fmt"
	"sync"
	"net/http"
	"strings"
	"flag"

	"github.com/mgutz/ansi"
)

type Task struct {
	BadBytes string  `json:"badBytes"`
	Result string    `json:"result"`
}

// TaskStore is a simple in-memory database of tasks; TaskStore methods are
// safe to call concurrently.
type TaskStore struct {
	sync.Mutex

	tasks  		map[string]Task
	serverIP 	string
	probeConfig map[string]string
}

func NewTaskStore(conf map[string]string) *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[string]Task)
	ts.serverIP = ""
	ts.probeConfig = conf
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

	go ts.backgroundScan(id, ts.probeConfig)

	return task
}

// Function that handles sample download and scans
func (ts *TaskStore) backgroundScan(id string, probeConfig map[string]string) {
	printLog(logInfo, fmt.Sprintf("[%s] Downloading sample from Server", ansi.ColorFunc("default+hb")(id)))
	
    err := GetSample(ts.serverIP, id)
    if err != nil {
        printLog(logError, fmt.Sprintf("[%s] Failed to download sample: %v", ansi.ColorFunc("default+hb")(id), err))
        return
    }

	// Scan sample
    scanResult, badBytes := scanSample(id, probeConfig)

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
	var configPath string
	flag.StringVar(&configPath, "config", "", "Config file to be used.")
	flag.Parse()

	// Retrieve configPath configuration
	conf, err := GetConf(configPath)
	if err != nil {
		printLog(logError, "Failed to read config file")
		return
	}

	mux := http.NewServeMux()
	server := NewTaskStore(conf)

	mux.HandleFunc("GET /task/{id}", server.getTaskStatusHandler)

	log.Fatal(http.ListenAndServe("localhost:9090", mux))
}