package download

import (
	"download/file_manager"
	"download/http_client"
	"download/utils/progress_bar"
	"io"
	"net/http"
	"path"
)

type SingleDownload struct {
	client *http_client.HTTPClient
	file   *file_manager.FileManager
}

func NewSingleDownload(client *http_client.HTTPClient, file *file_manager.FileManager) *SingleDownload {
	return &SingleDownload{
		client: client,
		file:   file,
	}
}
func (d *SingleDownload) Download(url string, filename string) error {
	if filename == "" {
		filename = path.Base(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	return d.singleDownload(resp, filename)
}
func (d *SingleDownload) singleDownload(resp *http.Response, filename string) error {
	destFile, err := d.file.CreateDestFile(filename)
	if err != nil {
		return err
	}
	defer destFile.Close()

	pb := progress_bar.NewProgressBar()
	bar := pb.CreateBar(resp.ContentLength)

	proxyReader := bar.ProxyReader(resp.Body)
	defer proxyReader.Close()

	err = d.BufferCopy(destFile, proxyReader)

	if err != nil {
		return err
	}
	return nil
}

func (d *SingleDownload) BufferCopy(dest io.Writer, src io.Reader) error {
	buf := make([]byte, 32*1024)
	_, err := io.CopyBuffer(dest, src, buf)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
