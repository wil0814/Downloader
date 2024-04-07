package file_manager

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateDestFile(t *testing.T) {
	fileManager := NewFileManager()

	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	destFile, err := fileManager.CreateDestFile(tempFile.Name())
	if err != nil {
		t.Fatalf("CreateDestFile() returned error: %v", err)
	}
	defer destFile.Close()

	fileInfo, err := os.Stat(destFile.Name())
	if err != nil {
		t.Fatalf("Error getting file info: %v", err)
	}
	if !fileInfo.Mode().IsRegular() {
		t.Error("Expected file to be regular file")
	}
}
func TestCreatDir(t *testing.T) {
	fileManager := NewFileManager()

	tempFile, err := os.CreateTemp("", "test_file")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	expectedPartDir := fileManager.GetDirName(tempFile.Name())

	os.RemoveAll(expectedPartDir)

	err = fileManager.CreatDir(tempFile.Name())
	if err != nil {
		t.Fatalf("CreatDir() returned error: %v", err)
	}

	if _, err := os.Stat(expectedPartDir); os.IsNotExist(err) {
		t.Errorf("Expected part directory %q to be created, but it does not exist", expectedPartDir)
	}
}

func TestWriteFile(t *testing.T) {
	fileManager := NewFileManager()

	tempDir, err := os.MkdirTemp(".", "test_dir")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fileName := filepath.Base(tempDir)

	data := []byte("Hello, world!")

	reader := bytes.NewReader(data)

	err = fileManager.WriteFile(fileName, 0, reader)
	if err != nil {
		t.Fatalf("WriteFile return error: %v", err)
	}

	partFileName := fileManager.GetPartFileName(fileName, 0)

	content, err := os.ReadFile(partFileName)
	if err != nil {
		t.Fatalf("Error reading file content: %v", err)
	}

	if string(content) != string(data) {
		t.Errorf("Expected file content %s, got %s", string(data), string(content))
	}
}

func TestMerge(t *testing.T) {
	fileManager := NewFileManager()

	concurrency := 3
	partContent := "Hello, world!"
	for i := 0; i < concurrency; i++ {
		partFileName := fileManager.GetPartFileName("test_file.txt", i)

		dirName := filepath.Dir(partFileName)

		if err := os.MkdirAll(dirName, 0755); err != nil {
			t.Fatalf("Error creating directory: %v", err)
		}

		err := os.WriteFile(partFileName, []byte(partContent), 0666)
		if err != nil {
			t.Fatalf("Error creating part file: %v", err)
		}
	}
	defer os.RemoveAll(filepath.Dir(fileManager.GetPartFileName("test_file.txt", 0)))

	err := fileManager.Merge("test_file.txt", concurrency)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}
	defer os.Remove("test_file.txt")

	data, err := os.ReadFile("test_file.txt")
	if err != nil {
		t.Fatalf("Error reading merged file: %v", err)
	}

	expectedContent := strings.Repeat(partContent, concurrency)
	if string(data) != expectedContent {
		t.Errorf("Merged file content incorrect, expected: %s, got: %s", expectedContent, string(data))
	}
}

func TestGetPartFileName(t *testing.T) {
	fileManager := NewFileManager()

	filename := "example.txt"
	partNum := 2
	expected := fmt.Sprintf("%s/%s-%d", fileManager.GetDirName(filename), filename, partNum)
	actual := fileManager.GetPartFileName(filename, partNum)
	if actual != expected {
		t.Errorf("GetPartFileName(%q, %d) returned %q, expected %q", filename, partNum, actual, expected)
	}
}

func TestGetPartDir(t *testing.T) {
	fileManager := NewFileManager()

	filename := "example.txt"
	expected := strings.SplitN(filename, ".", 2)[0]
	actual := fileManager.GetDirName(filename)
	if actual != expected {
		t.Errorf("GetDirName(%q) returned %q, expected %q", filename, actual, expected)
	}
}
