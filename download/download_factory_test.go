package download

import (
	"download/file_manager"
	"download/http_client"
	"net/http"
	"testing"
)

func TestSingleDownloadFactory_CanHandle(t *testing.T) {
	factory := NewSingleDownloadFactory()

	resp := &http.Response{}

	if !factory.CanHandle(resp) {
		t.Error("SingleDownloadFactory should handle any response")
	}
}
func TestSingleDownloadFactory_CreateDownloader(t *testing.T) {
	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()

	factory := NewSingleDownloadFactory()

	downloader := factory.CreateDownloader(client, fileManager, 0)

	_, ok := downloader.(*SingleDownload)
	if !ok {
		t.Error("Expected SingleDownload downloader")
	}
}

func TestMultiDownloadFactory_CanHandle(t *testing.T) {
	factory := NewMultiDownloadFactory()

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Accept-Ranges": []string{"bytes"}},
	}

	if !factory.CanHandle(resp) {
		t.Error("MultiDownloadFactory should handle response with status OK and support for range requests")
	}
}

func TestMultiDownloadFactory_CreateDownloader(t *testing.T) {
	client := http_client.NewHTTPClient()
	fileManager := file_manager.NewFileManager()

	factory := NewMultiDownloadFactory()

	downloader := factory.CreateDownloader(client, fileManager, 4)

	_, ok := downloader.(*MultiDownloader)
	if !ok {
		t.Error("Expected MultiDownloader downloader")
	}
}
