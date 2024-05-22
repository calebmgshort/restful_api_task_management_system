package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
}

var (
	tasks  = make(map[int]Task)
	nextID = 1
	mu     sync.Mutex
)

func createTask(responseWriter http.ResponseWriter, request *http.Request) {
	var task Task
	json.NewDecoder(request.Body).Decode(&task)
	
	task.Description = "pending"

	mu.Lock()
	task.ID = nextID
	tasks[nextID] = task
	nextID++
	mu.Unlock()

	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(task)
}

func getTask(responseWriter http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))
	mu.Lock()
	task, ok := tasks[id]
	mu.Unlock()

	if !ok {
			http.NotFound(responseWriter, request)
			return
	}

	json.NewEncoder(responseWriter).Encode(task)
}

func main() {
	router := chi.NewRouter()
	router.Post("/tasks", createTask)
	router.Get("/tasks/{id}", getTask)
	// router.Put("/tasks/{id}", updateTask)
	// router.Delete("/tasks/{id}", deleteTask)
	// router.Get("/tasks", getTasks)
	fmt.Println("listening on localhost port 8080")
	http.ListenAndServe(":8080", router)
}