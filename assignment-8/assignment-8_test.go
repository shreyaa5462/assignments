package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var testTracker *TaskTracker

// Setup: connect to DB once
func init() {
	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	testTracker = NewTaskTracker(db)
}

func TestHTTPPostTask(t *testing.T) {
	payload := []byte(`{"description":"Test task"}`)

	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	httppostTask(w, req, testTracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	var task Task
	json.NewDecoder(resp.Body).Decode(&task)

	if task.Description != "Test task" {
		t.Errorf("Expected description 'Test task', got '%s'", task.Description)
	}
}

func TestHTTPGetAllTasks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	w := httptest.NewRecorder()

	httpListallTask(w, req, testTracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestHTTPPutTask(t *testing.T) {
	// First, add a task
	task, _ := testTracker.AddTask("Task to mark done")

	req := httptest.NewRequest(http.MethodPut, "/task?id="+strconv.Itoa(task.ID), nil)
	w := httptest.NewRecorder()

	httpPutTask(w, req, testTracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestHTTPDeleteTask(t *testing.T) {
	// Add task to delete
	task, _ := testTracker.AddTask("Task to delete")

	req := httptest.NewRequest(http.MethodDelete, "/task?id="+strconv.Itoa(task.ID), nil)
	w := httptest.NewRecorder()

	httpDeletetask(w, req, testTracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestHTTPGetTaskByID(t *testing.T) {
	// Add task
	task, _ := testTracker.AddTask("Task to fetch by ID")

	// Simulate /task/{id} with PathValue
	req := httptest.NewRequest(http.MethodGet, "/task/"+strconv.Itoa(task.ID), nil)

	// Manually set path value (Go 1.22+ only, or else simulate)
	req.SetPathValue("id", strconv.Itoa(task.ID))

	w := httptest.NewRecorder()
	httpListbyId(w, req, testTracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}

	var fetched Task
	json.NewDecoder(resp.Body).Decode(&fetched)

	if fetched.ID != task.ID {
		t.Errorf("Expected task ID %d, got %d", task.ID, fetched.ID)
	}
}
