package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestAddTask(t *testing.T) {
	tracker := NewTaskTracker()
	description1 := "Buy groceries"
	output := captureOutput(func() {
		tracker.AddTask(description1)
	})
	if !strings.Contains(output, "Task Added: 1 - Buy groceries") {
		t.Errorf("AddTask failed for '%s': expected confirmation message, got '%s'", description1, output)
	}
	if len(tracker.tasks) != 1 {
		t.Fatalf("AddTask failed: expected 1 task, got %d", len(tracker.tasks))
	}
	if tracker.tasks[0].ID != 1 || tracker.tasks[0].Description != description1 || tracker.tasks[0].Completed != false {
		t.Errorf("AddTask failed: task data incorrect. Expected {1, %s, false}, got %+v", description1, tracker.tasks[0])
	}
	description2 := "Read Go documentation"
	output = captureOutput(func() {
		tracker.AddTask(description2)
	})

	if !strings.Contains(output, "Task Added: 2 - Read Go documentation") {
		t.Errorf("AddTask failed for '%s': expected confirmation message, got '%s'", description2, output)
	}
	if len(tracker.tasks) != 2 {
		t.Fatalf("AddTask failed: expected 2 tasks, got %d", len(tracker.tasks))
	}
	if tracker.tasks[1].ID != 2 || tracker.tasks[1].Description != description2 || tracker.tasks[1].Completed != false {
		t.Errorf("AddTask failed: task data incorrect. Expected {2, %s, false}, got %+v", description2, tracker.tasks[1])
	}
}

func TestListTasks(t *testing.T) {
	tracker := NewTaskTracker()

	// Test Case 1: No pending tasks initially
	output := captureOutput(func() {
		tracker.ListTasks()
	})
	expectedOutput := "\nPending Tasks:\nNo pending tasks.\n"
	if output != expectedOutput {
		t.Errorf("ListTasks failed for empty tracker: expected \n'%s'\n, got \n'%s'\n", expectedOutput, output)
	}

	// Add some tasks
	tracker.AddTask("Task A")
	tracker.AddTask("Task B")
	tracker.AddTask("Task C")
	// Complete one task
	tracker.CompleteTask(2) // Mark Task B as completed

	// Test Case 2: List with some pending and some completed
	output = captureOutput(func() {
		tracker.ListTasks()
	})
	expectedOutput = "\nPending Tasks:\n1: Task A\n3: Task C\n" // Task B should not be listed
	if output != expectedOutput {
		t.Errorf("ListTasks failed with mixed tasks: expected \n'%s'\n, got \n'%s'\n", expectedOutput, output)
	}

	// Test Case 3: All tasks completed
	tracker.CompleteTask(1)
	tracker.CompleteTask(3)
	output = captureOutput(func() {
		tracker.ListTasks()
	})
	expectedOutput = "\nPending Tasks:\nNo pending tasks.\n"
	if output != expectedOutput {
		t.Errorf("ListTasks failed when all tasks completed: expected \n'%s'\n, got \n'%s'\n", expectedOutput, output)
	}
}

// TestCompleteTask tests the functionality of marking tasks as completed.
func TestCompleteTask(t *testing.T) {
	tracker := NewTaskTracker()
	tracker.AddTask("Task to complete 1") // ID 1
	tracker.AddTask("Task to complete 2") // ID 2
	tracker.AddTask("Task to complete 3") // ID 3

	// Test Case 1: Successfully complete an existing task
	output := captureOutput(func() {
		tracker.CompleteTask(2)
	})
	if !strings.Contains(output, "Marking task 2 as completed: Task to complete 2") {
		t.Errorf("CompleteTask failed for existing task: expected completion message, got '%s'", output)
	}
	if !tracker.tasks[1].Completed { // Assuming ID 2 is at index 1
		t.Errorf("CompleteTask failed: Task ID 2 not marked as completed")
	}

	// Test Case 2: Attempt to complete an already completed task
	output = captureOutput(func() {
		tracker.CompleteTask(2) // Try to complete ID 2 again
	})
	if !strings.Contains(output, "Task 2 is already completed.") {
		t.Errorf("CompleteTask failed for already completed task: expected 'already completed' message, got '%s'", output)
	}

	// Test Case 3: Attempt to complete a non-existent task
	output = captureOutput(func() {
		tracker.CompleteTask(999) // Non-existent ID
	})
	if !strings.Contains(output, "Task with ID 999 not found.") {
		t.Errorf("CompleteTask failed for non-existent task: expected 'not found' message, got '%s'", output)
	}

	// Verify that other tasks remain unchanged
	if tracker.tasks[0].Completed || tracker.tasks[2].Completed {
		t.Errorf("CompleteTask unexpectedly modified other tasks.")
	}
}
