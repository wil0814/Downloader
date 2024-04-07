package download

import (
	"download/file_manager"
	"download/http_client"
	"net/http"
)

type Downloader interface {
	Download(url string, filename string) error
}

type Download struct {
	downloader Downloader
}

func NewDownload(url string, concurrency int) (*Download, error) {
	client := http_client.NewHTTPClient()
	file := file_manager.NewFileManager()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	factories := []DownloaderFactory{
		NewMultiDownloadFactory(),
		NewSingleDownloadFactory(), // Fallback
	}

	var downloader Downloader
	for _, factory := range factories {
		if factory.CanHandle(resp) {
			downloader = factory.CreateDownloader(client, file, concurrency)
			break
		}
	}

	return &Download{
		downloader: downloader,
	}, nil
}

func (d *Download) Download(url string, filename string) error {
	return d.downloader.Download(url, filename)
}
