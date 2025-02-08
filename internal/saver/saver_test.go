package saver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name       string
		outputFile string
		wantErr    bool
	}{
		{
			name:       "valid file path",
			outputFile: "test_output.txt",
			wantErr:    false,
		},
		{
			name:       "empty file path",
			outputFile: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if instance != nil {
					instance.Close()
					if tt.outputFile != "" {
						os.Remove(tt.outputFile)
					}
				}
				instance = nil
			}()

			err := Init(tt.outputFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && instance == nil {
				t.Error("Init() failed to initialize instance")
			}
		})
	}
}

func TestSaveKeyPair(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test_keys.txt")

	tests := []struct {
		name        string
		setup       func()
		publicKey   string
		privateKey  string
		wantErr     bool
		wantContent string
	}{
		{
			name:        "successful save",
			setup:       func() { Init(tempFile) },
			publicKey:   "pub123",
			privateKey:  "priv123",
			wantErr:     false,
			wantContent: "Public: pub123\nPrivate: priv123\n\n",
		},
		{
			name:        "uninitialized saver",
			setup:       func() { instance = nil },
			publicKey:   "pub123",
			privateKey:  "priv123",
			wantErr:     true,
			wantContent: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if instance != nil {
					instance.Close()
				}
				instance = nil
				os.Remove(tempFile)
			}()

			tt.setup()

			err := SaveKeyPair(tt.publicKey, tt.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				content, err := os.ReadFile(tempFile)
				if err != nil {
					t.Fatalf("Failed to read test file: %v", err)
				}
				if string(content) != tt.wantContent {
					t.Errorf("SaveKeyPair() content = %q, want %q", string(content), tt.wantContent)
				}
			}
		})
	}
}

func TestClose(t *testing.T) {
	tempFile := filepath.Join(t.TempDir(), "test_close.txt")

	tests := []struct {
		name    string
		setup   func()
		wantErr bool
	}{
		{
			name:    "successful close",
			setup:   func() { Init(tempFile) },
			wantErr: false,
		},
		{
			name:    "already closed",
			setup:   func() { instance = nil },
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify instance is nil after closing
			if instance != nil {
				t.Error("Close() failed to set instance to nil")
			}

			os.Remove(tempFile)
		})
	}
}
