package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestReadDatabases(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testdb")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testData := "db1\ndb2\ndb3\n"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

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

func createEmptyFile(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "emptyfile")
	if err != nil {
		t.Fatalf("Failed to create empty temp file: %v", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

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

func TestRenameDatabaseInDump(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testdump")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	initialContent := "CREATE DATABASE `testdb`;\nUSE `testdb`;"
	if err := os.WriteFile(tmpFile.Name(), []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	dbName := "testdb"
	err = renameDatabaseInDump(tmpFile.Name(), dbName)
	if err != nil {
		t.Fatalf("renameDatabaseInDump failed: %v", err)
	}

	updatedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read updated file content: %v", err)
	}

	expectedContent := strings.ReplaceAll(initialContent, fmt.Sprintf("`%s`", dbName), fmt.Sprintf("`R4_%s`", dbName))

	if string(updatedContent) != expectedContent {
		t.Errorf("Expected content:\n%s\nGot:\n%s", expectedContent, string(updatedContent))
	}
}
