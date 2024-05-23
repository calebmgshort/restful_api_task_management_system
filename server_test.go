package main

import (
	"bytes"
	"encoding/json"
	"reflect"

	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestCreateTask(t *testing.T) {

	var router *chi.Mux

	tests := []struct {
		name           string
		inputTask      any
		expectedStatus int
		expectedTask   Task
	}{
		{"Basic create with title and description",
			Task{Title: "wash dishes", Description: "make 'em clean"},
			http.StatusCreated,
			Task{ID: 1, Title: "wash dishes", Description: "make 'em clean", Status: "Pending"},
		},
		{"Extra properties are ignored",
			struct {
				Title       string
				Description string
				Other       int
			}{Title: "wash dishes", Description: "make 'em clean", Other: 5},
			http.StatusCreated,
			Task{ID: 1, Title: "wash dishes", Description: "make 'em clean", Status: "Pending"},
		},
		{"Status is ignored",
			Task{Title: "wash dishes", Description: "make 'em clean", Status: "my custom status"},
			http.StatusCreated,
			Task{ID: 1, Title: "wash dishes", Description: "make 'em clean", Status: "Pending"},
		},
		{"title must be a string",
			struct {
				Title       int
				Description string
			}{Title: 3, Description: "make 'em clean"},
			http.StatusBadRequest,
			Task{},
		},
		{"title must not be blank",
			Task{Title: "", Description: "make 'em clean"},
			http.StatusBadRequest,
			Task{},
		},
		{"description must be a string",
			struct {
				Title       string
				Description int
			}{Title: "title", Description: 4},
			http.StatusBadRequest,
			Task{},
		},
		{"description must not be blank",
			Task{Title: "", Description: "make 'em clean"},
			http.StatusBadRequest,
			Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router = createRouter()

			body, _ := json.Marshal(tt.inputTask)

			req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			if res.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, res.Code)
			}

			if tt.expectedStatus < 400 {
				var createdTask Task
				json.NewDecoder(res.Body).Decode(&createdTask)

				if createdTask != tt.expectedTask {
					t.Errorf("Expected task %+v, got %+v", tt.expectedTask, createdTask)
				}
			}

		})
	}
}

func TestGetTask(t *testing.T) {
	router := createRouter()

	t.Run("id must exist", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/tasks/1", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		if res.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", res.Code)
		}
	})

	t.Run("if id exists, fetches the task", func(t *testing.T) {
		title := "wash dishes"
		description := "make 'em clean"
		task := Task{Title: title, Description: description}

		body, _ := json.Marshal(task)

		// Create task
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		// TODO: tests shouldn't depend on ids incrementing. The id should be taken from the create
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
	})

}

func TestUpdateTask(t *testing.T) {

	var router *chi.Mux

	tests := []struct {
		name           string
		taskToUpdate   any
		expectedStatus int
		expectedTask   Task
	}{
		{"Basic update for title, description and status",
			Task{ID: 1, Title: "wash the dishes", Description: "with lots of soap", Status: "In Progress"},
			http.StatusOK,
			Task{ID: 1, Title: "wash the dishes", Description: "with lots of soap", Status: "In Progress"},
		},
		{"Extra properties are ignored",
			struct {
				ID          int
				Title       string
				Description string
				Status      string
				Other       int
			}{ID: 1, Title: "wash the dishes", Description: "with lots of soap", Status: "Completed", Other: 5},
			http.StatusOK,
			Task{ID: 1, Title: "wash the dishes", Description: "with lots of soap", Status: "Completed"},
		},
		{"title must be a string",
			struct {
				ID          int
				Title       int
				Description string
				Status      string
			}{ID: 1, Title: 3, Description: "make 'em clean", Status: "In Progress"},
			http.StatusBadRequest,
			Task{},
		},
		{"title must not be blank",
			Task{ID: 1, Title: "", Description: "make 'em clean", Status: "In Progress"},
			http.StatusBadRequest,
			Task{},
		},
		{"description must be a string",
			struct {
				ID          int
				Title       string
				Description int
				Status      string
			}{ID: 1, Title: "title", Description: 4, Status: "In Progress"},
			http.StatusBadRequest,
			Task{},
		},
		{"description must not be blank",
			Task{ID: 1, Title: "", Description: "make 'em clean", Status: "In Progress"},
			http.StatusBadRequest,
			Task{},
		},
		// TODO: should probably test status more thoroughly
		{"status must one of the three provided values",
			Task{ID: 1, Title: "", Description: "make 'em clean", Status: "invalid status"},
			http.StatusBadRequest,
			Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router = createRouter()

			body, _ := json.Marshal(Task{Title: "wash dishes", Description: "make 'em clean"})

			req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			body, _ = json.Marshal(tt.taskToUpdate)

			// Update task
			req, _ = http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(body))
			res = httptest.NewRecorder()
			router.ServeHTTP(res, req)

			if res.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, res.Code)
			}

			if tt.expectedStatus < 400 {
				var updatedTask Task
				json.NewDecoder(res.Body).Decode(&updatedTask)

				if updatedTask != tt.expectedTask {
					t.Errorf("Expected task %+v, got %+v", tt.expectedTask, updatedTask)
				}
			}

		})
	}

	t.Run("id must exist", func(t *testing.T) {
		router = createRouter()

		body, _ := json.Marshal(Task{Title: "wash dishes", Description: "make 'em clean"})

		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		body, _ = json.Marshal(Task{ID: 2, Title: "", Description: "make 'em clean", Status: "In Progress"})

		// Update task
		req, _ = http.NewRequest("PUT", "/tasks/2", bytes.NewBuffer(body))
		res = httptest.NewRecorder()
		router.ServeHTTP(res, req)

		if res.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, res.Code)
		}

	})
}

func TestDeleteItem(t *testing.T) {
	router := createRouter()

	t.Run("id must exist", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		if res.Code != http.StatusNotFound {
			t.Errorf("Expected status 404 Not Found, got %d", res.Code)
		}
	})

	t.Run("if id exists, deletes the task", func(t *testing.T) {
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
	})

}

// TODO: consider adding more tests where the task list is longer than 1, and when a task has been deleted
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

	// Fetch tasks
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
