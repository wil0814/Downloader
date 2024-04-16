package download

import (
	"download/download/ftp_download"
	"download/download/ftp_download/ftp_handler"
	"download/download/http_download"
	"download/download/http_download/http_handler"
	"download/file_manager"
	"download/utils"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DownloaderFactory interface {
	CreateDownloader(flag utils.UserFlag) (Downloader, error)
}
type HTTPDownloaderFactory struct{}

func NewHTTPDownloaderFactory() *HTTPDownloaderFactory {
	return &HTTPDownloaderFactory{}
}
func (f *HTTPDownloaderFactory) CreateDownloader(userFlag utils.UserFlag) (Downloader, error) {
	client := http_handler.NewHTTPClient()
	file := file_manager.NewFileManager()
	resp, err := http.Get(userFlag.Path)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if !client.IsHttpStatusOK(resp) {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	if client.SupportRange(resp) {
		fmt.Println("resume download")
		return http_download.NewResumeDownload(client, file, &userFlag), nil
	} else {
		fmt.Println("standard download")
		return http_download.NewStandardDownload(&userFlag), nil
	}
}

type FTPDownloaderFactory struct{}

func NewFTPDownloaderFactory() *FTPDownloaderFactory {
	return &FTPDownloaderFactory{}
}
func (f *FTPDownloaderFactory) CreateDownloader(userFlag utils.UserFlag) (Downloader, error) {
	ftpURLParser := ftp_handler.NewFTPURLParser()
	ftpURL, err := ftpURLParser.Parse(userFlag.Path)
	if err != nil {
		return nil, err
	}
	return ftp_download.NewFtpDownload(&userFlag, ftpURL), nil
}
