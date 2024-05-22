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

// TODO: filter input, making sure the title and description are strings
func createTask(responseWriter http.ResponseWriter, request *http.Request) {
	var task Task
	json.NewDecoder(request.Body).Decode(&task)
	
	task.Status = "Pending"

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

	// TODO: reconsider if I really need to use the mutex if the transaction is only a single line
	mu.Lock()
	task, ok := tasks[id]
	mu.Unlock()

	if !ok {
			http.NotFound(responseWriter, request)
			return
	}

	json.NewEncoder(responseWriter).Encode(task)
}

// TODO: filter input
func updateTask(responseWriter http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))
	var updatedTask Task
	json.NewDecoder(request.Body).Decode(&updatedTask)

	mu.Lock()
	existingTask, ok := tasks[id]
	if ok {
			if(updatedTask.Status != "Pending" && updatedTask.Status != "In Progress" && updatedTask.Status != "Completed") {
				updatedTask.Status = existingTask.Status
			}
			updatedTask.ID = id
			tasks[id] = updatedTask
	}
	mu.Unlock()

	if !ok {
			http.NotFound(responseWriter, request)
			return
	}
	
	json.NewEncoder(responseWriter).Encode(request)
}

func deleteTask(responseWriter http.ResponseWriter, router *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(router, "id"))

	mu.Lock()
	delete(tasks, id)
	mu.Unlock()

	responseWriter.WriteHeader(http.StatusNoContent)
}

func getTasks(responseWriter http.ResponseWriter, request *http.Request) {

	// Convert map to list
	mu.Lock()
	taskList := make([]Task, 0, len(tasks))
	for  _, task := range tasks {
		taskList = append(taskList, task)
	}
	mu.Unlock()

	json.NewEncoder(responseWriter).Encode(taskList)
}

func main() {
	router := chi.NewRouter()
	router.Post("/tasks", createTask)
	router.Get("/tasks/{id}", getTask)
	router.Put("/tasks/{id}", updateTask)
	router.Delete("/tasks/{id}", deleteTask)
	router.Get("/tasks", getTasks)
	fmt.Println("listening on localhost port 8080")
	http.ListenAndServe(":8080", router)
}