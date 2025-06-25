package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Helper function to create a new TaskTracker for each test
func setupTracker() *TaskTracker {
	return NewTaskTracker()
}

// Test POST /task - Add new task
func TestHttpPostTask(t *testing.T) {
	tracker := setupTracker()

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid task",
			queryParam:     "task=" + url.QueryEscape("Buy groceries"),
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Missing task parameter",
			queryParam:     "",
			expectedStatus: http.StatusOK,
			expectedBody:   "Missing id parameter",
		},
		{
			name:           "Empty task parameter",
			queryParam:     "task=",
			expectedStatus: http.StatusOK,
			expectedBody:   "Missing id parameter",
		},
	}

	for _, tt := range tests {

		req := httptest.NewRequest("POST", "/task?"+tt.queryParam, nil)
		w := httptest.NewRecorder()

		httppostTask(w, req, tracker)

		if w.Code != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
		}

		if tt.expectedBody != "" && !strings.Contains(w.Body.String(), tt.expectedBody) {
			t.Errorf("Expected body to contain '%s', got '%s'", tt.expectedBody, w.Body.String())
		}

	}
}

// Test GET /task - List all tasks
func TestHttpListAllTask(t *testing.T) {
	tracker := setupTracker()

	// Add some test tasks
	tracker.AddTask("Task 1")
	tracker.AddTask("Task 2")

	req := httptest.NewRequest("GET", "/task", nil)
	w := httptest.NewRecorder()

	httpListallTask(w, req, tracker)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Task 1") || !strings.Contains(body, "Task 2") {
		t.Errorf("Expected response to contain both tasks, got: %s", body)
	}

	if !strings.Contains(body, "Pending Tasks:") {
		t.Errorf("Expected response to contain 'Pending Tasks:', got: %s", body)
	}
}

// Test GET /task - List all tasks when empty
func TestHttpListAllTaskEmpty(t *testing.T) {
	tracker := setupTracker()

	req := httptest.NewRequest("GET", "/task", nil)
	w := httptest.NewRecorder()

	httpListallTask(w, req, tracker)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "No pending tasks.") {
		t.Errorf("Expected 'No pending tasks.' in response, got: %s", body)
	}
}

// Test GET /task/{id} - Get task by ID
func TestHttpListById(t *testing.T) {
	tracker := setupTracker()
	tracker.AddTask("Test Task")

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid ID",
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedBody:   "Test Task",
		},
		{
			name:           "Invalid ID format",
			id:             "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid ID format",
		},
		{
			name:           "Non-existent ID",
			id:             "999",
			expectedStatus: http.StatusOK,
			expectedBody:   "Task not found ofr this ID",
		},
	}

	for _, tt := range tests {

		req := httptest.NewRequest("GET", "/task/"+tt.id, nil)
		req.SetPathValue("id", tt.id)
		w := httptest.NewRecorder()

		httpListbyId(w, req, tracker)

		if w.Code != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
		}

		if !strings.Contains(w.Body.String(), tt.expectedBody) {
			t.Errorf("Expected body to contain '%s', got '%s'", tt.expectedBody, w.Body.String())
		}

	}
}

// Test DELETE /task - Delete task
func TestHttpDeleteTask(t *testing.T) {
	tracker := setupTracker()
	tracker.AddTask("Task to delete")
	tracker.AddTask("Task to keep")

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		shouldContain  string
	}{
		{
			name:           "Valid deletion",
			queryParam:     "id=1",
			expectedStatus: http.StatusOK,
			shouldContain:  "Task to keep",
		},
		{
			name:           "Invalid ID format",
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Invalid ID format",
		},
	}

	for _, tt := range tests {

		req := httptest.NewRequest("DELETE", "/task?"+tt.queryParam, nil)
		w := httptest.NewRecorder()

		httpDeletetask(w, req, tracker)

		if w.Code != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
		}

		if !strings.Contains(w.Body.String(), tt.shouldContain) {
			t.Errorf("Expected body to contain '%s', got '%s'", tt.shouldContain, w.Body.String())
		}

	}
}

// Test PUT /task - Complete task
func TestHttpPutTask(t *testing.T) {
	tracker := setupTracker()
	tracker.AddTask("Task to complete")

	tests := []struct {
		name           string
		queryParam     string
		expectedStatus int
		shouldContain  string
	}{
		{
			name:           "Valid completion",
			queryParam:     "id=1",
			expectedStatus: http.StatusOK,
			shouldContain:  "Marking task 1 as completed",
		},
		{
			name:           "Invalid ID format",
			queryParam:     "id=abc",
			expectedStatus: http.StatusBadRequest,
			shouldContain:  "Invalid ID format",
		},
		{
			name:           "Non-existent task",
			queryParam:     "id=999",
			expectedStatus: http.StatusOK,
			shouldContain:  "Task with ID 999 not found",
		},
	}

	for _, tt := range tests {

		req := httptest.NewRequest("PUT", "/task?"+tt.queryParam, nil)
		w := httptest.NewRecorder()

		httpPutTask(w, req, tracker)

		if w.Code != tt.expectedStatus {
			t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
		}

		if !strings.Contains(w.Body.String(), tt.shouldContain) {
			t.Errorf("Expected body to contain '%s', got '%s'", tt.shouldContain, w.Body.String())
		}

	}
}

// Test completing already completed task
func TestHttpPutTaskAlreadyCompleted(t *testing.T) {
	tracker := setupTracker()
	tracker.AddTask("Task to complete twice")

	// Complete the task first time
	tracker.CompleteTask(1)

	req := httptest.NewRequest("PUT", "/task?id=1", nil)
	w := httptest.NewRecorder()

	httpPutTask(w, req, tracker)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedMsg := "Task 1 is already completed"
	if !strings.Contains(w.Body.String(), expectedMsg) {
		t.Errorf("Expected body to contain '%s', got '%s'", expectedMsg, w.Body.String())
	}
}

//// Integration test - Complete workflow
//func TestCompleteWorkflow(t *testing.T) {
//	tracker := setupTracker()
//
//	// 1. Add a task
//	taskDescription := url.QueryEscape("Integration test task")
//	req := httptest.NewRequest("POST", "/task?task="+taskDescription, nil)
//	w := httptest.NewRecorder()
//	httppostTask(w, req, tracker)
//
//	// 2. List all tasks and verify it's there
//	req = httptest.NewRequest("GET", "/task", nil)
//	w = httptest.NewRecorder()
//	httpListallTask(w, req, tracker)
//
//	if !strings.Contains(w.Body.String(), "Integration test task") {
//		t.Error("Task not found in list after adding")
//	}
//
//	// 3. Get specific task by ID
//	req = httptest.NewRequest("GET", "/task/1", nil)
//	req.SetPathValue("id", "1")
//	w = httptest.NewRecorder()
//	httpListbyId(w, req, tracker)
//
//	if !strings.Contains(w.Body.String(), "Integration test task") {
//		t.Error("Task not found when querying by ID")
//	}
//
//	// 4. Complete the task
//	req = httptest.NewRequest("PUT", "/task?id=1", nil)
//	w = httptest.NewRecorder()
//	httpPutTask(w, req, tracker)
//
//	if !strings.Contains(w.Body.String(), "Marking task 1 as completed") {
//		t.Error("Task completion message not correct")
//	}
//
//	// 5. Verify task is no longer in pending list
//	req = httptest.NewRequest("GET", "/task", nil)
//	w = httptest.NewRecorder()
//	httpListallTask(w, req, tracker)
//
//	if strings.Contains(w.Body.String(), "Integration test task") {
//		t.Error("Completed task still appears in pending list")
//	}
//}
//
//// Benchmark tests
//func BenchmarkHttpPostTask(b *testing.B) {
//	tracker := setupTracker()
//
//	for i := 0; i < b.N; i++ {
//		taskParam := url.QueryEscape("Benchmark task")
//		req := httptest.NewRequest("POST", "/task?task="+taskParam, nil)
//		w := httptest.NewRecorder()
//		httppostTask(w, req, tracker)
//	}
//}
//
//func BenchmarkHttpListAllTask(b *testing.B) {
//	tracker := setupTracker()
//	// Add some tasks for realistic benchmark
//	for i := 0; i < 100; i++ {
//		tracker.AddTask("Benchmark task")
//	}
//
//	for i := 0; i < b.N; i++ {
//		req := httptest.NewRequest("GET", "/task", nil)
//		w := httptest.NewRecorder()
//		httpListallTask(w, req, tracker)
//	}
//}
