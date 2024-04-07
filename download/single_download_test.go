package download

import (
	"bytes"
	"download/file_manager"
	"download/http_client"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestSingleDownload_Download(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_download")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()
	downloader := NewSingleDownload(client, fileManager)

	expectedContent := []byte("Hello, world!")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(expectedContent)
	}))
	defer server.Close()

	err = downloader.Download(server.URL, path.Join(tempDir, "test_download.txt"))
	if err != nil {
		t.Fatalf("Download() returned error: %v", err)
	}

	downloadedContent, err := os.ReadFile(path.Join(tempDir, "test_download.txt"))
	if err != nil {
		t.Fatalf("Error reading downloaded file: %v", err)
	}
	if !bytes.Equal(downloadedContent, expectedContent) {
		t.Error("Downloaded content does not match expected content")
	}
}

func TestSingleDownload(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_single_download")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()
	downloader := NewSingleDownload(client, fileManager)

	expectedContent := []byte("Hello, world!")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(expectedContent)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Error creating test HTTP response: %v", err)
	}
	defer resp.Body.Close()

	destFile, err := os.CreateTemp(tempDir, "test_download.txt")
	if err != nil {
		t.Fatalf("Error creating destination file: %v", err)
	}
	defer destFile.Close()

	err = downloader.singleDownload(resp, destFile.Name())
	if err != nil {
		t.Fatalf("singleDownload() returned error: %v", err)
	}

	downloadedContent, err := os.ReadFile(destFile.Name())
	if err != nil {
		t.Fatalf("Error reading downloaded file: %v", err)
	}
	if !bytes.Equal(downloadedContent, expectedContent) {
		t.Error("Downloaded content does not match expected content")
	}
}

func TestSingleBufferCopy(t *testing.T) {
	srcData := []byte("Hello, world!")
	destBuffer := bytes.NewBuffer(nil)

	err := NewSingleDownload(nil, nil).BufferCopy(destBuffer, bytes.NewReader(srcData))
	if err != nil {
		t.Fatalf("BufferCopy() returned error: %v", err)
	}

	if !bytes.Equal(destBuffer.Bytes(), srcData) {
		t.Error("Copied content does not match expected content")
	}
}
