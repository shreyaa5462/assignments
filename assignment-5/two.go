package main

import (
	"fmt"
	"strings"
)

type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{}

func (cl ConsoleLogger) Log(message string) {
	fmt.Printf("Console: %s\n", message)
}

type FileLogger struct {
	loggedMessages []string
}

func (fl *FileLogger) Log(message string) {
	fl.loggedMessages = append(fl.loggedMessages, message)
	fmt.Printf("File: %s (Simulated write to file)\n", message)
}

func (fl FileLogger) GetFileContent() string {
	if len(fl.loggedMessages) == 0 {
		return "No entries logged to file yet."
	}
	return strings.Join(fl.loggedMessages, "\n")
}

type RemoteLogger struct{}

func (rl RemoteLogger) Log(message string) {
	fmt.Printf("Remote: %s (Simulated send to remote server)\n", message)
}

func LogAll(loggers []Logger, message string) {
	fmt.Println("\n--- Logging Message to All Loggers ---")
	for _, logger := range loggers {
		logger.Log(message)
	}
}

func main() {
	fmt.Println("--- Go Logging Application ---")

	consoleLogger := ConsoleLogger{}
	fileLogger := &FileLogger{}
	remoteLogger := RemoteLogger{}

	allLoggers := []Logger{
		consoleLogger,
		fileLogger,
		remoteLogger,
	}

	LogAll(allLoggers, "Hello!")

	LogAll(allLoggers, "User logged in successfully.")

	fmt.Println("\n--- Logging to Individual Loggers ---")
	consoleLogger.Log("This is a direct console message.")
	fileLogger.Log("This is a direct file message.")
	remoteLogger.Log("This is a direct remote message.")

	fmt.Println("\n--- Content of the Simulated Log File ---")
	fmt.Println(fileLogger.GetFileContent())

	fmt.Println("\n--- Application End ---")
}
