package main

import (
	"bufio"   // Used for efficient line-by-line file reading
	"fmt"     // Used for formatted input/output (printing to console)
	"log"     // Used for logging fatal errors (e.g., if the file cannot be opened)
	"os"      // Used for accessing command-line arguments and file operations
	"strings" // Used for string manipulation, specifically checking prefixes
)

func main() {
	// Check if exactly one argument (the file path) is provided after the program name.
	if len(os.Args) < 2 {
		prinUsage() // If no file path is provided, print usage instructions.
		os.Exit(1)  // Exit the program with an error code (1 indicates an error).
	}
	// Get the log file path from the command-line arguments.
	logFilePath := os.Args[1]
	// Read a File (Simulate logs): Open the specified log file.
	file, err := os.Open("log.txt")
	if err != nil {
		// If there's an error opening the file (e.g., file not found, permission denied),
		log.Fatalf("Error: Failed to open log file '%s': %v\n", logFilePath, err)
	}
	//  Control Flow: Use defer to ensure the file is closed.
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Warning: Error closing file '%s': %v\n", logFilePath, closeErr)
		}
	}()
	// Initialize counters for different log levels and total lines.
	infoCount := 0
	warningCount := 0
	errorCount := 0
	// Create a new scanner to read the file line by line.
	// bufio.NewScanner is very efficient for reading large files line by line.
	scanner := bufio.NewScanner(file)
	// Control Flow: Use for + range to iterate over file lines.
	// scanner.Scan() advances the scanner to the next token
	// and returns true if a line was successfully read, false otherwise
	for scanner.Scan() {
		line := scanner.Text() // Get the current line as a string.
		switch {
		// We use strings.HasPrefix to check if the line starts with a specific log level tag.
		case strings.HasPrefix(line, "[INFO]"):
			infoCount++
		case strings.HasPrefix(line, "[WARNING]"):
			warningCount++
		case strings.HasPrefix(line, "[ERROR]"):
			errorCount++
		}
	}
	// Check for any errors that occurred during the scanning process
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file '%s': %v\n", logFilePath, err)
		os.Exit(1) // Exit with error if scanning failed.
	}
	// Summary Report: Print the counts and optional bonus features.
	printSummary(logFilePath, infoCount, warningCount, errorCount)
}

// printUsage displays instructions on how to use the CLI tool.
func prinUsage() {
	fmt.Println("Usage: go run main.go <log_file_path>")
	fmt.Println("Example: go run main.go log.txt")
	fmt.Println("\nCounts ERROR, WARNING, and INFO messages in the specified log file.")
}

// printSummary displays the analysis report, including counts, percentages, and analysis time.
func printSummary(logFilePath string, info, warning, errors int) {
	fmt.Printf("Log Analysis of file: %s\n\n", logFilePath)
	fmt.Printf("INFO: %d entries", info)
	fmt.Println()
	fmt.Printf("WARNING: %d entries", warning)
	fmt.Println()
	fmt.Printf("ERROR: %d entries", errors)
	fmt.Println()
}
