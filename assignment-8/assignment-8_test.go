package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var tracker *TaskTracker

// ✅ Setup DB connection once for tests
func init() {
	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	tracker = NewTaskTracker(db)
}

func TestAddTask(t *testing.T) {
	payload := []byte(`{"description":"Test task"}`)

	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	httppostTask(w, req, tracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	var task Task
	json.NewDecoder(resp.Body).Decode(&task)

	if task.Description != "Test task" {
		t.Errorf("Expected task description 'Test task', got '%s'", task.Description)
	}
}

func TestGetTasks(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	w := httptest.NewRecorder()

	httpListallTask(w, req, tracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestCompleteTask(t *testing.T) {
	// First, create a task
	task, _ := tracker.AddTask("Task to complete")

	req := httptest.NewRequest(http.MethodPut, "/task?id="+strconv.Itoa(task.ID), nil)
	w := httptest.NewRecorder()

	httpPutTask(w, req, tracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestDeleteTask(t *testing.T) {
	// Add a task to delete
	task, _ := tracker.AddTask("Task to delete")

	req := httptest.NewRequest(http.MethodDelete, "/task?id="+strconv.Itoa(task.ID), nil)
	w := httptest.NewRecorder()

	httpDeletetask(w, req, tracker)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
