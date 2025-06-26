package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type TaskTracker struct {
	db *sql.DB
}

// NewTaskTracker creates a new TaskTracker with database connection
func NewTaskTracker(db *sql.DB) *TaskTracker {
	return &TaskTracker{
		db: db,
	}
}

// InitDB initializes the database connection and creates the tasks table
func InitDB() (*sql.DB, error) {
	// Update these credentials according to your MySQL setup
	// Example: "root:yourpassword@tcp(localhost:3306)/taskdb"
	dsn := "root:shreya123@tcp(localhost:3306)/shreya"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Create tasks table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INT AUTO_INCREMENT PRIMARY KEY,
		description TEXT NOT NULL,
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// AddTask adds a new task to the database
func (tt *TaskTracker) AddTask(description string) (Task, error) {
	query := "INSERT INTO tasks (description, completed) VALUES (?, ?)"
	result, err := tt.db.Exec(query, description, false)
	if err != nil {
		return Task{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, err
	}

	return Task{
		ID:          int(id),
		Description: description,
		Completed:   false,
	}, nil
}

// ListTasks returns all pending tasks as a formatted string
func (tt *TaskTracker) ListTasks() (string, error) {
	query := "SELECT id, description FROM tasks WHERE completed = FALSE ORDER BY id"
	rows, err := tt.db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	s := "Pending Tasks:\n"
	foundPending := false

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Description)
		if err != nil {
			return "", err
		}
		s += fmt.Sprintf("%d: %s\n", task.ID, task.Description)
		foundPending = true
	}

	if !foundPending {
		s += "No pending tasks."
	}

	return s, nil
}

// GetPendingTasks returns all pending tasks as a slice
func (tt *TaskTracker) GetPendingTasks() ([]Task, error) {
	query := "SELECT id, description, completed FROM tasks WHERE completed = FALSE ORDER BY id"
	rows, err := tt.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Description, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTaskByID retrieves a task by its ID
func (tt *TaskTracker) GetTaskByID(id int) (Task, error) {
	query := "SELECT id, description, completed FROM tasks WHERE id = ?"
	row := tt.db.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Description, &task.Completed)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

// CompleteTask marks a task as completed
func (tt *TaskTracker) CompleteTask(id int) (bool, string, error) {
	// First check if task exists and get its details
	task, err := tt.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Sprintf("Task with ID %d not found.", id), nil
		}
		return false, "", err
	}

	if task.Completed {
		return false, fmt.Sprintf("Task %d is already completed.", id), nil
	}

	// Update the task to completed
	query := "UPDATE tasks SET completed = TRUE WHERE id = ?"
	_, err = tt.db.Exec(query, id)
	if err != nil {
		return false, "", err
	}

	return true, fmt.Sprintf("Marking task %d as completed: %s", id, task.Description), nil
}

// DeleteTask removes a task from the database
func (tt *TaskTracker) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	result, err := tt.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return nil
}

// HTTP Handlers

func httppostTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	var input Task
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	task, err := tracker.AddTask(input.Description)
	if err != nil {
		http.Error(w, "Failed to add task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func httpListallTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	pendingTasks, err := tracker.GetPendingTasks()
	if err != nil {
		http.Error(w, "Failed to retrieve tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pendingTasks)
}

func httpDeletetask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = tracker.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Task %d deleted successfully", id)))
}

func httpListbyId(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	task, err := tracker.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found for this ID", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func httpPutTask(w http.ResponseWriter, r *http.Request, tracker *TaskTracker) {
	queryValues := r.URL.Query()
	idStr := queryValues.Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	success, msg, err := tracker.CompleteTask(id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	w.Write([]byte(msg))
}

// CLI Functions (keeping for backward compatibility)

func displayMenu() {
	fmt.Println("\n--- Personal Task Tracker ---")
	fmt.Println("1. Add a new task")
	fmt.Println("2. List all pending tasks")
	fmt.Println("3. Mark a task as completed")
	fmt.Println("4. Exit")
	fmt.Print("Choose an option: ")
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	// Initialize database
	db, err := InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	tracker := NewTaskTracker(db)

	// Setup HTTP routes (compatible with older Go versions)
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			httpListallTask(w, r, tracker)
		case "POST":
			httppostTask(w, r, tracker)
		case "DELETE":
			httpDeletetask(w, r, tracker)
		case "PUT":
			httpPutTask(w, r, tracker)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/task/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			httpListbyId(w, r, tracker)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server starting on :8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
