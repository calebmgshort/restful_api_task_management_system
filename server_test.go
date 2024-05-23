package main

import (
	"bytes"
	"encoding/json"
	"reflect"

	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTask(t *testing.T) {
	router := createRouter()

	title := "wash dishes"
	description := "make 'em clean"
	task := Task{Title: title, Description: description}

	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("Expected status 201 Created, got %d", res.Code)
	}

	var createdTask Task
	json.NewDecoder(res.Body).Decode(&createdTask)

	expectedTask := Task{ID: 1, Title: title, Description: description, Status: "Pending"}

	if createdTask != expectedTask {
		t.Errorf("Expected task %+v, got %+v", expectedTask, createdTask)
	}
}

func TestGetTask(t *testing.T) {
	router := createRouter()

	title := "wash dishes"
	description := "make 'em clean"
	task := Task{Title: title, Description: description}

	body, _ := json.Marshal(task)

	// Create task
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// Fetch task
	req, _ = http.NewRequest("GET", "/tasks/1", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", res.Code)
	}

	var fetchedTask Task
	json.NewDecoder(res.Body).Decode(&fetchedTask)

	expectedTask := Task{ID: 1, Title: title, Description: description, Status: "Pending"}

	if fetchedTask != expectedTask {
		t.Errorf("Expected task %+v, got %+v", expectedTask, fetchedTask)
	}
}

func TestUpdateTask(t *testing.T) {
	router := createRouter()

	task := Task{Title: "wash dishes", Description: "make 'em clean"}
	body, _ := json.Marshal(task)

	// Create task
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	task.Title = "wash the dishes"
	task.Description = "with lots of soap"
	task.Status = "In Progress"
	body, _ = json.Marshal(task)

	// Update task
	req, _ = http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", res.Code)
	}

	var updatedTask Task
	json.NewDecoder(res.Body).Decode(&updatedTask)

	expectedTask := Task{ID: 1, Title: task.Title, Description: task.Description, Status: task.Status}

	if updatedTask != expectedTask {
		t.Errorf("Expected task %+v, got %+v", expectedTask, updatedTask)
	}
}

func TestDeleteItem(t *testing.T) {
	router := createRouter()

	task := Task{Title: "wash dishes", Description: "make 'em clean"}
	body, _ := json.Marshal(task)

	// Create task
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// Delete task
	req, _ = http.NewRequest("DELETE", "/tasks/1", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 No Content, got %d", res.Code)
	}

	// Try to get the deleted item
	req, _ = http.NewRequest("GET", "/tasks/1", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 Not Found, got %d", res.Code)
	}
}

func TestGetTasks(t *testing.T) {
	router := createRouter()

	title := "wash dishes"
	description := "make 'em clean"
	task := Task{Title: title, Description: description}

	body, _ := json.Marshal(task)

	// Create task
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// Fetch task
	req, _ = http.NewRequest("GET", "/tasks", nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", res.Code)
	}

	var fetchedTasks []Task
	json.NewDecoder(res.Body).Decode(&fetchedTasks)

	expectedTasks := []Task{{ID: 1, Title: title, Description: description, Status: "Pending"}}

	if !reflect.DeepEqual(fetchedTasks, expectedTasks) {
		t.Errorf("Expected task %+v, got %+v", expectedTasks, fetchedTasks)
	}
}
