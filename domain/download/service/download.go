package service

import (
	"download/application/model"
	"download/domain/download"
	"download/domain/download/ftp"
	"download/domain/download/http"
	"fmt"
	"log"
	hp "net/http"
)

type Download struct {
	downloadServices map[model.ProtocolType]download.DownloadInterface
}
type Configuration func(d *Download) error

func NewDownloadService(cfgs ...Configuration) *Download {
	d := &Download{
		downloadServices: make(map[model.ProtocolType]download.DownloadInterface),
	}
	for _, cfg := range cfgs {
		err := cfg(d)
		if err != nil {
			log.Fatal(err)
		}
	}
	return d
}

func WithHTTPDownloaderFactory(userFlag *model.UserFlag) Configuration {
	var service *http.HttpDownload
	var err error
	supportResume := checkCanResumeDownload(*userFlag)
	if supportResume {
		service, err = http.NewHttpService(http.WithResumeDownload(userFlag))
	} else {
		service, err = http.NewHttpService(http.WithStandardDownload(userFlag))
	}
	if err != nil {
		log.Fatal(err)
	}
	return withHTTPDownloader(service)
}

func withHTTPDownloader(d *http.HttpDownload) Configuration {
	return func(os *Download) error {
		os.downloadServices[model.ProtocolHTTP] = d.Download
		return nil
	}
}

func WithFTPDownloader(userFlag *model.UserFlag) Configuration {
	ftpURLParser := ftp.NewFTPURLParser()
	ftpURL, err := ftpURLParser.Parse(userFlag.Path)
	if err != nil {
		log.Fatal(err)
	}
	ftpD := ftp.NewFtpDownload(userFlag, ftpURL)

	return func(os *Download) error {
		os.downloadServices[model.ProtocolFTP] = ftpD
		return nil
	}
}
func (d *Download) Download(key model.ProtocolType) error {
	service, ok := d.downloadServices[key]
	if !ok {
		return fmt.Errorf("no download service for protocol: %v", key)
	}

	err := service.Download()
	if err != nil {
		return err
	}

	return nil
}

func checkCanResumeDownload(userFlag model.UserFlag) bool {
	resp, err := hp.Get(userFlag.Path)
	if err != nil {
		log.Println(err)
	}
	if !isHttpStatusOK(resp) {
		log.Println("HTTP request failed with status code: %d", resp.StatusCode)
	}
	return supportRange(resp)
}
func isHttpStatusOK(resp *hp.Response) bool {
	return resp.StatusCode == hp.StatusOK
}

func supportRange(resp *hp.Response) bool {
	return resp.Header.Get("Accept-Ranges") == "bytes"
}
