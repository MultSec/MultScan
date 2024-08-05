package main

import (
	"encoding/json"
	"log"
	"fmt"
	"sync"
	"net/http"

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
	serverAddress string
}

func NewTaskStore() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[string]Task)
	ts.serverAddress = ""
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
	defer ts.Unlock()

	task, ok := ts.tasks[id]
	if ok {
		printLog(logInfo, fmt.Sprintf("[%s] Deliverying task", ansi.ColorFunc("default+hb")(id)))

		return task
	} else {
		printLog(logInfo, fmt.Sprintf("[%s] Creating task", ansi.ColorFunc("default+hb")(id)))

		// Creates a new task in the store.
		task := Task{
			BadBytes: 	"",
			Result:  	"Scanning",
		}
	
		ts.tasks[id] = task

		return task
	}
}

func (ts *TaskStore) getTaskStatusHandler(w http.ResponseWriter, req *http.Request) {	
	// Check if server address is set
	if ts.serverAddress == "" {
		printLog(logInfo, fmt.Sprintf("%s %s", ansi.ColorFunc("default+hb")("Server Address: "), ansi.ColorFunc("cyan")(req.RemoteAddr)))
		ts.serverAddress = req.RemoteAddr
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