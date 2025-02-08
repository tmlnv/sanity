// saver/saver.go
package saver

import (
	"fmt"
	"os"
	"sync"
)

var (
	instance *os.File
	mu       sync.Mutex
)

// Init initializes the key saver with the specified output file
func Init(outputFile string) error {
	mu.Lock()
	defer mu.Unlock()

	if outputFile == "" {
		return fmt.Errorf("output file path cannot be empty")
	}

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("could not open output file: %v", err)
	}

	instance = file
	return nil
}

// SaveKeyPair safely writes a public-private key pair to the output file
func SaveKeyPair(publicKey, privateKey string) error {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		return fmt.Errorf("saver not initialized")
	}

	entry := fmt.Sprintf("Public: %s\nPrivate: %s\n\n", publicKey, privateKey)
	_, err := instance.WriteString(entry)
	if err != nil {
		return fmt.Errorf("failed to write key pair: %w", err)
	}

	return nil
}

// Close properly closes the output file
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		err := instance.Close()
		instance = nil
		return err
	}
	return nil
}
