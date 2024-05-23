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

func (serverState *ServerState) createTask(responseWriter http.ResponseWriter, request *http.Request) {
	var task Task
	json.NewDecoder(request.Body).Decode(&task)

	// TODO: Return unique error messages for title and description invalid
	if task.Title == "" || task.Description == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

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

func (serverState *ServerState) updateTask(responseWriter http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))

	var updatedTask Task
	json.NewDecoder(request.Body).Decode(&updatedTask)

	serverState.mu.Lock()
	defer serverState.mu.Unlock()
	_, ok := serverState.tasks[id]
	if !ok {
		http.NotFound(responseWriter, request)
		return
	}

	// TODO: Return unique error messages for invalid title, description and status
	if updatedTask.Title == "" || updatedTask.Description == "" ||
		(updatedTask.Status != "Pending" && updatedTask.Status != "In Progress" && updatedTask.Status != "Completed") {
		// TODO: I'm not sure if this is the right status code for this
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedTask.ID = id
	serverState.tasks[id] = updatedTask

	json.NewEncoder(responseWriter).Encode(updatedTask)
}

func (serverState *ServerState) deleteTask(responseWriter http.ResponseWriter, request *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(request, "id"))

	serverState.mu.Lock()
	defer serverState.mu.Unlock()

	_, ok := serverState.tasks[id]
	if !ok {
		http.NotFound(responseWriter, request)
		return
	}

	delete(serverState.tasks, id)

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
