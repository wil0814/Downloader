package download

import (
	"download/file_manager"
	"download/http_client"
	"net/http"
)

type DownloaderFactory interface {
	CanHandle(resp *http.Response) bool
	CreateDownloader(client *http_client.HTTPClient, file *file_manager.FileManager, concurrency int) Downloader
}

type SingleDownloadFactory struct{}

func NewSingleDownloadFactory() *SingleDownloadFactory {
	return &SingleDownloadFactory{}
}
func (f *SingleDownloadFactory) CanHandle(_ *http.Response) bool {
	// Always return true as a fallback
	return true
}

func (f *SingleDownloadFactory) CreateDownloader(client *http_client.HTTPClient, file *file_manager.FileManager, _ int) Downloader {
	return NewSingleDownload(client, file)
}

type MultiDownloadFactory struct{}

func NewMultiDownloadFactory() *MultiDownloadFactory {
	return &MultiDownloadFactory{}
}
func (f *MultiDownloadFactory) CanHandle(resp *http.Response) bool {
	client := &http_client.HTTPClient{}
	return client.IsHttpStatusOK(resp) && client.SupportRange(resp)
}

func (f *MultiDownloadFactory) CreateDownloader(client *http_client.HTTPClient, file *file_manager.FileManager, concurrency int) Downloader {
	return NewMultiDownloader(client, file, concurrency)
}
