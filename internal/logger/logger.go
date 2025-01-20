package logger

import (
	"os"
	"sync"

	"github.com/charmbracelet/log"
)

var (
	logFile *os.File
	mu      sync.Mutex
)

// InitLogger initializes the logger.
func InitLogger(logFilePath string, enableCLI bool) {
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file", "error", err)
	}

	if !enableCLI {
		log.SetOutput(logFile)
	}
}

// LogResult logs the result to both the file and CLI.
func LogResult(publicKey, privateKey string) {
	mu.Lock()
	defer mu.Unlock()

	log.Info("Vanity address found", "publicKey", publicKey, "privateKey", privateKey)
	logFile.WriteString(publicKey + "\n")
}
