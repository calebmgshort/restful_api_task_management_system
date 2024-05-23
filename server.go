package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ServerState struct {
	tasks  map[int]Task
	nextID int
	mu     sync.Mutex
}

// TODO: filter input, making sure the title and description are strings and not empty
func (serverState *ServerState) createTask(responseWriter http.ResponseWriter, request *http.Request) {
	var task Task
	json.NewDecoder(request.Body).Decode(&task)

	task.Status = "Pending"

	serverState.mu.Lock()
	task.ID = serverState.nextID
	serverState.tasks[serverState.nextID] = task
	serverState.nextID++
	serverState.mu.Unlock()

	responseWriter.WriteHeader(http.StatusCreated)
	json.NewEncoder(responseWriter).Encode(task)
}

func (serverState *ServerState) getTask(responseWriter http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))

	// TODO: reconsider if I really need to use the mutex if the transaction is only a single line
	serverState.mu.Lock()
	task, ok := serverState.tasks[id]
	serverState.mu.Unlock()

	if !ok {
		http.NotFound(responseWriter, request)
		return
	}

	json.NewEncoder(responseWriter).Encode(task)
}

// TODO: filter input
func (serverState *ServerState) updateTask(responseWriter http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))
	var updatedTask Task
	json.NewDecoder(request.Body).Decode(&updatedTask)

	serverState.mu.Lock()
	existingTask, ok := serverState.tasks[id]
	if ok {
		if updatedTask.Status != "Pending" && updatedTask.Status != "In Progress" && updatedTask.Status != "Completed" {
			updatedTask.Status = existingTask.Status
		}
		updatedTask.ID = id
		serverState.tasks[id] = updatedTask
	}
	serverState.mu.Unlock()

	if !ok {
		http.NotFound(responseWriter, request)
		return
	}

	json.NewEncoder(responseWriter).Encode(updatedTask)
}

func (serverState *ServerState) deleteTask(responseWriter http.ResponseWriter, router *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(router, "id"))

	serverState.mu.Lock()
	delete(serverState.tasks, id)
	serverState.mu.Unlock()

	responseWriter.WriteHeader(http.StatusNoContent)
}

func (serverState *ServerState) getTasks(responseWriter http.ResponseWriter, request *http.Request) {

	// Convert map to list
	serverState.mu.Lock()
	taskList := make([]Task, 0, len(serverState.tasks))
	for _, task := range serverState.tasks {
		taskList = append(taskList, task)
	}
	serverState.mu.Unlock()

	json.NewEncoder(responseWriter).Encode(taskList)
}

func createRouter() *chi.Mux {
	router := chi.NewRouter()

	serverState := &ServerState{
		tasks:  make(map[int]Task),
		nextID: 1,
	}

	router.Post("/tasks", serverState.createTask)
	router.Get("/tasks/{id}", serverState.getTask)
	router.Put("/tasks/{id}", serverState.updateTask)
	router.Delete("/tasks/{id}", serverState.deleteTask)
	router.Get("/tasks", serverState.getTasks)
	return router
}

func main() {
	router := createRouter()
	fmt.Println("listening on localhost port 8080")
	http.ListenAndServe(":8080", router)
}
