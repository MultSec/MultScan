package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"net/http"
)

type Task struct {
	BadBytes string  `json:"badBytes"`
	Result string    `json:"result"`
}

// TaskStore is a simple in-memory database of tasks; TaskStore methods are
// safe to call concurrently.
type TaskStore struct {
	sync.Mutex

	tasks  map[string]Task
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[string]Task)
	return ts
}

// GetTask retrieves a task from the store, by id. If no such id exists, 
// then its created
func (ts *TaskStore) GetTask(id string) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	t, ok := ts.tasks[id]
	if ok {
		return t, nil
	} else {
		return Task{}, fmt.Errorf("task with id=%d not found", id)
	}
}

// CreateTask creates a new task in the store.
func (ts *TaskStore) CreateTask(id string, badBytes string, result string) {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		BadBytes: 	badBytes,
		Result:  	result
	}

	ts.tasks[id] = task
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

func reqHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the remote address from address
		remoteAddr := r.RemoteAddr

		// Download payload
	}
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task at %s\n", req.URL.Path)
  
	id, err := req.PathValue("id")
	if err != nil {
	  http.Error(w, "invalid id", http.StatusBadRequest)
	  return
	}
  
	task, err := ts.store.GetTask(id)
	if err != nil {
	  http.Error(w, err.Error(), http.StatusNotFound)
	  return
	}
  
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
	server := NewTaskServer()

	mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)

	log.Fatal(http.ListenAndServe("localhost:9090", mux))
}