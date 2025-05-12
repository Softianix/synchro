// End to End tests for the synchro CLI - test whole app behavior.
package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Devourian/synchro/internal/core"
)

func TestE2E(t *testing.T) {
	// --- ARRANGE ---
	// Create a temporary source and destination directories
	sourceDir, err := os.MkdirTemp("", "sourceDir")
	if err != nil {
		t.Fatalf("failed to create source directory: %v", err)
	}
	defer os.RemoveAll(sourceDir) // Automatic cleanup

	destinationDir, err := os.MkdirTemp("", "destinationDir")
	if err != nil {
		t.Fatalf("failed to create destination directory: %v", err)
	}
	defer os.RemoveAll(destinationDir) // Automatic cleanup

	// Create a test file in the destination directory
	fileName := "testfile.txt"
	destinationFilePath := filepath.Join(destinationDir, fileName)
	destinationFileContent := []byte("Old content")
	err = os.WriteFile(destinationFilePath, destinationFileContent, 0644)
	if err != nil {
		t.Fatalf("failed to create test file in destination dir: %v", err)
	}

	// Create a test file in the source directory
	time.Sleep(5 * time.Millisecond)
	filePath := filepath.Join(sourceDir, fileName)
	fileContent := []byte("Hello, World!")
	err = os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		t.Fatalf("failed to create test file in source dir: %v", err)
	}

	// Create additional test file in the source directory
	additionalFileName := "additionalfile.txt"
	additionalFilePath := filepath.Join(sourceDir, additionalFileName)
	additionalFileContent := []byte("Additional content")
	err = os.WriteFile(additionalFilePath, additionalFileContent, 0644)
	if err != nil {
		t.Fatalf("failed to create additional test file: %v", err)
	}

	//  Create additional test file in the destination directory
	additionalDestinationFilePath := filepath.Join(destinationDir, "additionaldestfile.txt")
	additionalDestinationFileContent := []byte("To be removed")
	err = os.WriteFile(additionalDestinationFilePath, additionalDestinationFileContent, 0644)
	if err != nil {
		t.Fatalf("failed to create additional test file in destination dir: %v", err)
	}

	// --- ACT ---
	core.Run([]string{"synchro", "sync", "--src", sourceDir, "--dst", destinationDir, "--delete_missing"})

	// --- ASSERT ---
	// Verify that the additional file was copied correctly
	additionalCopiedFilePath := filepath.Join(destinationDir, additionalFileName)
	additionalCopiedFileContent, err := os.ReadFile(additionalCopiedFilePath)
	if err != nil {
		t.Fatalf("failed to read copied file: %v", err)
	}
	if string(additionalCopiedFileContent) != string(additionalFileContent) {
		t.Fatalf("copied file content mismatch: expected %s, got %s", string(additionalFileContent), string(additionalCopiedFileContent))
	}

	// Verify that the original file in the destination directory was overwritten
	overwrittenFileContent, err := os.ReadFile(destinationFilePath)
	if err != nil {
		t.Fatalf("failed to read overwritten file: %v", err)
	}
	if string(overwrittenFileContent) != string(fileContent) {
		t.Fatalf("overwritten file content mismatch: expected %s, got %s", string(fileContent), string(overwrittenFileContent))
	}

	// Verify that the additional file in the destination directory was deleted
	_, err = os.ReadFile(additionalDestinationFilePath)
	if err == nil {
		t.Fatalf("expected additional file to be deleted, but it still exists")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("expected file not found error, got: %v", err)
	}
}
