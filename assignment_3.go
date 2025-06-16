package main

import (
	"bufio"   // For robust line reading from standard input
	"fmt"     // For formatted I/O (printing and reading)
	"os"      // For standard input/output streams
	"strconv" // For converting string to int
	"strings" // For string manipulation like trimming whitespace
)

// Task represents a single task in our tracker.
type Task struct {
	ID          int    // Unique identifier for the task
	Description string // Description of the task
	Completed   bool   // Status: true if completed, false otherwise (pending)
}

// TaskTracker manages the collection of tasks and generates unique IDs.
type TaskTracker struct {
	tasks     []Task     // Slice to store all tasks (both pending and completed)
	nextIDGen func() int // Function closure to generate unique task IDs
}

// idGenerator is a closure that generates unique sequential integer IDs.
// It encapsulates the 'id' counter, so it's not a global variable.
func idGenerator() func() int {
	id := 0 // The 'id' variable is closed over by the returned function.
	return func() int {
		id++ // Increment 'id' each time the returned function is called.
		return id
	}
}

// NewTaskTracker creates and initializes a new TaskTracker instance.
// It also sets up the unique ID generator.
func NewTaskTracker() *TaskTracker {
	return &TaskTracker{
		tasks:     []Task{},      // Initialize an empty slice of tasks
		nextIDGen: idGenerator(), // Assign the closure to generate IDs
	}
}

// AddTask adds a new task to the tracker.
// It uses a pointer receiver (*TaskTracker) because it modifies the TaskTracker's state (its 'tasks' slice).
func (tt *TaskTracker) AddTask(description string) {
	newID := tt.nextIDGen() // Get a unique ID from the closure
	newTask := Task{
		ID:          newID,
		Description: description,
		Completed:   false, // New tasks are always pending
	}
	tt.tasks = append(tt.tasks, newTask) // Add the new task to the slice
	fmt.Printf("Task Added: %d - %s\n", newTask.ID, newTask.Description)
}

// ListTasks displays all pending tasks.
// It uses a pointer receiver (*TaskTracker) because it operates on the TaskTracker's 'tasks' slice,
// even though it doesn't modify it directly in this function (good practice for methods operating on collections).
func (tt *TaskTracker) ListTasks() {
	fmt.Println("\nPending Tasks:")
	foundPending := false
	for _, task := range tt.tasks { // Iterate through all tasks
		if !task.Completed { // Only print tasks that are not completed
			fmt.Printf("%d: %s\n", task.ID, task.Description)
			foundPending = true
		}
	}
	if !foundPending {
		fmt.Println("No pending tasks.")
	}
}

// CompleteTask marks a task as completed given its ID.
// It uses a pointer receiver (*TaskTracker) because it modifies the state of a Task within the tracker's slice.
func (tt *TaskTracker) CompleteTask(id int) {
	taskFound := false
	for i := range tt.tasks { // Iterate using index to allow modification
		if tt.tasks[i].ID == id {
			if tt.tasks[i].Completed {
				fmt.Printf("Task %d is already completed.\n", id)
			} else {
				tt.tasks[i].Completed = true // Mark as completed
				fmt.Printf("Marking task %d as completed: %s\n", id, tt.tasks[i].Description)
			}
			taskFound = true
			break // Exit loop once task is found and updated
		}
	}
	if !taskFound {
		fmt.Printf("Task with ID %d not found.\n", id)
	}
}

// displayMenu prints the interactive menu options to the console.
func displayMenu() {
	fmt.Println("\n--- Personal Task Tracker ---")
	fmt.Println("1. Add a new task")
	fmt.Println("2. List all pending tasks")
	fmt.Println("3. Mark a task as completed")
	fmt.Println("4. Exit")
	fmt.Print("Choose an option: ")
}

// getUserInput reads a line of text from the standard input.
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n') // Read until newline
	return strings.TrimSpace(input)     // Remove leading/trailing whitespace, including newline
}

// main function orchestrates the CLI interaction.
func main() {
	tracker := NewTaskTracker() // Create a new instance of the TaskTracker

	// Main application loop
	for {
		displayMenu()
		choiceStr := getUserInput()
		choice, err := strconv.Atoi(choiceStr) // Convert user input to an integer
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number between 1 and 4.")
			continue // Skip to next loop iteration
		}

		switch choice {
		case 1: // Add Task
			fmt.Print("Enter task description: ")
			description := getUserInput()
			if description == "" {
				fmt.Println("Task description cannot be empty.")
				continue
			}
			tracker.AddTask(description)
		case 2: // List Tasks
			tracker.ListTasks()
		case 3: // Complete Task
			fmt.Print("Enter ID of task to mark as completed: ")
			idStr := getUserInput()
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid ID. Please enter a valid number.")
				continue
			}
			tracker.CompleteTask(id)
		case 4: // Exit
			fmt.Println("Exiting Task Tracker. Goodbye!")
			return // Exit the main function, terminating the program
		default: // Invalid option
			fmt.Println("Invalid option. Please choose a number between 1 and 4.")
		}
	}
}
