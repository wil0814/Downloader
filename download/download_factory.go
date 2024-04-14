package download

import (
	"download/download/flag"
	"download/download/ftp_download"
	"download/download/ftp_download/ftp_handler"
	"download/download/http_download"
	"download/download/http_download/http_handler"
	"download/file_manager"
	"fmt"
	"net/http"
)

type DownloaderFactory interface {
	CreateDownloader(flag flag.UserFlag) (Downloader, error)
}
type HTTPDownloaderFactory struct{}

func NewHTTPDownloaderFactory() *HTTPDownloaderFactory {
	return &HTTPDownloaderFactory{}
}
func (f *HTTPDownloaderFactory) CreateDownloader(flag flag.UserFlag) (Downloader, error) {
	client := http_handler.NewHTTPClient()
	file := file_manager.NewFileManager()
	resp, err := http.Get(flag.Path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !client.IsHttpStatusOK(resp) {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	if client.SupportRange(resp) {
		fmt.Println("resume download")
		return http_download.NewResumeDownload(client, file, &flag), nil
	} else {
		fmt.Println("standard download")
		return http_download.NewStandardDownload(&flag), nil
	}
}

type FTPDownloaderFactory struct{}

func NewFTPDownloaderFactory() *FTPDownloaderFactory {
	return &FTPDownloaderFactory{}
}
func (f *FTPDownloaderFactory) CreateDownloader(flag flag.UserFlag) (Downloader, error) {
	ftpURLParser := ftp_handler.NewFTPURLParser()
	ftpURL, err := ftpURLParser.Parse(flag.Path)
	if err != nil {
		return nil, err
	}
	return ftp_download.NewFtpDownload(&flag, ftpURL), nil
}
