package cmd

import (
	"os"
	"testing"
)

func TestReadDatabases(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testdb")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write test data to the temporary file
	testData := "db1\ndb2\ndb3\n"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Test cases
	tests := []struct {
		name     string
		filePath string
		want     []string
		wantErr  bool
	}{
		{
			name:     "valid file",
			filePath: tmpFile.Name(),
			want:     []string{"db1", "db2", "db3"},
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filePath: "nonexistent.txt",
			want:     nil,
			wantErr:  true,
		},
		{
			name:     "empty file",
			filePath: createEmptyFile(t),
			want:     []string{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readDatabases(tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readDatabases() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareStringSlices(got, tt.want) {
				t.Errorf("readDatabases() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create an empty file and return its path
func createEmptyFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "emptyfile")
	if err != nil {
		t.Fatalf("Failed to create empty temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

// Helper function to compare two string slices
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
