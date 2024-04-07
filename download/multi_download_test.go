package download

import (
	"download/file_manager"
	"download/http_client"
	"download/utils/progress_bar"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestMultiDownload_Download(t *testing.T) {
	// Prepare test data
	fileContent := []byte("test file content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			// Parse the range header
			startStr := strings.TrimPrefix(rangeHeader, "bytes=")
			parts := strings.Split(startStr, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])

			// Write the specified range of the file content
			w.WriteHeader(http.StatusPartialContent)
			if end >= len(fileContent) {
				end = len(fileContent) - 1
			}
			w.Write(fileContent[start : end+1])
		} else {
			// Write the whole file content
			w.Write(fileContent)
		}
	}))
	defer server.Close()

	// Create an instance of MultiDownloader
	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()
	downloader := NewMultiDownloader(client, fileManager, 2)

	// Execute the Download method
	fileName := "test_file.txt"
	url := server.URL
	err := downloader.Download(url, fileName)
	assert.NoError(t, err)
	defer os.RemoveAll(fileName)

	// Check if the downloaded file exists
	fmt.Println("Checking if file exists at:", fileName)
	_, err = os.Stat(fileName)
	assert.NoError(t, err)

	// Check if the content of the downloaded file is correct
	fmt.Println("Reading file from:", fileName)
	downloadedContent, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, downloadedContent)
}

func TestMultiDownloader_multiDownload(t *testing.T) {
	// Prepare test data
	fileContent := []byte("test file content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			// Parse the range header
			startStr := strings.TrimPrefix(rangeHeader, "bytes=")
			parts := strings.Split(startStr, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])

			// Write the specified range of the file content
			w.WriteHeader(http.StatusPartialContent)
			if end >= len(fileContent) {
				end = len(fileContent) - 1
			}
			w.Write(fileContent[start : end+1])
		} else {
			// Write the whole file content
			w.Write(fileContent)
		}
	}))
	defer server.Close()

	// Create an instance of MultiDownloader
	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()
	downloader := NewMultiDownloader(client, fileManager, 2)

	// Execute the multiDownload method
	fileName := "test_file.txt"
	url := server.URL
	err := downloader.multiDownload(url, fileName, len(fileContent))
	assert.NoError(t, err)
	defer os.RemoveAll(fileName)

	// Check if the downloaded file exists
	_, err = os.Stat(fileName)
	assert.NoError(t, err)

	// Check if the content of the downloaded file is correct
	downloadedContent, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, downloadedContent)
}

func TestMultiDownloader_startDownloadHandlers(t *testing.T) {
	// Prepare test data
	fileContent := []byte("test file content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			// Parse the range header
			startStr := strings.TrimPrefix(rangeHeader, "bytes=")
			parts := strings.Split(startStr, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])

			// Write the specified range of the file content
			w.WriteHeader(http.StatusPartialContent)
			if end >= len(fileContent) {
				end = len(fileContent) - 1
			}
			w.Write(fileContent[start : end+1])
		} else {
			// Write the whole file content
			w.Write(fileContent)
		}
	}))
	defer server.Close()

	// Create an instance of MultiDownloader
	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()
	downloader := NewMultiDownloader(client, fileManager, 2)

	// Prepare the parameters for startDownloadHandlers
	fileName := "test_file.txt"

	url := server.URL
	err := fileManager.CreatDir(fileName)
	if err != nil {
		t.Errorf("Error creating directory: %v", err)
	}
	defer os.RemoveAll(fileManager.GetDirName(fileName))

	contentLength := len(fileContent)
	partSize := contentLength / downloader.concurrency
	pb := progress_bar.NewProgressBar()
	bar := pb.CreateBar(int64(contentLength))

	// Execute the startDownloadHandlers method
	downloader.startDownloadHandlers(url, fileName, contentLength, partSize, bar)

	// Check if the downloaded file exists
	_, err = os.Stat(fileManager.GetPartFileName(fileName, 0))
	assert.NoError(t, err)

	// Check if the content of the downloaded file is correct
	downloadedContent, err := os.ReadFile(fileManager.GetPartFileName(fileName, 0))
	assert.NoError(t, err)
	assert.Equal(t, fileContent[:partSize+1], downloadedContent)
}

func TestMultiDownloader_downloadHandler(t *testing.T) {
	// Prepare test data
	fileContent := []byte("test file content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" {
			// Parse the range header
			startStr := strings.TrimPrefix(rangeHeader, "bytes=")
			parts := strings.Split(startStr, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])

			// Write the specified range of the file content
			w.WriteHeader(http.StatusPartialContent)
			if end >= len(fileContent) {
				end = len(fileContent) - 1
			}
			w.Write(fileContent[start : end+1])
		} else {
			// Write the whole file content
			w.Write(fileContent)
		}
	}))
	defer server.Close()

	// Create an instance of MultiDownloader
	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()
	downloader := NewMultiDownloader(client, fileManager, 2)

	// Prepare the parameters for downloadHandler
	fileName := "test_file.txt"
	url := server.URL
	err := fileManager.CreatDir(fileName)
	if err != nil {
		t.Errorf("Error creating directory: %v", err)
	}
	defer os.RemoveAll(fileManager.GetDirName(fileName))

	contentLength := len(fileContent)
	partSize := contentLength / downloader.concurrency
	var wg sync.WaitGroup
	wg.Add(1)
	pb := progress_bar.NewProgressBar()
	bar := pb.CreateBar(int64(contentLength))

	// Execute the downloadHandler method
	go downloader.downloadHandler(url, fileName, contentLength, partSize, &wg, 0, 0, bar)
	wg.Wait()

	// Check if the downloaded file exists
	_, err = os.Stat(fileManager.GetPartFileName(fileName, 0))
	assert.NoError(t, err)

	// Check if the content of the downloaded file is correct
	downloadedContent, err := os.ReadFile(fileManager.GetPartFileName(fileName, 0))
	assert.NoError(t, err)
	assert.Equal(t, fileContent[:partSize+1], downloadedContent)
}
