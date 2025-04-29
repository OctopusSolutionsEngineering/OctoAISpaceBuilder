package execute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMakeAllExecutable(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "execute_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test directory structure
	files := []string{
		"file1.txt",
		"file2.exe",
		"subdir/file3.sh",
		"subdir/file4.py",
		"subdir/deep/file5",
	}

	// Create the files
	for _, file := range files {
		fullPath := filepath.Join(tempDir, file)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Create the file with no execute permissions
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Call the function being tested
	err = MakeAllExecutable(tempDir)
	if err != nil {
		t.Errorf("MakeAllExecutable returned an error: %v", err)
	}

	// Check all files have executable permission
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Errorf("Error walking path %s: %v", path, err)
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if execute bit is set
		if info.Mode()&0111 == 0 {
			t.Errorf("File %s is not executable", path)
		}

		return nil
	})

	if err != nil {
		t.Errorf("Error walking directory: %v", err)
	}
}
