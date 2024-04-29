package http

import (
	"download/application/model"
	"download/domain/upload/service"
	"download/infrastructure"
	"errors"
	"fmt"
	"github.com/vbauerster/mpb/v8"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

type DownloadParams struct {
	filename      string
	contentLength int
	partSize      int
	bar           *mpb.Bar
}
type ResumeDownload struct {
	fileManager infrastructure.FileManagerInterface
	userFlag    *model.UserFlag
}

func NewResumeDownload(flag *model.UserFlag) *ResumeDownload {
	fileManager := infrastructure.NewFileManager()
	return &ResumeDownload{
		fileManager: fileManager,
		userFlag:    flag,
	}
}

func (d *ResumeDownload) Download() error {
	return d.retrieveHTTP(d.getFileName())
}
func (d *ResumeDownload) getFileName() string {
	if d.userFlag.FileName == "" {
		return path.Base(d.userFlag.Path)
	} else {
		return d.userFlag.FileName
	}
}
func (d *ResumeDownload) retrieveHTTP(fileName string) error {
	resp, err := http.Get(d.userFlag.Path)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	return d.resumeDownload(fileName, int(resp.ContentLength))
}

func (d *ResumeDownload) resumeDownload(fileName string, contentLength int) error {
	concurrency := d.userFlag.Concurrency
	partSize := contentLength / concurrency
	err := d.fileManager.CreatDir(fileName)
	if err != nil {
		return err
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Println(err)
		}
	}(d.fileManager.GetDirName(fileName))

	pb := infrastructure.NewProgressBar()
	bar := pb.CreateBar()
	params := &DownloadParams{
		filename:      fileName,
		contentLength: contentLength,
		partSize:      partSize,
		bar:           bar,
	}
	d.initiateDownloadWorkers(params)

	// Set the total to -1 to indicate that the Download is complete
	bar.SetTotal(-1, true)
	// Wait for all progress bar to complete
	pb.Progress.Wait()

	err = d.fileManager.Merge(fileName, concurrency)
	if err != nil {
		return err
	}

	if string(service.DestinationS3) == d.userFlag.UploadDestination {
		return d.uploadToS3(fileName)
	}

	if err != nil {
		return err
	}
	return nil
}

func (d *ResumeDownload) initiateDownloadWorkers(params *DownloadParams) {
	concurrency := d.userFlag.Concurrency
	var wg sync.WaitGroup
	wg.Add(concurrency)

	segmentStart := 0
	for segmentIndex := 0; segmentIndex < concurrency; segmentIndex++ {

		segmentEnd := segmentStart + params.partSize
		if segmentIndex == concurrency-1 {
			segmentEnd = params.contentLength
		}

		go d.downloadWorker(params, &wg, segmentIndex, segmentStart, segmentEnd)

		segmentStart += params.partSize + 1
	}
	wg.Wait()
}

func (d *ResumeDownload) downloadWorker(params *DownloadParams, wg *sync.WaitGroup, segmentIndex int, segmentStart int, segmentEnd int) {
	defer wg.Done()

	resp, err := d.DownloadPartial(d.userFlag.Path, segmentStart, segmentEnd)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	partFile, err := d.fileManager.CreateDestFile(d.fileManager.GetPartFileName(params.filename, segmentIndex))
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {

		}

		if n > 0 {
			params.bar.SetTotal(params.bar.Current()+int64(n), false)
			_, err := partFile.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
			params.bar.IncrBy(n)
		}

		if err == io.EOF {
			break
		}
	}
}
func (d *ResumeDownload) uploadToS3(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}(fileName)
	defer file.Close()
	u := service.NewUpload(fileName, file, d.userFlag, service.WithS3Uploader())
	return u.Upload(service.DestinationS3)
}
func (d *ResumeDownload) DownloadPartial(url string, start int, end int) (*http.Response, error) {
	resp, err := d.createRangeRequest(url, start, end)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *ResumeDownload) createRangeRequest(url string, start, end int) (*http.Response, error) {
	if start >= end {
		return nil, errors.New("invalid range")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return resp, nil
}
